package config

import (
	"path/filepath"
	"testing"
)

func TestManagerInitAndLoad(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	m, err := NewManager()
	if err != nil {
		t.Fatalf("new manager: %v", err)
	}
	// override path to temp dir
	m.configPath = filepath.Join(tmp, "config.yaml")
	m.v.AddConfigPath(tmp)
	m.v.SetConfigFile(m.configPath)

	if err := m.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}

	cfg, err := m.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.LogLevel != "info" {
		t.Fatalf("expected default log level info, got %s", cfg.LogLevel)
	}
	if cfg.LogFormat != "text" {
		t.Fatalf("expected default log format text, got %s", cfg.LogFormat)
	}
}

func TestManagerSetAndGet(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	m, err := NewManager()
	if err != nil {
		t.Fatalf("new manager: %v", err)
	}
	m.configPath = filepath.Join(tmp, "config.yaml")
	m.v.SetConfigFile(m.configPath)

	if err := m.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}

	if err := m.Set("log_level", "debug"); err != nil {
		t.Fatalf("set: %v", err)
	}
	if got := m.Get("log_level"); got != "debug" {
		t.Fatalf("expected debug, got %v", got)
	}
}
