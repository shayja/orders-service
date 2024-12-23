package utils

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

// GenerateJWT generates a signed JWT token containing the user ID and expiration time.
func GenerateJWT(userId string, secretKey string) (string, error) {
	// Define the claims (the payload of the JWT)
	claims := jwt.MapClaims{
		"sub": userId,       // user ID (subject)
		"exp": time.Now().Add(time.Hour * 24).Unix(), // expiration time (1 day)
		"iat": time.Now().Unix(), // issued at time
	}

	// Create the token with the claims and sign it using the HMAC method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	return tokenString, nil
}