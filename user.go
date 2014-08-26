package main

import (
	"gopkg.in/mgo.v2/bson"
)

/*
	This struct describes an user.
	email and password is required.

	POST /api/login
	{
		"email":"foo@bar",
		"password":"foobar"
	}
*/
type User struct {
	//If UserID is empty, don't create json with an empty field.
	UserID   bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Email    string        `json:"email" binding:"required"`
	Password string        `json:"password" binding:"required"`
}
