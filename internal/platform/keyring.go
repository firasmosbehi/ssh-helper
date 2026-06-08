package platform

import (
	"github.com/zalando/go-keyring"
)

const keyringService = "ssh-helper"

// SetKeyringSecret stores a secret in the OS keyring.
func SetKeyringSecret(account, secret string) error {
	return keyring.Set(keyringService, account, secret)
}

// GetKeyringSecret retrieves a secret from the OS keyring.
func GetKeyringSecret(account string) (string, error) {
	return keyring.Get(keyringService, account)
}

// DeleteKeyringSecret removes a secret from the OS keyring.
func DeleteKeyringSecret(account string) error {
	return keyring.Delete(keyringService, account)
}

// KeyringSecretExists checks whether a secret exists in the keyring.
func KeyringSecretExists(account string) bool {
	_, err := GetKeyringSecret(account)
	return err == nil
}
