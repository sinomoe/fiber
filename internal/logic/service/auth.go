package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sinomoe/fiber/internal/logic/dto"
)

type Auth struct {
	secret []byte
}

var (
	authSvc  *Auth
	authOnce sync.Once
)

func InitAuth(secret string) {
	authOnce.Do(func() {
		authSvc = &Auth{
			secret: []byte(secret),
		}
	})
}

func GetAuth() *Auth {
	return authSvc
}

func (a Auth) BuildToken(username string) (string, error) {
	claims := dto.AuthClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secret)
}

func (a Auth) ParseToken(token string) (claims dto.AuthClaims, err error) {
	if _, err = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secret, nil
	}); err != nil {
		return
	}
	return
}
