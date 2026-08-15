package cli

import (
	"errors"
	"os"

	"github.com/zalando/go-keyring"
)

const keyringService = "bloberry"

// getToken resolves the active token: BLOBERRY_TOKEN (CI) wins, then the OS
// keychain. The keychain stores the refresh token; the access token is derived
// from it on demand by auth refresh.
func getToken() (string, bool, error) {
	if t := os.Getenv("BLOBERRY_TOKEN"); t != "" {
		return t, true, nil
	}
	t, err := keyring.Get(keyringService, "refresh_token")
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", false, nil
		}
		return "", false, err
	}
	return t, false, nil
}

func setRefreshToken(token string) error {
	return keyring.Set(keyringService, "refresh_token", token)
}

func deleteRefreshToken() error {
	return keyring.Delete(keyringService, "refresh_token")
}
