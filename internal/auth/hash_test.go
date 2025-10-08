package auth

import (
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	pwd := "s3cr3tP@ssw0rd"

	hash, err := HashPassword(pwd)
	if err != nil {
		t.Fatalf("unexpected error from HashPassword: %v", err)
	}

	if hash == "" {
		t.Fatalf("expected non-empty hash")
	}

	// Correct password should validate
	if err := CheckPasswordHash(pwd, hash); err != nil {
		t.Fatalf("password did not validate: %v", err)
	}

	// Wrong password should fail
	if err := CheckPasswordHash("wrong-password", hash); err == nil {
		t.Fatalf("expected validation to fail for wrong password")
	}
}
