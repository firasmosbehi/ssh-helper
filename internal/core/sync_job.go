package core

// SyncJob represents a persisted rsync job.
type SyncJob struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Source   string   `json:"source"`
	Dest     string   `json:"dest"`
	Flags    []string `json:"flags,omitempty"`
	Excludes []string `json:"excludes,omitempty"`
	DryRun   bool     `json:"dry_run,omitempty"`
}
