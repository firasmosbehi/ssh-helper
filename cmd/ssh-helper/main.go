package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	root := &cobra.Command{
		Use:   "ssh-helper",
		Short: "CLI + TUI + GUI tool for managing SSH and rsync",
		Long: `ssh-helper helps you manage SSH connections, keys, configs,
and run rsync transfers through a unified CLI, TUI, and native GUI.
It also exposes and consumes MCP servers.`,
		Version: version,
	}

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
