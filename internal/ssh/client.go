package ssh

import (
	"fmt"
	"net"
	"os"

	"github.com/firasmosbehi/ssh-helper/internal/platform"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/term"
)

// ConnectOptions holds parameters for an SSH connection.
type ConnectOptions struct {
	Host     string
	Port     int
	User     string
	Identity string
}

// Connect establishes an interactive SSH session.
// DialClient creates an SSH client without starting an interactive session.
func DialClient(opts ConnectOptions) (*ssh.Client, error) {
	if opts.Port == 0 {
		opts.Port = 22
	}
	authMethods, err := buildAuthMethods(opts.Identity)
	if err != nil {
		return nil, err
	}
	cfg := &ssh.ClientConfig{
		User:            opts.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	return ssh.Dial("tcp", addr, cfg)
}

func Connect(opts ConnectOptions) error {
	if opts.Port == 0 {
		opts.Port = 22
	}

	authMethods, err := buildAuthMethods(opts.Identity)
	if err != nil {
		return err
	}

	cfg := &ssh.ClientConfig{
		User:            opts.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	client, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return fmt.Errorf("ssh dial %s: %w", addr, err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("make raw terminal: %w", err)
	}
	defer term.Restore(fd, oldState)

	w, h, err := term.GetSize(fd)
	if err != nil {
		w, h = 80, 24
	}

	if err := session.RequestPty("xterm-256color", h, w, ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}); err != nil {
		return fmt.Errorf("request pty: %w", err)
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		return fmt.Errorf("start shell: %w", err)
	}
	return session.Wait()
}

func buildAuthMethods(identity string) ([]ssh.AuthMethod, error) {
	var methods []ssh.AuthMethod

	if signers, err := getAgentSigners(); err == nil && len(signers) > 0 {
		methods = append(methods, ssh.PublicKeys(signers...))
	}

	if identity != "" {
		signer, err := parsePrivateKey(identity)
		if err != nil {
			return nil, fmt.Errorf("parse identity %s: %w", identity, err)
		}
		methods = append(methods, ssh.PublicKeys(signer))
	}

	methods = append(methods, ssh.PasswordCallback(func() (string, error) {
		return promptPassword("Password: ")
	}))

	return methods, nil
}

func getAgentSigners() ([]ssh.Signer, error) {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return nil, nil
	}
	conn, err := net.Dial("unix", sock)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return agent.NewClient(conn).Signers()
}

func parsePrivateKey(path string) (ssh.Signer, error) {
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		if _, ok := err.(*ssh.PassphraseMissingError); ok {
			// Try keyring first
			if pass, kerr := platform.GetKeyringSecret(path); kerr == nil {
				signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(pass))
				if err == nil {
					return signer, nil
				}
			}
			pass, perr := promptPassword(fmt.Sprintf("Enter passphrase for %s: ", path))
			if perr != nil {
				return nil, perr
			}
			signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(pass))
			if err == nil {
				_ = platform.SetKeyringSecret(path, pass)
			}
		}
	}
	return signer, err
}

func promptPassword(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	return string(b), err
}
