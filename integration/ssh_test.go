//go:build integration

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/firasmosbehi/ssh-helper/internal/core"
	"github.com/firasmosbehi/ssh-helper/internal/rsync"
	"github.com/firasmosbehi/ssh-helper/internal/ssh"
)

// TestSSHConnectivity tests connecting to the Docker-based SSH fixture.
func TestSSHConnectivity(t *testing.T) {
	opts := ssh.ConnectOptions{
		Host:     getEnv("SSH_HOST", "127.0.0.1"),
		Port:     getEnvInt("SSH_PORT", 2222),
		User:     getEnv("SSH_USER", "testuser"),
		Password: getEnv("SSH_PASSWORD", "testpass"),
	}
	if err := ssh.TestConnectivity(opts); err != nil {
		t.Fatalf("connectivity failed: %v", err)
	}
}

// TestSSHRunCommand tests remote command execution.
func TestSSHRunCommand(t *testing.T) {
	opts := ssh.ConnectOptions{
		Host:     getEnv("SSH_HOST", "127.0.0.1"),
		Port:     getEnvInt("SSH_PORT", 2222),
		User:     getEnv("SSH_USER", "testuser"),
		Password: getEnv("SSH_PASSWORD", "testpass"),
	}
	out, err := ssh.RunCommand(opts, "whoami")
	if err != nil {
		t.Fatalf("run command failed: %v", err)
	}
	if out != "testuser\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

// TestRSyncTransfer tests a basic rsync transfer to the fixture.
func TestRSyncTransfer(t *testing.T) {
	src := t.TempDir()
	dst := fmt.Sprintf("%s@%s:%s",
		getEnv("SSH_USER", "testuser"),
		getEnv("SSH_HOST", "127.0.0.1"),
		"/data/",
	)

	f, err := os.Create(filepath.Join(src, "hello.txt"))
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("hello")
	f.Close()

	job := core.SyncJob{
		Source: src + "/",
		Dest:   dst,
	}
	runner := rsync.Runner{Job: job}
	res := runner.Run(nil, rsync.RunOptions{})
	if res.Error != nil {
		t.Fatalf("rsync failed: %v", res.Error)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		var i int
		fmt.Sscanf(v, "%d", &i)
		return i
	}
	return fallback
}
