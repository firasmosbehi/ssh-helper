package main

import (
	"fmt"

	"github.com/firasmosbehi/ssh-helper/internal/config"
	"github.com/firasmosbehi/ssh-helper/internal/platform"
	"github.com/spf13/cobra"
)

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage ssh-helper configuration",
	}

	cmd.AddCommand(
		newConfigInitCommand(),
		newConfigGetCommand(),
		newConfigSetCommand(),
		newConfigEditCommand(),
	)

	return cmd
}

func newConfigInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a default configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := config.NewManager()
			if err != nil {
				return err
			}
			if err := mgr.Init(); err != nil {
				return err
			}
			fmt.Printf("Config initialized at %s\n", mgr.Path())
			return nil
		},
	}
}

func newConfigGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfgMgr == nil {
				return fmt.Errorf("config not loaded")
			}
			out, err := cfgMgr.Render(args[0])
			if err != nil {
				return err
			}
			fmt.Println(out)
			return nil
		},
	}
}

func newConfigSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfgMgr == nil {
				return fmt.Errorf("config not loaded")
			}
			if err := cfgMgr.Set(args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("Set %s = %s\n", args[0], args[1])
			return nil
		},
	}
}

func newConfigEditCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Open the configuration file in your default editor",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfgMgr == nil {
				return fmt.Errorf("config not loaded")
			}
			if err := platform.OpenInEditor(cfgMgr.Path()); err != nil {
				return err
			}
			return nil
		},
	}
}
