package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bradschwartz/docker-credential-ghcr-login/auth"
	"github.com/bradschwartz/docker-credential-ghcr-login/ghcr"
	"github.com/docker/docker-credential-helpers/credentials"
)

var (
	version        string
	requiredScopes = "write:packages"
	hostname       = "github.com"
)

func main() {
	var versionFlag bool
	flag.BoolVar(&versionFlag, "v", false, "print version and exit")
	flag.Parse()

	if versionFlag {
		if version == "" {
			fmt.Println("dev")
			os.Exit(0)
		}

		fmt.Println(version)
		os.Exit(0)
	}

	err := auth.EnsureValidTokenForHost(hostname, requiredScopes, version)
	if err != nil {
		log.Fatalf("There was an error while logging in: %s", err)
	}

	// Have a token with required scopes!
	credentials.Serve(ghcr.Ghcr{})
}
