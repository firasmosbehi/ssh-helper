package core

// Host represents an SSH destination.
type Host struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Hostname     string   `json:"hostname"`
	User         string   `json:"user,omitempty"`
	Port         int      `json:"port,omitempty"`
	IdentityFile string   `json:"identity_file,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Group        string   `json:"group,omitempty"`
}
