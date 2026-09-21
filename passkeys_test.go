package passkeys

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/descope/virtualwebauthn"
)

func TestCredentialRoundTripAndAssertion(t *testing.T) {
	t.Parallel()
	p := &Passkeys{}
	credential := p.NewCredential()
	encoded, err := json.Marshal(credential)
	if err != nil {
		t.Fatal(err)
	}
	var restored virtualwebauthn.Credential
	if err := json.Unmarshal(encoded, &restored); err != nil {
		t.Fatal(err)
	}
	rp := p.NewRelyingParty("Test", "localhost", "http://localhost:8080")
	response, err := p.CreateAssertionResponse(rp, restored, "user-123",
		`{"publicKey":{"challenge":"dGVzdC1jaGFsbGVuZ2U","rpId":"localhost","userVerification":"required"}}`)
	if err != nil {
		t.Fatal(err)
	}
	var assertion struct {
		Response struct {
			UserHandle        string `json:"userHandle"`
			AuthenticatorData string `json:"authenticatorData"`
			Signature         string `json:"signature"`
		} `json:"response"`
	}
	if err := json.Unmarshal([]byte(response), &assertion); err != nil {
		t.Fatal(err)
	}
	handle, err := base64.RawURLEncoding.DecodeString(assertion.Response.UserHandle)
	if err != nil || string(handle) != "user-123" {
		t.Fatalf("unexpected user handle %q: %v", handle, err)
	}
	data, err := base64.RawURLEncoding.DecodeString(assertion.Response.AuthenticatorData)
	if err != nil || len(data) < 37 {
		t.Fatalf("invalid authenticator data: %v", err)
	}
	// User presence, user verification, backup eligibility and backup state.
	if data[32]&0x1d != 0x1d || assertion.Response.Signature == "" {
		t.Fatal("assertion is missing expected flags or signature")
	}
}

func TestMalformedOptionsReturnErrors(t *testing.T) {
	t.Parallel()
	p := &Passkeys{}
	rp := p.NewRelyingParty("Test", "localhost", "http://localhost:8080")
	credential := p.NewCredential()
	if _, err := p.CreateAttestationResponse(rp, credential, "invalid JSON"); err == nil {
		t.Fatal("expected attestation parsing error")
	}
	if _, err := p.CreateAssertionResponse(rp, credential, "user", "invalid JSON"); err == nil {
		t.Fatal("expected assertion parsing error")
	}
}
