package main

import (
	"crypto/rand"
	"encoding/binary"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// User represents a user in the WebAuthn system.
type User struct {
	id          uint64
	name        string
	displayName string
	credentials []webauthn.Credential
}

// NewUser creates a new user with the given name and display name.
func NewUser(name string, displayName string) *User {
	user := &User{}
	user.id = randomUint64()
	user.name = name
	user.displayName = displayName
	return user
}

// randomUint64 generates a random 64-bit unsigned integer.
func randomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

// WebAuthnID returns the user's ID as a byte slice for WebAuthn operations.
func (u User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.id))
	return buf
}

// WebAuthnName returns the user's username for WebAuthn operations.
func (u User) WebAuthnName() string {
	return u.name
}

// WebAuthnDisplayName returns the user's display name for WebAuthn operations.
func (u User) WebAuthnDisplayName() string {
	return u.displayName
}

// WebAuthnIcon returns an empty string as user icons are not supported.
func (u User) WebAuthnIcon() string {
	return ""
}

// AddCredential adds a new WebAuthn credential to the user's credentials list.
func (u *User) AddCredential(cred webauthn.Credential) {
	u.credentials = append(u.credentials, cred)
}

// WebAuthnCredentials returns all WebAuthn credentials associated with the user.
func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

// CredentialExcludeList returns a list of credential descriptors to prevent duplicate registrations.
func (u User) CredentialExcludeList() []protocol.CredentialDescriptor {
	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}
	return credentialExcludeList
}
