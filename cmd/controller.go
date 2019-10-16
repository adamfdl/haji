package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/adamfdl/tenx/domain"

	"github.com/adamfdl/tenx/third_party/harvest"
	"github.com/adamfdl/tenx/third_party/jira"
	"github.com/jmoiron/sqlx"

	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type controller struct {
	log *zap.Logger
	db  *sqlx.DB
}

func newController(log *zap.Logger, db *sqlx.DB) controller {
	return controller{
		log: log,
		db:  db,
	}
}

func (c controller) hAuthJira(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, jira.OAuthConsentURL(), http.StatusTemporaryRedirect)
}

func (c controller) hAuthJiraCallback(w http.ResponseWriter, r *http.Request) {
	token, err := jira.OAuthExchange(r.FormValue("code"))
	if err != nil {
		badRequestResponse(w)
		return
	}

	writeJSONResponse(w, http.StatusOK, token)
}

func (c controller) hAuthHarvest(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, harvest.OAuthConsentURL(), http.StatusTemporaryRedirect)
}

func (c controller) hAuthHarvestCallback(w http.ResponseWriter, r *http.Request) {
	token, err := harvest.OAuthExchange(r.FormValue("code"))
	if err != nil {
		badRequestResponse(w)
		return
	}

	writeJSONResponse(w, http.StatusOK, token)
}

func (c controller) hUserSignUp(w http.ResponseWriter, r *http.Request) {
	type signUpIn struct {
		HarvestToken     oauth2.Token `json:"harvest_token"`
		JiraToken        oauth2.Token `json:"jira_token"`
		Email            string       `json:"email"`
		FcmID            string       `json:"fcm_id"`
		SeniorityLevelID int          `json:"seniority_level_id"`
	}

	in := signUpIn{}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		c.log.Error("signup.unmarshal", zap.Error(err))
		badRequestResponse(w)
		return
	}

	tx, _ := c.db.Begin()

	// TODO: Move this to service layer
	var userID int
	err := tx.QueryRow(`
	INSERT INTO users (seniority_level_id, email, fcm_id)
		VALUES
			($1, $2, $3) RETURNING id`, in.SeniorityLevelID, in.Email, in.FcmID).Scan(&userID)
	if err != nil {
		tx.Rollback()
		c.log.Error("signup.insert.user", zap.Error(err))
		internalServerErrorResponse(w)
		return
	}

	_, err = tx.Exec(`
	INSERT INTO harvest_tokens (access_token, refresh_token, expiry, users_id)
		VALUES
			($1, $2, $3, $4)`, in.HarvestToken.AccessToken, in.HarvestToken.RefreshToken, in.HarvestToken.Expiry, userID)
	if err != nil {
		tx.Rollback()
		c.log.Error("signup.insert.harvest_token", zap.Error(err))
		internalServerErrorResponse(w)
		return
	}

	_, err = tx.Exec(`
	INSERT INTO jira_tokens (access_token, refresh_token, expiry, users_id)
		VALUES 
			($1, $2, $3, $4)`, in.JiraToken.AccessToken, in.JiraToken.RefreshToken, in.JiraToken.Expiry, userID)
	if err != nil {
		tx.Rollback()
		c.log.Error("signup.insert.jira_token", zap.Error(err))
		internalServerErrorResponse(w)
		return
	}

	tx.Commit()

	var user domain.User
	err = c.db.Get(&user, `
		SELECT u.id, u.email, sl."level" "seniority_level", u.fcm_id, se.velocity "story_velocity" FROM users u 
			INNER JOIN seniority_levels sl ON sl.id = u.seniority_level_id
			INNER JOIN seniority_estimation se ON se.seniority_levels_id = sl.id
		WHERE u.id = $1`, userID)
	if err != nil {
		c.log.Error("signup.get.user", zap.Error(err))
	}

	// Set jwt claims
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":              user.ID,
		"email":           user.Email,
		"fcm_id":          user.FcmID,
		"seniority_level": user.SeniorityLevel,
		"story_velocity":  user.StoryVelocity,
	})

	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		c.log.Error("signup.jwt.sign", zap.Error(err))
		internalServerErrorResponse(w)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"token": tokenString,
	})
}
