package core

import "time"

// JobRun represents an execution of a sync job.
type JobRun struct {
	ID        string    `json:"id"`
	JobID     string    `json:"job_id"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at,omitempty"`
	Status    string    `json:"status"`
	Log       string    `json:"log,omitempty"`
}
