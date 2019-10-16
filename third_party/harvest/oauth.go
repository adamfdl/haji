package harvest

import (
	"golang.org/x/oauth2"
)

const (
	// TODO: Update this to not be static
	csrfMitigation = "asdj1g2i3uy12t87s6db19726b31"
)

var (
	oauthEndpoint = oauth2.Endpoint{
		AuthURL:  "https://id.getharvest.com/oauth2/authorize",
		TokenURL: "https://id.getharvest.com/api/v2/oauth2/token",
	}
	oauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/auth/harvest/callback",
		ClientID:     "ie8OTchfVT26XHooCuiM8UVW",
		ClientSecret: "6gKKC2ItOTkqsn26NUKSVIv0wZ4LLKNE1wune9TVyfEQRQ6o5g_9bwgNEaDUu5yyt_SXhrebUkrEk0Y1l4tjmA",
		Scopes:       []string{"harvest:all"},
		Endpoint:     oauthEndpoint,
	}
)

func OAuthConsentURL() string {
	return oauthConfig.AuthCodeURL(csrfMitigation)
}

func OAuthExchange(code string) (*oauth2.Token, error) {
	return oauthConfig.Exchange(oauth2.NoContext, code)
}

func OauthRefresh() {
}
