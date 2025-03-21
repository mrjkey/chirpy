package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordHashing(t *testing.T) {
	password := "test"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Hashing failed")
		t.Fail()
	}
	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("Compare failed")
		t.Fail()
	}
	// sanity check
	err = CheckPasswordHash("no", hash)
	if err == nil {
		t.Errorf("Compare succeeded when it should have failed")
		t.Fail()
	}
}

func TestJWTs(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "i'm a secret token"
	expiresIn := time.Second * 2

	signedString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("failed to make jwt: %v", err.Error())
	}

	retrievedUUID, err := ValidateJWT(signedString, tokenSecret)
	if err != nil {
		t.Fatalf("failed to validate jwt: %v", err.Error())
	}

	if retrievedUUID != userID {
		t.Fatal("uuid does not match")
	}

	time.Sleep(expiresIn)

	_, err = ValidateJWT(signedString, tokenSecret)
	if err == nil {
		t.Fatalf("token should have expired by now")
	}
}

func TestGetBearerToken(t *testing.T) {
	token := "MyToken"
	tokenString := fmt.Sprintf("Bearer %v", token)
	headers := http.Header{}
	headers.Set("Authorization", tokenString)
	returnedString, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("error getting token: %v", err.Error())
	}
	if returnedString != token {
		t.Fatal("token does not match")
	}
}

func TestMakeRefreshToken(t *testing.T) {
	token, err := MakeRefreshToken()
	if err != nil {
		t.Fatal("it failed somehow")
	}
	if token == "" {
		t.Fatalf("The token is empty: %v", token)
	}
}
