package configs

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("G_CLIENT"),
		ClientSecret: os.Getenv("G_SECRET"),
		RedirectURL:  os.Getenv("G_REDIRECT"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}
