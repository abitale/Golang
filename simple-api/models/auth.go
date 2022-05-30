package models

import "github.com/dgrijalva/jwt-go"

type AuthCustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
