package models

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

var EMPTY_JWT_CLAIMS = JWTClaims{}
