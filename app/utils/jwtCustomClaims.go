package utils

import (
	"github.com/dgrijalva/jwt-go"
)

// MySigningKey key for jwt secret.
const MySigningKey string = "MySecret1234"

// MyCustomClaims for claim jwt payload.
type MyCustomClaims struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	ID       string `json:"id,omitempty"`
	jwt.StandardClaims
}
