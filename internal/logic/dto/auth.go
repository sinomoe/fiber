package dto

import "github.com/golang-jwt/jwt"

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type AuthClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type ParseTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type ParseTokenResponse = AuthClaims
