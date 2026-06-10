package main

import (
	"github.com/firasmosbehi/ssh-helper/internal/gui"
	"github.com/spf13/cobra"
)

func newGUICommand() *cobra.Command {
	return &cobra.Command{
		Use:   "gui",
		Short: "Launch the native desktop GUI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return gui.Run()
		},
	}
}
