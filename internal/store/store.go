package store

import (
	"github.com/firasmosbehi/ssh-helper/internal/core"
)

// Store defines persistence operations for application entities.
type Store interface {
	ListHosts() ([]core.Host, error)
	SaveHost(host core.Host) error
	DeleteHost(id string) error

	ListKeys() ([]core.Key, error)
	SaveKey(key core.Key) error
	DeleteKey(id string) error

	ListSyncJobs() ([]core.SyncJob, error)
	SaveSyncJob(job core.SyncJob) error
	DeleteSyncJob(id string) error
}
