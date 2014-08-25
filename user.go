package main

/*
	This struct describes an user.
	email and password is required.

	POST /api/login
	{
		"email":"foo@bar",
		"password":"foobar"
	}
*/
type user struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
