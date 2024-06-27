package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func GetUserID(jwtToken, secret string) (uuid.UUID, error) {
	token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("token is invalid")
	}

	claims := token.Claims.(jwt.MapClaims)

	idRaw, ok := claims["id"].(string)
	if !ok {
		return uuid.Nil, errors.New("token is invalid")
	}

	id, err := uuid.Parse(idRaw)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
