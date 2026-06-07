package platform

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigDirUsesUserConfigDir(t *testing.T) {
	dir, err := ConfigDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir == "" {
		t.Fatal("expected non-empty dir")
	}
	if filepath.Base(dir) != "ssh-helper" {
		t.Fatalf("expected base ssh-helper, got %s", filepath.Base(dir))
	}
}

func TestEnsureDirCreatesDirectory(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "nested", "dir")
	if err := EnsureDir(target); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("expected directory to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected path to be a directory")
	}
}
