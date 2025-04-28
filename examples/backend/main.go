package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
)

func main() {
	// Initialize WebAuthn
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "WebAuthn Demo",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:8080"},
	})
	if err != nil {
		log.Fatal("failed to create WebAuthn:", err)
	}

	userDB := NewUserDB()
	sessionDB := NewSessionDB()
	server := NewWebAuthnServer(webAuthn, userDB, sessionDB)

	router := gin.Default()
	router.GET("/register/start/:username", server.RegisterStart)
	router.POST("/register/finish/:username", server.RegisterFinish)
	router.GET("/login/start/:username", server.LoginStart)
	router.POST("/login/finish/:username", server.LoginFinish)

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
