package auth

import (
	"errors"
	"time"

	"example.com/simple-api/models"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secretkey")

func GenerateJWT(email string) (tokenString string, err error) {
	expTime := time.Now().Add(1 * time.Hour)
	claims := &models.AuthCustomClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&models.AuthCustomClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*models.AuthCustomClaims)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token already expired")
		return
	}
	return
}
