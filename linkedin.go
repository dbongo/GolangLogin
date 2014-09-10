package main

import (
	"github.com/golang/oauth2"
)

func newLinkedInConf() (*oauth2.Config, error) {
	linkedInConf, err := oauth2.NewConfig(&oauth2.Options{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:3000/api/callback/linkedin",
		Scopes:       []string{"email"},
	},
		"",
		"")
	if err != nil {
		return nil, err
	}
	return linkedInConf, nil
}
