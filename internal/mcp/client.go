package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/firasmosbehi/ssh-helper/internal/config"
	"github.com/mark3labs/mcp-go/client"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
)

func newStdioClient(srv config.MCPClientConfig) (*client.Client, error) {
	env := make([]string, 0, len(srv.Env))
	for k, v := range srv.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return client.NewStdioMCPClient(srv.Command, env, srv.Args...)
}

func initClient(ctx context.Context, cl *client.Client) error {
	if err := cl.Start(ctx); err != nil {
		return fmt.Errorf("start client: %w", err)
	}
	if _, err := cl.Initialize(ctx, mcpgo.InitializeRequest{
		Params: mcpgo.InitializeParams{
			ProtocolVersion: "2024-11-05",
			ClientInfo: mcpgo.Implementation{
				Name:    "ssh-helper",
				Version: "1.0.0",
			},
		},
	}); err != nil {
		return fmt.Errorf("initialize: %w", err)
	}
	return nil
}

// ListTools connects to an external MCP server and returns its tools.
func ListTools(ctx context.Context, srv config.MCPClientConfig) ([]mcpgo.Tool, error) {
	cl, err := newStdioClient(srv)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}
	defer cl.Close()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := initClient(ctx, cl); err != nil {
		return nil, err
	}
	res, err := cl.ListTools(ctx, mcpgo.ListToolsRequest{})
	if err != nil {
		return nil, fmt.Errorf("list tools: %w", err)
	}
	return res.Tools, nil
}

// CallToolString connects to an external MCP server, invokes a tool, and returns the result as a string.
func CallToolString(ctx context.Context, srv config.MCPClientConfig, toolName, argsJSON string) (string, error) {
	cl, err := newStdioClient(srv)
	if err != nil {
		return "", fmt.Errorf("create client: %w", err)
	}
	defer cl.Close()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := initClient(ctx, cl); err != nil {
		return "", err
	}
	var toolArgs map[string]interface{}
	if argsJSON != "" {
		if err := json.Unmarshal([]byte(argsJSON), &toolArgs); err != nil {
			return "", fmt.Errorf("parse args: %w", err)
		}
	}
	res, err := cl.CallTool(ctx, mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Name:      toolName,
			Arguments: toolArgs,
		},
	})
	if err != nil {
		return "", fmt.Errorf("call tool: %w", err)
	}
	b, _ := json.MarshalIndent(res, "", "  ")
	return string(b), nil
}

// CallTool connects to an external MCP server and invokes a tool, printing the result to stdout.
func CallTool(ctx context.Context, srv config.MCPClientConfig, toolName, argsJSON string) error {
	cl, err := newStdioClient(srv)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}
	defer cl.Close()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := initClient(ctx, cl); err != nil {
		return err
	}
	var toolArgs map[string]interface{}
	if argsJSON != "" {
		if err := json.Unmarshal([]byte(argsJSON), &toolArgs); err != nil {
			return fmt.Errorf("parse args: %w", err)
		}
	}
	res, err := cl.CallTool(ctx, mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Name:      toolName,
			Arguments: toolArgs,
		},
	})
	if err != nil {
		return fmt.Errorf("call tool: %w", err)
	}
	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))
	return nil
}
