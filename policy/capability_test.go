package policy_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	"github.com/aegistrace/aegistrace/policy"
)

func TestCapabilityScopeAndExpiry(t *testing.T) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 31, 2, 0, 0, 0, time.UTC)
	token, err := policy.Issue(private, "scheduler", "worker-1", []string{"github:read", "ticket:*"}, time.Minute, now)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := policy.Verify(token, public, "ticket:create", now.Add(30*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if claims.Subject != "worker-1" {
		t.Fatal("wrong subject")
	}
	if _, err := policy.Verify(token, public, "github:write", now); err == nil {
		t.Fatal("ungranted scope accepted")
	}
	if _, err := policy.Verify(token, public, "ticket:create", now.Add(time.Minute)); err == nil {
		t.Fatal("expired token accepted")
	}
}
