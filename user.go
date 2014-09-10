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
	//If MongoID is empty, don't create json with an empty field.
	MongoID  bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Email    string        `json:"email" binding:"required"`
	Password string        `json:"password" binding:"required"`
}

// Use this type if you need to show the user for a user.
type ShowLocalUser struct {
	MongoID bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Email   string        `json:"email" binding:"required"`
}

// Use this type if you need to update an user.
type UpdateLocalUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Type for a new facebook user.
type FacebookUser struct {
	MongoID bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	// Used for updating profile pic.
	ID        string   `json:"id" bson:"facebookID"`
	Fullname  string   `json:"name"`
	Email     string   `json:"email"`
	Education []string `json:"education"`
}

// Type for a new google user.
type GoogleUser struct {
	MongoID  bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	ID       string        `json:"id" bson:"googleID"`
	Fullname string        `json:"name"`
	Email    string        `json:"email"`
}

// Type for a new linkedin user.
type LinkedinUser struct {
	MongoID  bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	ID       string        `json:"id" bson:"linkedinID"`
	Fullname string        `json:"name"`
	Email    string        `json:"email"`
}
