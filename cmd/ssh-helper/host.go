package main

import (
	"fmt"

	"github.com/firasmosbehi/ssh-helper/internal/core"
	"github.com/firasmosbehi/ssh-helper/internal/ssh"
	"github.com/spf13/cobra"
)

func newHostCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host",
		Short: "Manage SSH hosts",
	}
	cmd.AddCommand(
		newHostListCommand(),
		newHostShowCommand(),
		newHostAddCommand(),
		newHostRemoveCommand(),
		newHostEditCommand(),
	)
	return cmd
}

func newHostListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List hosts from ~/.ssh/config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := ssh.ParseConfig(appConfig.SSHConfigPath)
			if err != nil {
				return err
			}
			hosts := ssh.HostsFromConfig(cfg)
			for _, h := range hosts {
				port := h.Port
				if port == 0 {
					port = 22
				}
				fmt.Printf("%s\t%s@%s:%d\n", h.Name, h.User, h.Hostname, port)
			}
			return nil
		},
	}
}

func newHostShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show details for a host",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := ssh.ParseConfig(appConfig.SSHConfigPath)
			if err != nil {
				return err
			}
			host, ok := ssh.GetHost(cfg, args[0])
			if !ok {
				return fmt.Errorf("host %q not found", args[0])
			}
			fmt.Printf("Name:         %s\n", host.Name)
			fmt.Printf("HostName:     %s\n", host.Hostname)
			fmt.Printf("User:         %s\n", host.User)
			fmt.Printf("Port:         %d\n", host.Port)
			fmt.Printf("IdentityFile: %s\n", host.IdentityFile)
			return nil
		},
	}
}

func newHostAddCommand() *cobra.Command {
	var host core.Host
	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Add a host to ~/.ssh/config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			host.Name = args[0]
			if _, err := ssh.BackupConfig(appConfig.SSHConfigPath); err != nil {
				return err
			}
			return ssh.AddHost(appConfig.SSHConfigPath, host)
		},
	}
	cmd.Flags().StringVar(&host.Hostname, "hostname", "", "remote hostname or IP")
	cmd.Flags().StringVar(&host.User, "user", "", "remote user")
	cmd.Flags().IntVar(&host.Port, "port", 0, "remote port")
	cmd.Flags().StringVar(&host.IdentityFile, "identity", "", "private key path")
	return cmd
}

func newHostRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a host from ~/.ssh/config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := ssh.BackupConfig(appConfig.SSHConfigPath); err != nil {
				return err
			}
			return ssh.RemoveHost(appConfig.SSHConfigPath, args[0])
		},
	}
}

func newHostEditCommand() *cobra.Command {
	var hostname, user, identity string
	var port int
	cmd := &cobra.Command{
		Use:   "edit <name>",
		Short: "Edit a host in ~/.ssh/config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := ssh.ParseConfig(appConfig.SSHConfigPath)
			if err != nil {
				return err
			}
			host, ok := ssh.GetHost(cfg, args[0])
			if !ok {
				return fmt.Errorf("host %q not found", args[0])
			}
			if cmd.Flags().Changed("hostname") {
				host.Hostname = hostname
			}
			if cmd.Flags().Changed("user") {
				host.User = user
			}
			if cmd.Flags().Changed("port") {
				host.Port = port
			}
			if cmd.Flags().Changed("identity") {
				host.IdentityFile = identity
			}
			return ssh.EditHost(appConfig.SSHConfigPath, host)
		},
	}
	cmd.Flags().StringVar(&hostname, "hostname", "", "remote hostname or IP")
	cmd.Flags().StringVar(&user, "user", "", "remote user")
	cmd.Flags().IntVar(&port, "port", 0, "remote port")
	cmd.Flags().StringVar(&identity, "identity", "", "private key path")
	return cmd
}
