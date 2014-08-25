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
	Id       string `gorethink:"id,omitempty"`
	Email    string `gorethink:"email" binding:"required"`
	Password string `gorethink:"password" binding:"required"`
}
