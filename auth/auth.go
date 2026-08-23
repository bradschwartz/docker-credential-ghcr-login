package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cli/cli/v2/pkg/cmd/auth/shared"
	"github.com/cli/go-gh/pkg/auth"
	ghconfig "github.com/cli/go-gh/pkg/config"
	"github.com/cli/oauth"
)

const ghClientID = "178c6fc778ccc68e1d6a"
const ghClientSecret = "34ddeff2b558a23d38fba8a6de74f086ede1cc0b"

// EnsureValidTokenForHost will check for an existing token. If one is not found,
// or found without the required scopes (`write:packages`), a new login flow will
// run to get a new one.
func EnsureValidTokenForHost(hostname string, requiredScopes string, version string) error {
	token, tokenSource := auth.TokenForHost(hostname)
	if !hasRequiredScopes(hostname, requiredScopes, token) {
		log.Printf(
			"Token found did not have required scopes: %s. Source: %s\n",
			requiredScopes,
			tokenSource)
		err := loginFlow(hostname, requiredScopes)
		if err != nil {
			return err
		}
	}
	return nil
}

func hasRequiredScopes(hostname string, requiredScopes string, token string) bool {
	scopes, _ := shared.GetScopes(http.DefaultClient, hostname, token)
	return strings.Contains(scopes, requiredScopes)
}

func loginFlow(hostname string, requiredScopes string) error {
	host, err := oauth.NewGitHubHost("https://" + hostname)
	if err != nil {
		return fmt.Errorf("error setting up OAuth host: %w", err)
	}

	flow := &oauth.Flow{
		Host:         host,
		ClientID:     ghClientID,
		ClientSecret: ghClientSecret,
		Scopes:       []string{requiredScopes},
		Stdin:        os.Stdin,
		Stdout:       os.Stderr,
	}

	token, err := flow.DetectFlow()
	if err != nil {
		return fmt.Errorf("OAuth flow failed: %w", err)
	}

	username, err := shared.GetCurrentLogin(http.DefaultClient, hostname, token.Token)
	if err != nil {
		return fmt.Errorf("error getting username: %w", err)
	}

	cfg, err := ghconfig.Read()
	if err != nil {
		return fmt.Errorf("error reading gh config: %w", err)
	}

	cfg.Set([]string{"hosts", hostname, "user"}, username)
	cfg.Set([]string{"hosts", hostname, "oauth_token"}, token.Token)

	if err := ghconfig.Write(cfg); err != nil {
		return fmt.Errorf("error writing gh config: %w", err)
	}

	return nil
}
