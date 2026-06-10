//go:build !cgo

package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newGUICommand() *cobra.Command {
	return &cobra.Command{
		Use:   "gui",
		Short: "Launch the native desktop GUI (disabled in this build)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("GUI requires CGO; build with CGO_ENABLED=1 or install a platform-specific release")
		},
	}
}
