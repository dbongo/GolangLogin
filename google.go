package main

import (
	"github.com/golang/oauth2"
)

func newGoogleConf() (*oauth2.Config, error) {
	googleConf, err := oauth2.NewConfig(&oauth2.Options{
		ClientID:     "979574808763-b093jvno4lgovmbusq5j402s3neo2tm2.apps.googleusercontent.com",
		ClientSecret: "2ILaAFIY2MsCwKWBHqIabXNA",
		RedirectURL:  "http://localhost:3000/api/callback/google",
		Scopes:       []string{"email"},
	},
		"https://accounts.google.com/o/oauth2/auth",
		"https://accounts.google.com/o/oauth2/token")
	if err != nil {
		return nil, err
	}
	return googleConf, nil
}
