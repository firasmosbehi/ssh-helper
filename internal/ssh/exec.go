package ssh

import (
	"fmt"
)

// RunCommand connects to a host and executes a single non-interactive command.
func RunCommand(opts ConnectOptions, cmd string) (string, error) {
	client, err := DialClient(opts)
	if err != nil {
		return "", fmt.Errorf("dial: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("new session: %w", err)
	}
	defer session.Close()

	out, err := session.CombinedOutput(cmd)
	return string(out), err
}

// TestConnectivity tries to open and close an SSH connection.
func TestConnectivity(opts ConnectOptions) error {
	client, err := DialClient(opts)
	if err != nil {
		return err
	}
	client.Close()
	return nil
}
