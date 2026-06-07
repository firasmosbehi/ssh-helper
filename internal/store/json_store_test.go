package store

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/firasmosbehi/ssh-helper/internal/core"
)

func TestJSONStoreHostCRUD(t *testing.T) {
	dir := t.TempDir()
	s := NewJSONStore(dir)

	h := core.Host{ID: "h1", Name: "web1", Hostname: "web1.example.com", User: "root"}
	if err := s.SaveHost(h); err != nil {
		t.Fatalf("save host: %v", err)
	}

	hosts, err := s.ListHosts()
	if err != nil {
		t.Fatalf("list hosts: %v", err)
	}
	if len(hosts) != 1 {
		t.Fatalf("expected 1 host, got %d", len(hosts))
	}

	h.User = "admin"
	if err := s.SaveHost(h); err != nil {
		t.Fatalf("update host: %v", err)
	}
	hosts, _ = s.ListHosts()
	if hosts[0].User != "admin" {
		t.Fatalf("expected updated user admin, got %s", hosts[0].User)
	}

	if err := s.DeleteHost("h1"); err != nil {
		t.Fatalf("delete host: %v", err)
	}
	hosts, _ = s.ListHosts()
	if len(hosts) != 0 {
		t.Fatalf("expected 0 hosts, got %d", len(hosts))
	}
}

func TestJSONStoreFilesCreated(t *testing.T) {
	dir := t.TempDir()
	s := NewJSONStore(dir)
	if err := s.SaveHost(core.Host{ID: "h1", Name: "web1", Hostname: "web1.example.com"}); err != nil {
		t.Fatalf("save host: %v", err)
	}
	path := filepath.Join(dir, "hosts.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %s to exist: %v", path, err)
	}
}

func TestJSONStoreSyncJobCRUD(t *testing.T) {
	dir := t.TempDir()
	s := NewJSONStore(dir)

	j := core.SyncJob{ID: "j1", Name: "backup", Source: "/a", Dest: "user@host:/b"}
	if err := s.SaveSyncJob(j); err != nil {
		t.Fatalf("save job: %v", err)
	}
	jobs, err := s.ListSyncJobs()
	if err != nil {
		t.Fatalf("list jobs: %v", err)
	}
	if len(jobs) != 1 || jobs[0].Name != "backup" {
		t.Fatalf("unexpected jobs: %+v", jobs)
	}
	if err := s.DeleteSyncJob("j1"); err != nil {
		t.Fatalf("delete job: %v", err)
	}
	jobs, _ = s.ListSyncJobs()
	if len(jobs) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(jobs))
	}
}
