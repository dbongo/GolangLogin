package main

type KTH struct {
	Username string `xml:"cas:user"`
}

type KTH_USER struct {
	Image      string `json:"url"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	Email      string `json:"email"`
}
