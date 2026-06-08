package ssh

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// AddKeyToAgent adds a private key to the running SSH agent.
func AddKeyToAgent(path string) error {
	key, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read key: %w", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("parse key: %w", err)
	}

	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return fmt.Errorf("SSH_AUTH_SOCK not set")
	}
	conn, err := net.Dial("unix", sock)
	if err != nil {
		return fmt.Errorf("dial agent: %w", err)
	}
	defer conn.Close()

	return agent.NewClient(conn).Add(agent.AddedKey{PrivateKey: signer})
}

// AgentSigners returns the signers available from the SSH agent.
func AgentSigners() ([]ssh.Signer, error) {
	return getAgentSigners()
}
