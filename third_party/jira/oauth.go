package jira

import (
	"golang.org/x/oauth2"
)

const (
	// TODO: Update this to not be static
	csrfMitigation = "asdj1g2i3uy12t87s6db19726b31"
)

var (
	oauthEndpoint = oauth2.Endpoint{
		AuthURL:  "https://auth.atlassian.com/authorize",
		TokenURL: "https://auth.atlassian.com/oauth/token",
	}
	oauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/auth/jira/callback",
		ClientID:     "N6XTi79CmCvhnWYszR70dypEQFtHA3qH",
		ClientSecret: "DWEQC92LBtFAchTsuHnFcf_0if7t8KhFxgM_68dNHKYxZ--n6PwSRFyEckem6X9n",
		Scopes:       []string{"read:jira-user", "read:jira-work", "offline_access"},
		Endpoint:     oauthEndpoint,
	}
)

func OAuthConsentURL() string {
	return oauthConfig.AuthCodeURL(csrfMitigation,
		oauth2.SetAuthURLParam("prompt", "consent"),
		oauth2.SetAuthURLParam("audience", "api.atlassian.com"))
}

func OAuthExchange(code string) (*oauth2.Token, error) {
	return oauthConfig.Exchange(oauth2.NoContext, code)
}

// Get token from in-memory database
func getTokenFromSomewhere() (*oauth2.Token, error) {
	return &oauth2.Token{}, nil
}

func test() {
	// oauthToken, err := getTokenFromSomewhere()
	// if err != nil {
	// 	// Handle error gracefully
	// }

	//"/rest/api/3/search"
	// jiraClient, err := jira.NewClient(oauthConfig.Client(oauth2.NoContext, oauthToken), baseURL)
	// if err != nil {
	// 	// Handle error gracefully
	// }

	// req, err := jiraClient.NewRequest(http.MethodGet, "/rest/api/3/search", nil)
	// if err != nil {
	// 	// Handle error gracefully
	// }

	// issues := new([]jira.Issue)
	// _, err := jiraClient.Do(req, issues)
	// if err != nil {
	// 	// Handle error gracefully
	// }

}
