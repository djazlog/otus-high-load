package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"otus-project/internal/model"
	"time"
)

const (
	accessTokenSecretKey  = "kjhbdsfkgjhKJHBKJHbdfsg-sf-asdf"
	accessTokenExpiration = 5 * time.Minute
)

func GenerateToken(userId string) (string, error) {
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        userId,
			ExpiresAt: time.Now().Add(accessTokenExpiration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(accessTokenSecretKey))
}

func VerifyToken(tokenStr string) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.Errorf("unexpected token signing method")
			}

			return []byte(accessTokenSecretKey), nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.Errorf("invalid token claims")
	}

	return claims, nil
}
