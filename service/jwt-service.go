package service

import (
	"errors"
	"messanger/models"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWTData struct {
	jwt.StandardClaims
	models.TokenParams
}

func GenerateToken(user_id int64, login string) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
		TokenParams: models.TokenParams{
			Login: login,
			ID:    user_id,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)

	return tokenString, err
}

func DecodeToken(stringToken string) (int64, error) {

	jwtData := JWTData{}
	token, err := jwt.ParseWithClaims(stringToken, &jwtData, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("Token is invalid")
	}
	return jwtData.TokenParams.ID, nil
}
