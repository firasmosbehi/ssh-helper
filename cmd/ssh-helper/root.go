package main

import (
	"fmt"
	"os"

	"github.com/firasmosbehi/ssh-helper/internal/config"
	"github.com/firasmosbehi/ssh-helper/internal/logger"
	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	verbose   int
	jsonLogs  bool
	cfgMgr    *config.Manager
	appConfig *config.Config
)

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "ssh-helper",
		Short: "CLI + TUI + GUI tool for managing SSH and rsync",
		Long: `ssh-helper helps you manage SSH connections, keys, configs,
and run rsync transfers through a unified CLI, TUI, and native GUI.
It also exposes and consumes MCP servers.`,
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip config loading for config init command.
			if cmd.Name() == "init" && cmd.Parent() != nil && cmd.Parent().Name() == "config" {
				return nil
			}
			mgr, err := config.NewManager()
			if err != nil {
				return fmt.Errorf("initialize config: %w", err)
			}
			cfgMgr = mgr
			cfg, err := cfgMgr.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			appConfig = cfg

			level := cfg.LogLevel
			if verbose > 0 {
				level = "debug"
			}
			format := cfg.LogFormat
			if jsonLogs {
				format = "json"
			}
			if err := logger.Init(logger.Options{Level: level, Format: format}); err != nil {
				return fmt.Errorf("initialize logger: %w", err)
			}
			return nil
		},
	}

	root.PersistentFlags().CountVarP(&verbose, "verbose", "v", "increase verbosity (use -vv for debug)")
	root.PersistentFlags().BoolVar(&jsonLogs, "json", false, "output logs as JSON")

	root.AddCommand(newConfigCommand())
	root.AddCommand(newHostCommand())
	root.AddCommand(newConnectCommand())
	root.AddCommand(newKeyCommand())
	root.AddCommand(newSyncCommand())
	root.AddCommand(newTUICommand())

	return root
}

func main() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
