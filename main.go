package main

import (
	"log"

	"github.com/bradschwartz/docker-credential-ghcr-login/auth"
)

var (
	version        string
	requiredScopes = "write:packages"
	hostname       = "github.com"
)

func main() {
	err := auth.EnsureValidTokenForHost(hostname, requiredScopes, version)
	if err != nil {
		log.Fatalf("There was an error while logging in: %s", err)
	}

	// Have a token with required scopes!
}
