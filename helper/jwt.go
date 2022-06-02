package helper

import (
	"auth/config"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"
)

const ()

type JwtCustomClaims struct {
	Username   string `json:"username"`
	UserID     uint   `json:"user_id"`
	Group      string `json:"group"`
	Permission string `json:"permission"`
	jwt.StandardClaims
}

var IsAuth = middleware.JWTWithConfig(middleware.JWTConfig{
	Claims:     &JwtCustomClaims{},
	SigningKey: []byte(config.JWT_SECREET_KEY),
})
