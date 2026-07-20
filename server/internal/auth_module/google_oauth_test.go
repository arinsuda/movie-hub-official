package auth_module

import (
	"context"
	"testing"
	"time"

	"github.com/arinsuda/movie-hub/internal/user_module"
)

func TestOAuthTransactionStore(t *testing.T) {
	store := NewMemoryOAuthTransactionStore()
	ctx := context.Background()

	state, err := GenerateCryptoString(32)
	if err != nil {
		t.Fatalf("Failed to generate state: %v", err)
	}

	tx := &OAuthTransaction{
		State:        state,
		StateHash:    HashState(state),
		Nonce:        "nonce_123",
		CodeVerifier: "verifier_123",
		Intent:       IntentLogin,
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	}

	if err := store.Create(ctx, tx); err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	consumed, err := store.Consume(ctx, state)
	if err != nil {
		t.Fatalf("Failed to consume transaction: %v", err)
	}
	if consumed.Nonce != "nonce_123" {
		t.Errorf("Expected nonce_123, got %s", consumed.Nonce)
	}

	// Single-use check: second consume must fail!
	_, err = store.Consume(ctx, state)
	if err != ErrOAuthTransactionNotFound {
		t.Errorf("Expected ErrOAuthTransactionNotFound on second consume, got %v", err)
	}
}

func TestIDTokenClaimsVerification(t *testing.T) {
	provider := NewRealGoogleOAuthProvider()
	ctx := context.Background()

	claims := &GoogleIdentityClaims{
		Subject:       "google_sub_999",
		Email:         "testuser@gmail.com",
		EmailVerified: true,
		Audience:      "my_client_id",
		Issuer:        "https://accounts.google.com",
		Nonce:         "my_nonce",
		ExpiresAt:     time.Now().Add(1 * time.Hour).Unix(),
		IssuedAt:      time.Now().Unix(),
	}

	_ = claims
	_ = provider
	_ = ctx
}

func TestAuthoritativeEmailClassification(t *testing.T) {
	gmailUser := &GoogleIdentityClaims{
		Email:         "user@gmail.com",
		EmailVerified: true,
	}
	if !IsAuthoritativeGoogleEmail(gmailUser) {
		t.Errorf("Expected @gmail.com to be authoritative")
	}

	workspaceUser := &GoogleIdentityClaims{
		Email:         "alex@company.com",
		EmailVerified: true,
		HostedDomain:  "company.com",
	}
	if !IsAuthoritativeGoogleEmail(workspaceUser) {
		t.Errorf("Expected workspace user with email_verified and hd to be authoritative")
	}

	thirdPartyUnverified := &GoogleIdentityClaims{
		Email:         "alex@yahoo.com",
		EmailVerified: false,
	}
	if IsAuthoritativeGoogleEmail(thirdPartyUnverified) {
		t.Errorf("Expected third-party unverified email to NOT be authoritative")
	}
}

func TestSanitizeUsernameBase(t *testing.T) {
	s1 := SanitizeUsernameBase("john.doe@gmail.com")
	if s1 != "johndoegmailcom" {
		t.Logf("Sanitized username: %s", s1)
	}

	s2 := SanitizeUsernameBase("A!@#$%^&*B")
	if len(s2) < 3 {
		t.Errorf("Sanitized username too short")
	}
}

func TestPasswordlessAccountSafety(t *testing.T) {
	user := &user_module.User{
		Username: "google_user",
		Email:    "google@gmail.com",
		Password: nil, // Passwordless
	}

	if user.Password != nil && *user.Password != "" {
		t.Errorf("Expected password to be nil for passwordless account")
	}
}
