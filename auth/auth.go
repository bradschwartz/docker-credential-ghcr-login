package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/cli/cli/v2/pkg/cmd/auth/login"
	"github.com/cli/cli/v2/pkg/cmd/auth/shared"
	"github.com/cli/cli/v2/pkg/cmd/factory"
	"github.com/cli/go-gh/pkg/auth"
)

// EnsureValidTokenForHost wil check for an existing token. If one is not found,
// or found without the required scopes (`write:packages`), a new login flow will
// run to get a new one
func EnsureValidTokenForHost(hostname string, requiredScopes string, version string) error {
	token, tokenSource := auth.TokenForHost(hostname)
	if !hasRequiredScopes(hostname, requiredScopes, token) {
		log.Printf("Token found did not have required scopes. Source: %s\n", tokenSource)
		err := loginFlow(hostname, requiredScopes, version)
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

func loginFlow(hostname string, requiredScopes string, version string) error {
	cmdFactory := factory.New(version)

	loginCmd := login.NewCmdLogin(cmdFactory, nil)
	// We do this to not have Cobra redefine the help command and error out
	loginCmd.PersistentFlags().BoolP("help", "", false, "")

	flag := loginCmd.Flag("scopes")
	flag.Value.Set(requiredScopes)

	flag = loginCmd.Flag("hostname")
	flag.Value.Set(hostname)

	return loginCmd.Execute()
}
