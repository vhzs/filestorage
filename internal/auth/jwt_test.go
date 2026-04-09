package auth

import "testing"

func TestTokenRoundtrip(t *testing.T) {
	secret := "test-secret-123"
	token, err := GenerateToken(42, secret)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	t.Logf("token: %s", token) // чтобы видеть в -v

	uid, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if uid != 42 {
		t.Errorf("uid = %d, want 42", uid)
	}
}

func TestTokenWrongSecret(t *testing.T) {
	token, _ := GenerateToken(1, "secret-a")
	_, err := ParseToken(token, "secret-b")
	if err == nil {
		t.Fatal("expected error with wrong secret")
	}
}

func TestTokenGarbage(t *testing.T) {
	_, err := ParseToken("not.a.real.token", "whatever")
	if err == nil {
		t.Fatal("expected error for garbage input")
	}
}
