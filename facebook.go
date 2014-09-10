package main

import (
	"github.com/golang/oauth2"
)

func newFacebookConf() (*oauth2.Config, error) {
	facebookConf, err := oauth2.NewConfig(&oauth2.Options{
		ClientID:     "742060722554202",
		ClientSecret: "14c8d1bc1bc23ecd670aa41d6b5348de",
		RedirectURL:  "http://localhost:3000/api/callback/facebook",
		Scopes:       []string{"email", "user_education_history"},
	},
		"https://graph.facebook.com/oauth/authorize",
		"https://graph.facebook.com/oauth/access_token")
	if err != nil {
		return nil, err
	}
	return facebookConf, nil
}
