package platform

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Editor returns the user's preferred editor, with a sensible OS fallback.
func Editor() string {
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	if e := os.Getenv("VISUAL"); e != "" {
		return e
	}
	if runtime.GOOS == "windows" {
		return "notepad"
	}
	return "vi"
}

// OpenInEditor opens the given file in the user's preferred editor and waits.
func OpenInEditor(path string) error {
	cmd := exec.Command(Editor(), path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}
	return nil
}
