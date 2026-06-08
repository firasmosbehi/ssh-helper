package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/firasmosbehi/ssh-helper/internal/core"
	sshcfg "github.com/kevinburke/ssh_config"
)

// ParseConfig reads and parses the SSH config file at path.
func ParseConfig(path string) (*sshcfg.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return sshcfg.Decode(strings.NewReader(""))
		}
		return nil, fmt.Errorf("open ssh config: %w", err)
	}
	defer f.Close()
	cfg, err := sshcfg.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("parse ssh config: %w", err)
	}
	return cfg, nil
}

// HostsFromConfig extracts core.Host entries from the parsed config.
func HostsFromConfig(cfg *sshcfg.Config) []core.Host {
	var hosts []core.Host
	for _, h := range cfg.Hosts {
		if len(h.Patterns) == 0 {
			continue
		}
		// Skip implicit or explicit wildcard Host * blocks
		if h.Patterns[0].String() == "*" {
			continue
		}
		host := core.Host{}
		host.Name = h.Patterns[0].String()
		if v, _ := cfg.Get(host.Name, "HostName"); v != "" {
			host.Hostname = v
		} else {
			host.Hostname = host.Name
		}
		if v, _ := cfg.Get(host.Name, "User"); v != "" {
			host.User = v
		}
		if v, _ := cfg.Get(host.Name, "Port"); v != "" {
			fmt.Sscanf(v, "%d", &host.Port)
		}
		if v, _ := cfg.Get(host.Name, "IdentityFile"); v != "" {
			host.IdentityFile = v
		}
		hosts = append(hosts, host)
	}
	return hosts
}

// GetHost looks up a single host by name from the parsed config.
func GetHost(cfg *sshcfg.Config, name string) (core.Host, bool) {
	for _, h := range cfg.Hosts {
		if len(h.Patterns) == 0 {
			continue
		}
		if h.Patterns[0].String() == "*" {
			continue
		}
		if h.Patterns[0].String() == name {
			return hostFromConfig(cfg, h), true
		}
	}
	return core.Host{}, false
}

func hostFromConfig(cfg *sshcfg.Config, h *sshcfg.Host) core.Host {
	host := core.Host{}
	if len(h.Patterns) > 0 {
		host.Name = h.Patterns[0].String()
	}
	if v, _ := cfg.Get(host.Name, "HostName"); v != "" {
		host.Hostname = v
	} else {
		host.Hostname = host.Name
	}
	if v, _ := cfg.Get(host.Name, "User"); v != "" {
		host.User = v
	}
	if v, _ := cfg.Get(host.Name, "Port"); v != "" {
		fmt.Sscanf(v, "%d", &host.Port)
	}
	if v, _ := cfg.Get(host.Name, "IdentityFile"); v != "" {
		host.IdentityFile = v
	}
	return host
}

// BackupConfig creates a timestamped backup of the SSH config file.
func BackupConfig(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read config for backup: %w", err)
	}
	backup := path + ".backup." + time.Now().Format("20060102-150405")
	if err := os.WriteFile(backup, data, 0o600); err != nil {
		return "", fmt.Errorf("write backup: %w", err)
	}
	return backup, nil
}

// AddHost appends a Host block to the SSH config file.
func AddHost(path string, host core.Host) error {
	if err := ensureNewline(path); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open config for append: %w", err)
	}
	defer f.Close()

	var b strings.Builder
	b.WriteString(fmt.Sprintf("\nHost %s\n", host.Name))
	if host.Hostname != "" && host.Hostname != host.Name {
		b.WriteString(fmt.Sprintf("    HostName %s\n", host.Hostname))
	}
	if host.User != "" {
		b.WriteString(fmt.Sprintf("    User %s\n", host.User))
	}
	if host.Port != 0 {
		b.WriteString(fmt.Sprintf("    Port %d\n", host.Port))
	}
	if host.IdentityFile != "" {
		b.WriteString(fmt.Sprintf("    IdentityFile %s\n", host.IdentityFile))
	}
	_, err = f.WriteString(b.String())
	if err != nil {
		return fmt.Errorf("append host: %w", err)
	}
	return nil
}

// RemoveHost removes a Host block from the SSH config file.
func RemoveHost(path string, name string) error {
	lines, err := readLines(path)
	if err != nil {
		return err
	}
	start, end := findHostBlock(lines, name)
	if start == -1 {
		return fmt.Errorf("host %q not found in config", name)
	}
	newLines := append(lines[:start], lines[end:]...)
	return writeLines(path, newLines)
}

// EditHost updates an existing Host block.
func EditHost(path string, host core.Host) error {
	if _, err := BackupConfig(path); err != nil {
		return err
	}
	if err := RemoveHost(path, host.Name); err != nil {
		return err
	}
	return AddHost(path, host)
}

func findHostBlock(lines []string, name string) (start, end int) {
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "Host ") && !strings.HasPrefix(trimmed, "Host=") {
			continue
		}
		fields := strings.Fields(trimmed)
		if len(fields) >= 2 && fields[1] == name {
			start = i
			for j := i + 1; j < len(lines); j++ {
				t := strings.TrimSpace(lines[j])
				if strings.HasPrefix(t, "Host ") || strings.HasPrefix(t, "Host=") || strings.HasPrefix(t, "Match ") || strings.HasPrefix(t, "Match=") {
					return start, j
				}
			}
			return start, len(lines)
		}
	}
	return -1, -1
}

func ensureNewline(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			return err
		}
		return os.WriteFile(path, []byte{}, 0o600)
	}
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if len(data) > 0 && data[len(data)-1] != '\n' {
		return os.WriteFile(path, append(data, '\n'), info.Mode().Perm())
	}
	return nil
}

func readLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func writeLines(path string, lines []string) error {
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o600)
}
