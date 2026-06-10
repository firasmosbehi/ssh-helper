//go:build !cgo

package gui

import "fmt"

// Run returns an error when CGO is not available.
func Run() error {
	return fmt.Errorf("GUI requires CGO; build with CGO_ENABLED=1 or install a platform-specific release")
}
