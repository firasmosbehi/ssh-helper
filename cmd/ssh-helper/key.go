package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/firasmosbehi/ssh-helper/internal/ssh"
	"github.com/spf13/cobra"
)

func newKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "Manage SSH keys",
	}
	cmd.AddCommand(
		newKeyListCommand(),
		newKeyGenerateCommand(),
		newKeyRemoveCommand(),
		newKeyCopyIDCommand(),
		newKeyAddToAgentCommand(),
	)
	return cmd
}

func sshDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ssh")
}

func newKeyListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List SSH key pairs",
		RunE: func(cmd *cobra.Command, args []string) error {
			keys, err := ssh.ListKeys(sshDir())
			if err != nil {
				return err
			}
			for _, k := range keys {
				fmt.Printf("%s\t%s\t%s\n", k.Name, k.Type, k.Fingerprint)
			}
			return nil
		},
	}
}

func newKeyGenerateCommand() *cobra.Command {
	var keyType string
	cmd := &cobra.Command{
		Use:   "generate <name>",
		Short: "Generate a new SSH key pair",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return ssh.GenerateKey(sshDir(), args[0], keyType)
		},
	}
	cmd.Flags().StringVar(&keyType, "type", "ed25519", "key type (ed25519 or rsa)")
	return cmd
}

func newKeyRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove an SSH key pair",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return ssh.RemoveKey(filepath.Join(sshDir(), args[0]))
		},
	}
}

func newKeyCopyIDCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "copy-id <key-name> <host>",
		Short: "Copy a public key to a remote host",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			pubPath := filepath.Join(sshDir(), args[0]+".pub")
			// Resolve host from config or raw
			cfg, _ := ssh.ParseConfig(appConfig.SSHConfigPath)
			opts := ssh.ConnectOptions{Host: args[1], Port: 22}
			if h, ok := ssh.GetHost(cfg, args[1]); ok {
				opts.Host = h.Hostname
				opts.User = h.User
				opts.Port = h.Port
				opts.Identity = h.IdentityFile
			}
			if opts.User == "" {
				if u, err := user.Current(); err == nil {
					opts.User = u.Username
				}
			}
			if opts.Port == 0 {
				opts.Port = 22
			}
			client, err := ssh.DialClient(opts)
			if err != nil {
				return err
			}
			defer client.Close()
			return ssh.CopyID(client, pubPath)
		},
	}
}

func newKeyAddToAgentCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "add-to-agent <name>",
		Short: "Add a private key to the SSH agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return ssh.AddKeyToAgent(filepath.Join(sshDir(), args[0]))
		},
	}
}
