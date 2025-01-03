package auth

import (
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/big"
	"strings"
	"time"
)

const (
	// TournamentAPIRequesterKey is the key used to store the user ID in the request context in the tournament API handlers
	TournamentAPIRequesterKey = "tournamentAPIUser"
	// PasswordGeneratorDefaultLength is the length of the generated password for API users
	PasswordGeneratorDefaultLength = 16 // TODO: This should be a configuration option
)

type Auth struct {
	secretKey []byte
}

func NewAuth(key string) *Auth {
	return &Auth{
		secretKey: []byte(key),
	}
}

// ValidateJWT validates a JWT token
// TODO: REVIEW THIS
func (a *Auth) ValidateJWT(tokenString string) (string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return a.secretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("could not parse the token: %w", err)
	}
	return token.Claims.(jwt.MapClaims)["sub"].(string), nil
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

// CreateSecretKey creates a new secret key, or password, used for user's credentials
func CreateSecretKey(length int) (string, error) {
	// Define the character sets
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lower := "abcdefghijklmnopqrstuvwxyz"
	digits := "0123456789"
	special := "!@#$%^&*()-_=+[]{}|;:,.<>?/~`"

	// Combine all character sets into one
	allCharacters := upper + lower + digits + special

	var password strings.Builder

	// Generate each character for the password
	for i := 0; i < length; i++ {
		// Get a random index into the combined character set
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(allCharacters))))
		if err != nil {
			return "", fmt.Errorf("Error with rand.Int: %v", err)
		}

		// Append the random character to the password
		password.WriteByte(allCharacters[index.Int64()])
	}

	return password.String(), nil
}
