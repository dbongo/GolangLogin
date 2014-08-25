package main

/*
	This struct describes an user.
*/
type user struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
