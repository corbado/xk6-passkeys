package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// WebAuthnServer handles WebAuthn registration and authentication endpoints.
type WebAuthnServer struct {
	webAuthn  *webauthn.WebAuthn
	userDB    *UserDB
	sessionDB *SessionDB
}

// NewWebAuthnServer creates a new WebAuthn server with the given dependencies.
func NewWebAuthnServer(webAuthn *webauthn.WebAuthn, userDB *UserDB, sessionDB *SessionDB) *WebAuthnServer {
	return &WebAuthnServer{
		webAuthn:  webAuthn,
		userDB:    userDB,
		sessionDB: sessionDB,
	}
}

// RegisterStart initiates the WebAuthn registration process for a user.
func (s *WebAuthnServer) RegisterStart(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	user, err := s.userDB.GetUser(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}

	user = NewUser(username, username)
	s.userDB.PutUser(user)

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
	}

	options, sessionData, err := s.webAuthn.BeginRegistration(user, registerOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	challenge := options.Response.Challenge.String()

	s.sessionDB.SaveSession("register_"+challenge, sessionData)
	c.JSON(http.StatusOK, options)
}

// RegisterFinish completes the WebAuthn registration process for a user.
func (s *WebAuthnServer) RegisterFinish(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	user, err := s.userDB.GetUser(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	credentialCreationResponse, err := protocol.ParseCredentialCreationResponse(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	challenge := credentialCreationResponse.Response.CollectedClientData.Challenge

	sessionData, err := s.sessionDB.GetSession("register_" + challenge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if sessionData == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	// Using CreateCredential() instead of FinishRegistration() because we already have the
	// credentialCreationResponse. We had to to this to fix problems with "double consuming"
	// the request body.
	credential, err := s.webAuthn.CreateCredential(user, *sessionData, credentialCreationResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.AddCredential(*credential)
	s.userDB.PutUser(user)
	s.sessionDB.DeleteSession("register_" + challenge)

	c.JSON(http.StatusOK, gin.H{"status": "Registration Success"})
}

// LoginStart initiates the WebAuthn authentication process for a user.
func (s *WebAuthnServer) LoginStart(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	user, err := s.userDB.GetUser(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	options, sessionData, err := s.webAuthn.BeginLogin(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	challenge := options.Response.Challenge.String()
	s.sessionDB.SaveSession("login_"+challenge, sessionData)
	c.JSON(http.StatusOK, options)
}

// LoginFinish completes the WebAuthn authentication process for a user.
func (s *WebAuthnServer) LoginFinish(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	user, err := s.userDB.GetUser(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	credentialRequestResponse, err := protocol.ParseCredentialRequestResponse(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	challenge := credentialRequestResponse.Response.CollectedClientData.Challenge

	sessionData, err := s.sessionDB.GetSession("login_" + challenge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if sessionData == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	// Using ValidateLogin() instead of FinishLogin() because we already have the
	// credentialRequestResponse. We had to to this to fix problems with "double consuming"
	// the request body.
	_, err = s.webAuthn.ValidateLogin(user, *sessionData, credentialRequestResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.sessionDB.DeleteSession("login_" + challenge)
	c.JSON(http.StatusOK, gin.H{"status": "Login Success"})
}
