package main

import (
	"os/user"
	"strconv"
	"strings"

	"github.com/firasmosbehi/ssh-helper/internal/ssh"
	"github.com/spf13/cobra"
)

func newConnectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "connect <host>",
		Short: "Connect to an SSH host",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg := args[0]
			opts := ssh.ConnectOptions{}

			// Try to look up as a named host first
			cfg, err := ssh.ParseConfig(appConfig.SSHConfigPath)
			if err == nil {
				if h, ok := ssh.GetHost(cfg, arg); ok {
					opts.Host = h.Hostname
					opts.User = h.User
					opts.Port = h.Port
					opts.Identity = h.IdentityFile
				}
			}

			// If not found, parse user@host:port syntax
			if opts.Host == "" {
				opts.Host = arg
				if strings.Contains(arg, "@") {
					parts := strings.SplitN(arg, "@", 2)
					opts.User = parts[0]
					opts.Host = parts[1]
				}
				if strings.Contains(opts.Host, ":") {
					parts := strings.SplitN(opts.Host, ":", 2)
					opts.Host = parts[0]
					if p, err := strconv.Atoi(parts[1]); err == nil {
						opts.Port = p
					}
				}
			}

			if opts.User == "" {
				u, _ := user.Current()
				if u != nil {
					opts.User = u.Username
				}
			}
			if opts.Port == 0 {
				opts.Port = 22
			}

			return ssh.Connect(opts)
		},
	}
}
