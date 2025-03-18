package auth

import (
	"testing"
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
