package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Default is the package-level default logger.
var Default *slog.Logger

func init() {
	Default = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

// Options configure the logger.
type Options struct {
	Level  string
	Format string
	Output io.Writer
}

// New creates a structured logger from options.
func New(opts Options) (*slog.Logger, error) {
	level, err := parseLevel(opts.Level)
	if err != nil {
		return nil, err
	}
	if opts.Output == nil {
		opts.Output = os.Stderr
	}

	handlerOpts := &slog.HandlerOptions{Level: level}
	var handler slog.Handler
	switch strings.ToLower(opts.Format) {
	case "json":
		handler = slog.NewJSONHandler(opts.Output, handlerOpts)
	default:
		handler = slog.NewTextHandler(opts.Output, handlerOpts)
	}
	return slog.New(handler), nil
}

// Init sets the package-level Default logger.
func Init(opts Options) error {
	l, err := New(opts)
	if err != nil {
		return err
	}
	Default = l
	return nil
}

func parseLevel(s string) (slog.Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unknown log level %q", s)
	}
}
