package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/firasmosbehi/ssh-helper/internal/platform"
	"github.com/spf13/viper"
)

// MCPClientConfig holds the configuration for an external MCP server.
type MCPClientConfig struct {
	Command string            `mapstructure:"command"`
	Args    []string          `mapstructure:"args"`
	Env     map[string]string `mapstructure:"env"`
}

// Config is the user-facing configuration shape.
type Config struct {
	DataDir       string                    `mapstructure:"data_dir"`
	SSHConfigPath string                    `mapstructure:"ssh_config_path"`
	LogFormat     string                    `mapstructure:"log_format"`
	LogLevel      string                    `mapstructure:"log_level"`
	MCPClients    map[string]MCPClientConfig `mapstructure:"mcp_servers"`
}

// Manager wraps Viper and the on-disk config file path.
type Manager struct {
	v          *viper.Viper
	configPath string
}

// NewManager creates a Manager with config paths initialized.
func NewManager() (*Manager, error) {
	dir, err := platform.ConfigDir()
	if err != nil {
		return nil, err
	}
	v := viper.New()
	v.SetDefault("data_dir", dir)
	v.SetDefault("ssh_config_path", filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	v.SetDefault("log_format", "text")
	v.SetDefault("log_level", "info")
	v.SetDefault("mcp_servers", map[string]MCPClientConfig{})

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(dir)

	return &Manager{
		v:          v,
		configPath: filepath.Join(dir, "config.yaml"),
	}, nil
}

// Path returns the full path to the config file.
func (m *Manager) Path() string {
	return m.configPath
}

// Load reads the configuration from disk. If no file exists, defaults are returned.
func (m *Manager) Load() (*Config, error) {
	if err := m.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}
	var cfg Config
	if err := m.v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}

// Init writes the default configuration file if it does not already exist.
func (m *Manager) Init() error {
	if err := platform.EnsureDir(filepath.Dir(m.configPath)); err != nil {
		return err
	}
	if _, err := os.Stat(m.configPath); err == nil {
		return fmt.Errorf("config already exists at %s", m.configPath)
	}
	return m.writeDefault()
}

// Get returns the raw value for a config key.
func (m *Manager) Get(key string) interface{} {
	return m.v.Get(key)
}

// Set updates a config key and persists the file.
func (m *Manager) Set(key string, value interface{}) error {
	m.v.Set(key, value)
	return m.v.WriteConfigAs(m.configPath)
}

// Save persists the current in-memory config state to disk.
func (m *Manager) Save() error {
	if err := platform.EnsureDir(filepath.Dir(m.configPath)); err != nil {
		return err
	}
	return m.v.WriteConfigAs(m.configPath)
}

// Render returns the current config as indented JSON for display.
func (m *Manager) Render(key string) (string, error) {
	val := m.v.Get(key)
	b, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (m *Manager) writeDefault() error {
	if err := platform.EnsureDir(filepath.Dir(m.configPath)); err != nil {
		return err
	}
	return m.v.WriteConfigAs(m.configPath)
}
