package main

import (
	"errors"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//Gin middleware. Checks for valid tokens in the http auth header.
func tokenMiddleWare(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := jwt_lib.ParseFromRequest(c.Request, func(token *jwt_lib.Token) ([]byte, error) {
			return []byte(secret), nil
		})
		if err != nil {
			c.Fail(401, err)
		}
	}
}

//Generate a valid token. Put this is in the auth header when making calls to auth routes.
func generateToken(secret []byte, claims map[string]interface{}) (string, error) {
	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
	token.Claims = claims
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", errors.New("An error occured while generating token")
	}
	return tokenString, nil
}
