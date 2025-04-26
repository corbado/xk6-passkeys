// Package passkeys provides a k6 extension for passkeys load testing.
package passkeys

import (
	"github.com/descope/virtualwebauthn"
	"github.com/google/uuid"
	"go.k6.io/k6/js/modules"
)

const iCloudKeychainAaguid string = "fbfc3007-154e-4ecc-8c0b-6e020557d7bd"

func init() {
	modules.Register("k6/x/passkeys", new(Passkeys))
}

// Passkeys is the main struct for the passkeys module.
type Passkeys struct {
}

// NewCredential creates a new credential.
func (p *Passkeys) NewCredential() virtualwebauthn.Credential {
	return virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2)
}

// NewRelyingParty creates a new relying party.
func (p *Passkeys) NewRelyingParty(name string, id string, origin string) virtualwebauthn.RelyingParty {
	return virtualwebauthn.RelyingParty{Name: name, ID: id, Origin: origin}
}

// CreateAttestationResponse creates an attestation response.
func (p *Passkeys) CreateAttestationResponse(
	rp virtualwebauthn.RelyingParty,
	credential virtualwebauthn.Credential,
	attestationOptions string,
) string {
	aaguid, err := uuid.Parse(iCloudKeychainAaguid)
	if err != nil {
		panic(err)
	}

	authenticator := virtualwebauthn.NewAuthenticatorWithOptions(virtualwebauthn.AuthenticatorOptions{
		BackupEligible: true,
		BackupState:    true,
	})
	authenticator.Aaguid = [16]byte(aaguid)

	parsedAttestationOptions, err := virtualwebauthn.ParseAttestationOptions(attestationOptions)
	if err != nil {
		panic(err)
	}

	return virtualwebauthn.CreateAttestationResponse(rp, authenticator, credential, *parsedAttestationOptions)
}

// CreateAssertionResponse creates an assertion response.
func (p *Passkeys) CreateAssertionResponse(
	rp virtualwebauthn.RelyingParty,
	credential virtualwebauthn.Credential,
	userHandle string,
	assertionOptions string,
) string {
	aaguid, err := uuid.Parse(iCloudKeychainAaguid)
	if err != nil {
		panic(err)
	}

	authenticator := virtualwebauthn.NewAuthenticatorWithOptions(virtualwebauthn.AuthenticatorOptions{
		UserHandle:     []byte(userHandle),
		BackupEligible: true,
		BackupState:    true,
	})
	authenticator.Aaguid = [16]byte(aaguid)

	parsedAssertionOptions, err := virtualwebauthn.ParseAssertionOptions(assertionOptions)
	if err != nil {
		panic(err)
	}

	return virtualwebauthn.CreateAssertionResponse(rp, authenticator, credential, *parsedAssertionOptions)
}
