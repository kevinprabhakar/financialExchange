package util

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var TempAuthKey = "rb8vOJsvBfAgK3IkktBt"

func GetAccessToken(uid string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["uid"] = uid
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(TempAuthKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyAccessToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(TempAuthKey), nil
	})

	if err == nil && token.Valid {
		claims := token.Claims.(jwt.MapClaims)

		return claims["uid"].(string), nil
	} else {
		return "", err
	}
}
