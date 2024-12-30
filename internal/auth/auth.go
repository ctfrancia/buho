package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Auth struct {
	secretKey []byte
}

func NewAuth(key string) *Auth {
	return &Auth{
		secretKey: []byte(key),
	}
}

// CreateJWT generates a new JWT token
func (a *Auth) CreateJWT(userID string) (string, error) {
	// Define the claims
	claims := jwt.MapClaims{
		"sub": userID,                                // User ID (subject)
		"iat": time.Now().Unix(),                     // Issued at time
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expiration time (1 day from now)
	}

	// Create a new token using the HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		return "", fmt.Errorf("could not sign the token: %w", err)
	}

	return tokenString, nil
}
