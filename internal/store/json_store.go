package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/firasmosbehi/ssh-helper/internal/core"
)

// JSONStore persists entities as JSON files on disk.
type JSONStore struct {
	dir string
	mu  sync.RWMutex
}

// NewJSONStore creates a store backed by JSON files in dir.
func NewJSONStore(dir string) *JSONStore {
	return &JSONStore{dir: dir}
}

func (s *JSONStore) path(name string) string {
	return filepath.Join(s.dir, name+".json")
}

func (s *JSONStore) read(name string, out interface{}) error {
	path := s.path(name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read %s: %w", path, err)
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, out)
}

func (s *JSONStore) write(name string, data interface{}) error {
	path := s.path(name)
	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return fmt.Errorf("ensure store dir: %w", err)
	}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", name, err)
	}
	return os.WriteFile(path, b, 0o600)
}

// ListHosts returns all persisted hosts.
func (s *JSONStore) ListHosts() ([]core.Host, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var hosts []core.Host
	if err := s.read("hosts", &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}

// SaveHost inserts or updates a host by ID.
func (s *JSONStore) SaveHost(host core.Host) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	hosts, err := s.listHostsUnlocked()
	if err != nil {
		return err
	}
	hosts = upsert(hosts, host, func(a, b core.Host) bool { return a.ID == b.ID }, host)
	return s.write("hosts", hosts)
}

// DeleteHost removes a host by ID.
func (s *JSONStore) DeleteHost(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	hosts, err := s.listHostsUnlocked()
	if err != nil {
		return err
	}
	hosts = remove(hosts, func(h core.Host) bool { return h.ID == id })
	return s.write("hosts", hosts)
}

func (s *JSONStore) listHostsUnlocked() ([]core.Host, error) {
	var hosts []core.Host
	if err := s.read("hosts", &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}

// ListKeys returns all persisted keys.
func (s *JSONStore) ListKeys() ([]core.Key, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var keys []core.Key
	if err := s.read("keys", &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

// SaveKey inserts or updates a key by ID.
func (s *JSONStore) SaveKey(key core.Key) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys, err := s.listKeysUnlocked()
	if err != nil {
		return err
	}
	keys = upsert(keys, key, func(a, b core.Key) bool { return a.ID == b.ID }, key)
	return s.write("keys", keys)
}

// DeleteKey removes a key by ID.
func (s *JSONStore) DeleteKey(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys, err := s.listKeysUnlocked()
	if err != nil {
		return err
	}
	keys = remove(keys, func(k core.Key) bool { return k.ID == id })
	return s.write("keys", keys)
}

func (s *JSONStore) listKeysUnlocked() ([]core.Key, error) {
	var keys []core.Key
	if err := s.read("keys", &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

// ListSyncJobs returns all persisted sync jobs.
func (s *JSONStore) ListSyncJobs() ([]core.SyncJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var jobs []core.SyncJob
	if err := s.read("jobs", &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

// SaveSyncJob inserts or updates a sync job by ID.
func (s *JSONStore) SaveSyncJob(job core.SyncJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	jobs, err := s.listJobsUnlocked()
	if err != nil {
		return err
	}
	jobs = upsert(jobs, job, func(a, b core.SyncJob) bool { return a.ID == b.ID }, job)
	return s.write("jobs", jobs)
}

// DeleteSyncJob removes a sync job by ID.
func (s *JSONStore) DeleteSyncJob(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	jobs, err := s.listJobsUnlocked()
	if err != nil {
		return err
	}
	jobs = remove(jobs, func(j core.SyncJob) bool { return j.ID == id })
	return s.write("jobs", jobs)
}

func (s *JSONStore) listJobsUnlocked() ([]core.SyncJob, error) {
	var jobs []core.SyncJob
	if err := s.read("jobs", &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

// ListRuns returns all persisted job runs.
func (s *JSONStore) ListRuns() ([]core.JobRun, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var runs []core.JobRun
	if err := s.read("runs", &runs); err != nil {
		return nil, err
	}
	return runs, nil
}

// SaveRun inserts or updates a job run by ID.
func (s *JSONStore) SaveRun(run core.JobRun) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	runs, err := s.listRunsUnlocked()
	if err != nil {
		return err
	}
	runs = upsert(runs, run, func(a, b core.JobRun) bool { return a.ID == b.ID }, run)
	return s.write("runs", runs)
}

// DeleteRun removes a job run by ID.
func (s *JSONStore) DeleteRun(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	runs, err := s.listRunsUnlocked()
	if err != nil {
		return err
	}
	runs = remove(runs, func(r core.JobRun) bool { return r.ID == id })
	return s.write("runs", runs)
}

func (s *JSONStore) listRunsUnlocked() ([]core.JobRun, error) {
	var runs []core.JobRun
	if err := s.read("runs", &runs); err != nil {
		return nil, err
	}
	return runs, nil
}

func upsert[S ~[]E, E any](s S, item E, eq func(a, b E) bool, with E) S {
	for i, v := range s {
		if eq(v, item) {
			s[i] = with
			return s
		}
	}
	return append(s, with)
}

func remove[S ~[]E, E any](s S, match func(E) bool) S {
	out := s[:0]
	for _, v := range s {
		if !match(v) {
			out = append(out, v)
		}
	}
	return out
}
