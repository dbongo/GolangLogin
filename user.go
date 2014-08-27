package main

import (
	"gopkg.in/mgo.v2/bson"
)

/*
	These structs describes an user.
	Email and password is required.

	POST /api/login
	{
		"email":"foo@bar",
		"password":"foobar"
	}
*/
type LocalUser struct {
	//If UserID is empty, don't create json with an empty field.
	UserID   bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Email    string        `json:"email" binding:"required"`
	Password string        `json:"password" binding:"required"`
}

// Use this type if you need to show the user for a user.
type ShowLocalUser struct {
	UserID bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Email  string        `json:"email" binding:"required"`
}

// Use this type if you need to update an user.
type UpdateLocalUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Type for a new facebook user.
type FacebookUser struct {
	// Used for updating profile pic.
	FB_ID    string `json:"id"`
	Fullname string `json:"name"`
	Email    string `json:"email"`
}

// Type for a new google user.
type GoogleUser struct {
	Google_ID string `json:"id"`
	Fullname  string `json:"name"`
	Email     string `json:"email"`
}

// Type for a new linkedin user.
type LinkedinUser struct {
	Linkedin_ID string `json:"id"`
	Fullname    string `json:"name"`
	Email       string `json:"email"`
}
