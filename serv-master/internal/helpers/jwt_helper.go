package helpers

import "github.com/golang-jwt/jwt/v5"

// JWTClaims represents the claims issued by the api-gateway's auth token.
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
