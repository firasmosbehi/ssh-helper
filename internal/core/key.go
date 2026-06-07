package core

// Key represents an SSH identity key.
type Key struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	Fingerprint string `json:"fingerprint,omitempty"`
}
