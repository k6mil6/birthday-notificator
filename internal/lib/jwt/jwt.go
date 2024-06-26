package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/k6mil6/birthday-notificator/internal/model"
	"time"
)

type key int

const (
	KeyUserID key = iota
)

func NewToken(user model.User, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	return token.SignedString([]byte(secret))
}

func GetUserID(jwtToken string, secret string) (int64, error) {
	token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("token is invalid")
	}

	claims := token.Claims.(jwt.MapClaims)

	idFloat, ok := claims["id"].(int64)
	if !ok {
		return 0, errors.New("ID claim is not a number")
	}

	return idFloat, nil
}
