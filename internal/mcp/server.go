package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/firasmosbehi/ssh-helper/internal/config"
	"github.com/firasmosbehi/ssh-helper/internal/core"
	"github.com/firasmosbehi/ssh-helper/internal/rsync"
	"github.com/firasmosbehi/ssh-helper/internal/ssh"
	"github.com/firasmosbehi/ssh-helper/internal/store"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Serve starts the MCP server over stdio.
func Serve(cfg *config.Config) error {
	s := server.NewMCPServer("ssh-helper", "1.0.0")

	s.AddTool(mcp.NewTool("list_hosts",
		mcp.WithDescription("List SSH hosts from ~/.ssh/config"),
	), handleListHosts(cfg))

	s.AddTool(mcp.NewTool("show_host",
		mcp.WithDescription("Show details for a specific host"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Host alias")),
	), handleShowHost(cfg))

	s.AddTool(mcp.NewTool("connect_host",
		mcp.WithDescription("Test connectivity to a host"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Host alias")),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to proceed")),
	), handleConnectHost(cfg))

	s.AddTool(mcp.NewTool("run_command",
		mcp.WithDescription("Run a non-interactive command on a remote host"),
		mcp.WithString("host", mcp.Required(), mcp.Description("Host alias")),
		mcp.WithString("command", mcp.Required(), mcp.Description("Shell command to run")),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to proceed")),
	), handleRunCommand(cfg))

	s.AddTool(mcp.NewTool("transfer_files",
		mcp.WithDescription("Run an rsync transfer"),
		mcp.WithString("source", mcp.Required()),
		mcp.WithString("dest", mcp.Required()),
		mcp.WithBoolean("dry_run", mcp.Description("Perform a trial run")),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to proceed")),
	), handleTransferFiles(cfg))

	s.AddTool(mcp.NewTool("list_keys",
		mcp.WithDescription("List SSH key pairs"),
	), handleListKeys())

	s.AddTool(mcp.NewTool("generate_key",
		mcp.WithDescription("Generate a new SSH key pair"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Key file name")),
		mcp.WithString("type", mcp.DefaultString("ed25519"), mcp.Description("Key type: ed25519 or rsa")),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to proceed")),
	), handleGenerateKey())

	s.AddTool(mcp.NewTool("list_sync_jobs",
		mcp.WithDescription("List persisted rsync jobs"),
	), handleListSyncJobs())

	s.AddTool(mcp.NewTool("run_sync_job",
		mcp.WithDescription("Run a persisted rsync job"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Job ID")),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to proceed")),
	), handleRunSyncJob())

	return server.ServeStdio(s)
}

func handleListHosts(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := ssh.ParseConfig(cfg.SSHConfigPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("parse config: %v", err)), nil
		}
		hosts := ssh.HostsFromConfig(c)
		b, _ := json.MarshalIndent(hosts, "", "  ")
		return mcp.NewToolResultText(string(b)), nil
	}
}

func handleShowHost(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := req.GetString("name", "")
		c, err := ssh.ParseConfig(cfg.SSHConfigPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("parse config: %v", err)), nil
		}
		h, ok := ssh.GetHost(c, name)
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("host %q not found", name)), nil
		}
		b, _ := json.MarshalIndent(h, "", "  ")
		return mcp.NewToolResultText(string(b)), nil
	}
}

func handleConnectHost(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if !req.GetBool("confirm", false) {
			return mcp.NewToolResultError("confirmation required: set confirm=true"), nil
		}
		name := req.GetString("name", "")
		opts, err := resolveHost(cfg, name)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := ssh.TestConnectivity(opts); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("connection failed: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Successfully connected to %s", name)), nil
	}
}

func handleRunCommand(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if !req.GetBool("confirm", false) {
			return mcp.NewToolResultError("confirmation required: set confirm=true"), nil
		}
		host := req.GetString("host", "")
		cmd := req.GetString("command", "")
		if strings.ContainsAny(cmd, ";|&`$(){}[]\n\r") {
			return mcp.NewToolResultError("command contains invalid characters"), nil
		}
		opts, err := resolveHost(cfg, host)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		out, err := ssh.RunCommand(opts, cmd)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("%v\n%s", err, truncate(out, 2000))), nil
		}
		return mcp.NewToolResultText(truncate(out, 2000)), nil
	}
}

