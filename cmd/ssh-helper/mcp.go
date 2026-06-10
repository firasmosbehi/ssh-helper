package main

import (
	"context"
	"fmt"

	"github.com/firasmosbehi/ssh-helper/internal/mcp"
	"github.com/spf13/cobra"
)

func newMCPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "MCP server and client commands",
	}
	cmd.AddCommand(newMCPServeCommand())
	cmd.AddCommand(newMCPServersCommand())
	cmd.AddCommand(newMCPCallCommand())
	return cmd
}

func newMCPServeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server over stdio",
		RunE: func(cmd *cobra.Command, args []string) error {
			return mcp.Serve(appConfig)
		},
	}
}

func newMCPServersCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "servers",
		Short: "List configured external MCP servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			if appConfig == nil {
				return fmt.Errorf("config not loaded")
			}
			for name, srv := range appConfig.MCPClients {
				fmt.Printf("%s\t%s %v\n", name, srv.Command, srv.Args)
			}
			return nil
		},
	}
}

func newMCPCallCommand() *cobra.Command {
	var argsJSON string
	c := &cobra.Command{
		Use:   "call <server> <tool>",
		Short: "Call a tool on an external MCP server",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if appConfig == nil {
				return fmt.Errorf("config not loaded")
			}
			srv, ok := appConfig.MCPClients[args[0]]
			if !ok {
				return fmt.Errorf("server %q not configured", args[0])
			}
			return mcp.CallTool(context.Background(), srv, args[1], argsJSON)
		},
	}
	c.Flags().StringVar(&argsJSON, "args", "{}", "JSON arguments")
	return c
}
