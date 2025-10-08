package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestMakeAndValidateJWT(t *testing.T) {
	secret := "test-secret-key"
	uid := uuid.New()

	token, err := MakeJWT(uid, secret, 1*time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	if token == "" {
		t.Fatalf("expected token to be non-empty")
	}

	got, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}

	if got != uid {
		t.Fatalf("expected validated uid %v, got %v", uid, got)
	}
}
