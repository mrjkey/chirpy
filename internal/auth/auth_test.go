package auth

import (
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
