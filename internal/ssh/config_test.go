package ssh

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/firasmosbehi/ssh-helper/internal/core"
)

func TestHostsFromConfig(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config")
	data := `Host web1
    HostName 192.168.1.1
    User root
    Port 2222

Host web2
    HostName web2.example.com
`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := ParseConfig(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	hosts := HostsFromConfig(cfg)
	if len(hosts) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(hosts))
	}
	if hosts[0].Name != "web1" || hosts[0].User != "root" || hosts[0].Port != 2222 {
		t.Fatalf("unexpected web1: %+v", hosts[0])
	}
	if hosts[1].Name != "web2" {
		t.Fatalf("unexpected web2: %+v", hosts[1])
	}
}

func TestAddAndRemoveHost(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config")
	if err := os.WriteFile(path, []byte(""), 0o600); err != nil {
		t.Fatal(err)
	}

	h := core.Host{Name: "app1", Hostname: "10.0.0.1", User: "ubuntu", Port: 22}
	if err := AddHost(path, h); err != nil {
		t.Fatalf("add: %v", err)
	}

	cfg, _ := ParseConfig(path)
	hosts := HostsFromConfig(cfg)
	if len(hosts) != 1 || hosts[0].Name != "app1" {
		t.Fatalf("expected 1 host app1, got %+v", hosts)
	}

	if err := RemoveHost(path, "app1"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	cfg, _ = ParseConfig(path)
	hosts = HostsFromConfig(cfg)
	if len(hosts) != 0 {
		t.Fatalf("expected 0 hosts, got %d", len(hosts))
	}
}

func TestFindHostBlock(t *testing.T) {
	lines := []string{
		"Host a",
		"    HostName 1.1.1.1",
		"Host b",
		"    HostName 2.2.2.2",
		"",
	}
	start, end := findHostBlock(lines, "a")
	if start != 0 || end != 2 {
		t.Fatalf("expected block [0,2), got [%d,%d)", start, end)
	}
	start, end = findHostBlock(lines, "b")
	if start != 2 || end != 5 {
		t.Fatalf("expected block [2,5), got [%d,%d)", start, end)
	}
}
