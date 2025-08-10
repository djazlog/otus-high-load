package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"net/http"
	"otus-project/internal/model"
	"strings"
	"time"
)

const (
	accessTokenSecretKey  = "kjhbdsfkgjhKJHBKJHbdfsg-sf-asdf"
	accessTokenExpiration = 5 * time.Hour
	JWTClaimsContextKey   = "jwt_claims"
)

var (
	ErrNoAuthHeader      = errors.New("authorization header is missing")
	ErrInvalidAuthHeader = errors.New("authorization header is malformed")
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

func GetUserFromToken(req *http.Request) (*string, error) {
	// Now, we need to get the JWS from the request, to match the request expectations
	// against request contents.
	jws, err := GetJWSFromRequest(req)
	if err != nil {
		return nil, fmt.Errorf("getting jws: %w", err)
	}

	// if the JWS is valid, we have a JWT, which will contain a bunch of claims.
	cl, err := VerifyToken(jws)
	if err != nil {
		return nil, fmt.Errorf("token claims don't match: %w", err)
	}

	return &cl.UserId, nil
}

func GetJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	// Check for the Authorization header.
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	// We expect a header value of the form "Bearer <token>", with 1 space after
	// Bearer, per spec.
	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}