func handleTransferFiles(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if !req.GetBool("confirm", false) {
			return mcp.NewToolResultError("confirmation required: set confirm=true"), nil
		}
		source := req.GetString("source", "")
		dest := req.GetString("dest", "")
		if strings.HasPrefix(source, "-") || strings.HasPrefix(dest, "-") {
			return mcp.NewToolResultError("source/dest cannot start with '-'"), nil
		}
		job := core.SyncJob{
			Source: source,
			Dest:   dest,
			DryRun: req.GetBool("dry_run", false),
		}
		runner := rsync.Runner{Job: job}
		res := runner.Run(ctx, rsync.RunOptions{CaptureLog: true})
		if res.Error != nil {
			return mcp.NewToolResultError(fmt.Sprintf("%v\n%s", res.Error, truncate(res.Log, 2000))), nil
		}
		return mcp.NewToolResultText(truncate(res.Log, 2000)), nil
	}
}

func handleListKeys() server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		home, _ := os.UserHomeDir()
		keys, err := ssh.ListKeys(filepath.Join(home, ".ssh"))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("list keys: %v", err)), nil
		}
		b, _ := json.MarshalIndent(keys, "", "  ")
		return mcp.NewToolResultText(string(b)), nil
	}
}

func handleGenerateKey() server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if !req.GetBool("confirm", false) {
			return mcp.NewToolResultError("confirmation required: set confirm=true"), nil
		}
		name := req.GetString("name", "")
		if strings.ContainsAny(name, "/\\") || strings.Contains(name, "..") {
			return mcp.NewToolResultError("invalid key name"), nil
		}
		keyType := req.GetString("type", "ed25519")
		home, _ := os.UserHomeDir()
		if err := ssh.GenerateKey(filepath.Join(home, ".ssh"), name, keyType); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("generate key: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Generated %s key %s", keyType, name)), nil
	}
}

func handleListSyncJobs() server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		s, err := getStore()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		jobs, err := s.ListSyncJobs()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("list jobs: %v", err)), nil
		}
		b, _ := json.MarshalIndent(jobs, "", "  ")
		return mcp.NewToolResultText(string(b)), nil
	}
}

func handleRunSyncJob() server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if !req.GetBool("confirm", false) {
			return mcp.NewToolResultError("confirmation required: set confirm=true"), nil
		}
		id := req.GetString("id", "")
		s, err := getStore()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		jobs, err := s.ListSyncJobs()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("list jobs: %v", err)), nil
		}
		var job core.SyncJob
		for _, j := range jobs {
			if j.ID == id {
				job = j
				break
			}
		}
		if job.ID == "" {
			return mcp.NewToolResultError(fmt.Sprintf("job %q not found", id)), nil
		}
		runner := rsync.Runner{Job: job}
		res := runner.Run(ctx, rsync.RunOptions{CaptureLog: true})
		if res.Error != nil {
			return mcp.NewToolResultError(fmt.Sprintf("%v\n%s", res.Error, res.Log)), nil
		}
		return mcp.NewToolResultText(res.Log), nil
	}
}

func resolveHost(cfg *config.Config, name string) (ssh.ConnectOptions, error) {
	c, err := ssh.ParseConfig(cfg.SSHConfigPath)
	if err != nil {
		return ssh.ConnectOptions{}, err
	}
	opts := ssh.ConnectOptions{Host: name, Port: 22}
	if h, ok := ssh.GetHost(c, name); ok {
		opts.Host = h.Hostname
		opts.User = h.User
		opts.Port = h.Port
		opts.Identity = h.IdentityFile
	}
	if opts.Port == 0 {
		opts.Port = 22
	}
	return opts, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "... [truncated]"
}

func getStore() (store.Store, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return store.NewJSONStore(filepath.Join(dir, "ssh-helper")), nil
}
