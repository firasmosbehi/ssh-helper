package logger

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestNewTextLogger(t *testing.T) {
	var buf bytes.Buffer
	l, err := New(Options{Level: "debug", Format: "text", Output: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	l.Info("hello")
	if !strings.Contains(buf.String(), "hello") {
		t.Fatalf("expected log to contain hello, got %s", buf.String())
	}
}

func TestNewJSONLogger(t *testing.T) {
	var buf bytes.Buffer
	l, err := New(Options{Level: "info", Format: "json", Output: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	l.Info("hello")
	if !strings.Contains(buf.String(), `"msg":"hello"`) {
		t.Fatalf("expected JSON msg, got %s", buf.String())
	}
}

func TestParseLevel(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
	} {
		got, err := parseLevel(tc.in)
		if err != nil {
			t.Fatalf("parseLevel(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("parseLevel(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
