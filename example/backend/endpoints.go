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

	if user == nil {
		user = NewUser(username, username)
		s.userDB.PutUser(user)
	}

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
	}

	options, sessionData, err := s.webAuthn.BeginRegistration(user, registerOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.sessionDB.SaveSession("registration_"+username, sessionData)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	sessionData, err := s.sessionDB.GetSession("registration_" + username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if sessionData == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session not found"})
		return
	}

	credential, err := s.webAuthn.FinishRegistration(user, *sessionData, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.AddCredential(*credential)
	s.userDB.PutUser(user)
	s.sessionDB.DeleteSession("registration_" + username)

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	options, sessionData, err := s.webAuthn.BeginLogin(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.sessionDB.SaveSession("authentication_"+username, sessionData)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	sessionData, err := s.sessionDB.GetSession("authentication_" + username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if sessionData == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session not found"})
		return
	}

	_, err = s.webAuthn.FinishLogin(user, *sessionData, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.sessionDB.DeleteSession("authentication_" + username)
	c.JSON(http.StatusOK, gin.H{"status": "Login Success"})
}
