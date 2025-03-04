package auth

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ctfrancia/buho/internal/model"
	"github.com/ctfrancia/buho/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

type key int

const (
	// TournamentAPIRequesterKey is the key used to store the user ID in the request context in the tournament API handlers
	// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
	TournamentAPIRequesterKey key = iota
	// PasswordGeneratorDefaultLength is the length of the generated password for API users
	PasswordGeneratorDefaultLength = 16 // TODO: This should be a configuration option
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type a2params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

type Auth struct {
	privateKeyPath string
	publicKeyPath  string
}

func NewAuth(privKeyPath, PubKeyPath string) *Auth {
	return &Auth{
		privateKeyPath: privKeyPath,
		publicKeyPath:  PubKeyPath,
	}
}

// ValidateJWT validates a JWT token
func (a *Auth) ValidateJWT(tokenString string) (string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwt.ParseRSAPublicKeyFromPEM([]byte(a.publicKeyPath))
	})
	if err != nil {
		return "", fmt.Errorf("could not parse the token: %w", err)
	}
	return token.Claims.(jwt.MapClaims)["sub"].(string), nil
}

// CreateJWT generates a new JWT token
func (a *Auth) CreateJWT(user repository.Auth) (string, error) {
	// Define the claims
	claims := jwt.MapClaims{
		"sub": map[string]interface{}{
			"id":      user.ID,
			"email":   user.Email,
			"website": user.Website,
		},
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expiration time (1 day from now)
	}

	// Create a new token using the HS256 signing methuod SigningMethodHMAC
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	// ParseED25519PrivateKey reads an ED25519 private key from a PEM file
	key, err := ParseED25519PrivateKey(a.privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("could not parse the private key: %w", err)
	}

	// Sign the token with the secret key
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("could not sign the token: %w", err)
	}

	return tokenString, nil
}

// CreateRefreshToken generates a new refresh token
func (a *Auth) CreateRefreshToken(user *repository.Auth) (string, error) {
	// Generate a new secret key
	secretKey, err := CreateSecretKey(PasswordGeneratorDefaultLength)
	if err != nil {
		return "", fmt.Errorf("could not create a new secret key: %w", err)
	}

	// NOTE: the "rt_" prefix is used to indicate that this is a refresh token
	sk := fmt.Sprintf("rt_%s", secretKey)

	hashedSecretKey, err := Hash(sk)
	if err != nil {
		return "", fmt.Errorf("could not hash the secret key: %w", err)
	}

	user.RefreshToken = hashedSecretKey
	user.RefreshTokenExpiry = time.Now().Add(time.Hour * 24 * 7) // 1 week from now

	return sk, nil
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
			return "", fmt.Errorf("error with rand.Int: %v", err)
		}

		// Append the random character to the password
		password.WriteByte(allCharacters[index.Int64()])
	}

	return password.String(), nil
}

// Hash hashes the password
func Hash(password string) (string, error) {
	// Define the parameters for the Argon2 algorithm
	// FIXME: These should be configuration options
	p := a2params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
	hash, err := generateFromPassword(password, &p)
	if err != nil {
		return "", fmt.Errorf("error generating hash: %v", err)
	}

	return hash, nil
}

func generateFromPassword(password string, p *a2params) (encodedHash string, err error) {
	// Generate a cryptographically secure random salt.
	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// CompareHashAndPassword compares a hashed password with the plain text password
func CompareHashAndPassword(encodedHash, password string) (bool, error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}
	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (p *a2params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &a2params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}

// ParseED25519PublicKey reads an ED25519 public key from a PEM file
func ParseED25519PublicKey(filepath string) (ed25519.PublicKey, error) {
	// Read the key file
	keyBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %v", err)
	}

	// Decode the PEM-encoded key
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Parse the public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	// Type assert to ed25519.PublicKey
	ed25519Key, ok := pub.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ED25519 public key")
	}

	return ed25519Key, nil
}

// ParseED25519PrivateKey reads an ED25519 private key from a PEM file
func ParseED25519PrivateKey(filepath string) (ed25519.PrivateKey, error) {
	// Read the key file
	keyBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}

	// Decode the PEM-encoded key
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Parse the private key
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	// Type assert to ed25519.PrivateKey
	ed25519Key, ok := privateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an ED25519 private key")
	}

	return ed25519Key, nil
}

// VerifyJWTWithED25519 demonstrates verifying a JWT with an ED25519 public key
func VerifyJWTWithED25519(tokenString string, publicKeyPath string) (model.Subject, error) {
	// Parse the public key
	publicKey, err := ParseED25519PublicKey(publicKeyPath)
	if err != nil {
		return model.Subject{}, fmt.Errorf("failed to parse public key: %v", err)
	}

	// Parse and verify the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is EdDSA
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return model.Subject{}, fmt.Errorf("failed to parse token: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return model.Subject{}, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["sub"] == nil {
			return model.Subject{}, fmt.Errorf("invalid token or claims")
		}
		sub := claims["sub"].(map[string]interface{})
		return model.Subject{
			ID:      int(sub["id"].(float64)),
			Email:   sub["email"].(string),
			Website: sub["website"].(string),
		}, nil
	} else {
		return model.Subject{}, fmt.Errorf("invalid token or claims")
	}
}
