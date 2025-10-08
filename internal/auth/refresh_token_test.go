package auth

import (
	"testing"
)

func TestMakeRefreshToken(t *testing.T) {
	tok := MakeRefreshToken()
	if tok == "" {
		t.Fatalf("expected non-empty refresh token")
	}

	// token should be hex-encoded 32 bytes => 64 chars
	if len(tok) < 64 {
		t.Fatalf("unexpected refresh token length: got %d, want >=64", len(tok))
	}
}
