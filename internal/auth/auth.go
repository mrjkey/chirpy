package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 1)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
		Subject:   userID.String(),
	})
	signedString, err := token.SignedString([]byte(tokenSecret))
	return signedString, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Return the key used to sign the token
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("invalid token claims")
	}

	userIDStr := claims.Subject
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid user ID in token %v", err)
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return "", fmt.Errorf("token string not found in header")
	}
	stripped := strings.Replace(token, "Bearer ", "", 1)
	return stripped, nil
}

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// fmt.Println(i)
	// fmt.Println(b)

	output := hex.EncodeToString(b)
	return output, nil
}

func AuthorizeUser(headers http.Header, tokenSecret string) (uuid.UUID, error) {
	tokenString, err := GetBearerToken(headers)
	if err != nil {
		return uuid.UUID{}, err
	}

	userID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return "", fmt.Errorf("token string not found in header")
	}
	stripped := strings.Replace(token, "ApiKey ", "", 1)
	return stripped, nil
}

func AuthorizeApiKey(headers http.Header, key string) error {
	apiKey, err := GetAPIKey(headers)
	if err != nil {
		return err
	}

	if apiKey != key {
		return fmt.Errorf("api key does not match")
	}
	return nil
}
