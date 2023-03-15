package ghcr

import (
	"errors"

	"github.com/cli/go-gh/pkg/auth"
	"github.com/cli/go-gh/pkg/config"
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/docker/docker-credential-helpers/registryurl"
)

// Ghcr handles GitHub Container Registry auth with a local PAT
// Implements the helper interface: https://github.com/docker/docker-credential-helpers/blob/master/credentials/helper.go
type Ghcr struct{}

var errNotImplemented = errors.New("not implemented")

// Add adds new credentials to the keychain.
// TODO: Not implemented. Should be able to just run login flow?
func (ghcr Ghcr) Add(*credentials.Credentials) error {
	return errNotImplemented
}

// Delete removes credentials from the keychain.
// WONTFIX: Not implemented. Should be handled by `gh` CLI
func (ghcr Ghcr) Delete(serverURL string) error {
	return errNotImplemented
}

// Get returns the username and secret to use for a given registry server URL.
func (ghcr Ghcr) Get(serverURL string) (string, string, error) {
	url, err := registryurl.Parse(serverURL)
	if err != nil {
		return "", "", err
	}
	hostname := registryurl.GetHostname(url)
	token, _ := auth.TokenForHost(hostname)

	config, err := config.Read()
	if err != nil {
		return "", "", err
	}

	username, err := config.Get([]string{hostname, "user"})
	if err != nil {
		return "", "", err
	}

	return username, token, nil
}

// List returns the stored URLs and corresponding usernames.
func (ghcr Ghcr) List() (map[string]string, error) {
	config, err := config.Read()
	if err != nil {
		return map[string]string{}, err
	}
	hosts, _ := config.Keys([]string{"hosts"})
	auths := map[string]string{} // maps hostname -> username
	for _, host := range hosts {
		username, err := config.Get([]string{"hosts", host, "user"})
		if err == nil {
			auths[host] = username
		}
	}

	return auths, nil
}
