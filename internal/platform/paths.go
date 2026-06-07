package platform

import (
	"fmt"
	"os"
	"path/filepath"
)

// ConfigDir returns the directory used for configuration and data files.
func ConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}
	return filepath.Join(dir, "ssh-helper"), nil
}

// DataDir returns the directory used for persistent application data.
func DataDir() (string, error) {
	return ConfigDir()
}

// EnsureDir creates the directory if it does not exist.
func EnsureDir(path string) error {
	if err := os.MkdirAll(path, 0o700); err != nil {
		return fmt.Errorf("create directory %s: %w", path, err)
	}
	return nil
}
