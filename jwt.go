package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKeyUser = []byte("secret-key")
var secretKeyBusiness = []byte("secret-key-business")

func createToken(username string, utype string) (string, error) {
	var secretKey []byte
	if utype == "USER" {
		secretKey = secretKeyUser
	} else {
		secretKey = secretKeyBusiness
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string, utype string) (jwt.MapClaims, bool) {
	var secretKey []byte
	if utype == "USER" {
		secretKey = secretKeyUser
	} else {
		secretKey = secretKeyBusiness
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true

	}

	return nil, false
}
