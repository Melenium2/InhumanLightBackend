package jwtHelper

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/inhumanLightBackend/app/models"
)

const (
	secret = "21a481f3d02a91f8b6a9e9c1c2e1ce11d9e5d18fa23673bcdfdfedfa39ea393665aff383706b0c4666e092e048766138b605307042a208e6cc21fcb2e9ed8f24" 
)

func CreateJwtToken(u *models.User, time int64) (string, error) {
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS512, &jwt.StandardClaims{
		Id: string(u.ID),
		ExpiresAt: time,
		Subject: u.Role,
	})

	return jwt.SignedString([]byte(secret))
}

func ValidateJwtToken() error {
	return nil
}