package jwtHelper

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/inhumanLightBackend/app/models"
)

var (
	secret = []byte("21a481f3d02a91f8b6a9e9c1c2e1ce11d9e5d18fa23673bcdfdfedfa39ea393665aff383706b0c4666e092e048766138b605307042a208e6cc21fcb2e9ed8f24")
)

type claims struct {
	UserId int    `json:"user_id"`
	Access string `json:"access"`
	Type   string `json:"token_type"`
	jwt.StandardClaims
}

func Create(u *models.User, days uint8, tokenType string) (string, error) {
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS512, &claims{
		UserId: u.ID,
		Access: u.Role,
		Type:   tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(days)).Unix(),
		},
	})

	return jwt.SignedString(secret)
}

func Validate(token string) (*claims, error) {
	claims := &claims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil && !jwtToken.Valid {
		return nil, err
	}

	return claims, nil
}
