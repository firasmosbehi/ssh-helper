package main

import (
	"fmt"
	"os"
	"time"

	"github.com/firasmosbehi/ssh-helper/internal/core"
	"github.com/firasmosbehi/ssh-helper/internal/rsync"
	"github.com/firasmosbehi/ssh-helper/internal/store"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func newSyncCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Run and manage rsync jobs",
	}
	cmd.AddCommand(
		newSyncRunCommand(),
		newSyncJobCommand(),
	)
	return cmd
}

func newSyncRunCommand() *cobra.Command {
	var excludes []string
	var dryRun bool
	var flags []string
	cmd := &cobra.Command{
		Use:   "run <source> <dest>",
		Short: "Run a one-off rsync transfer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			job := core.SyncJob{
				Source:   args[0],
				Dest:     args[1],
				Excludes: excludes,
				DryRun:   dryRun,
				Flags:    flags,
			}
			runner := rsync.Runner{Job: job}
			res := runner.Run(cmd.Context(), rsync.RunOptions{
				Stdout:     os.Stdout,
				Stderr:     os.Stderr,
				OnProgress: func(p rsync.Progress) { /* CLI prints raw output */ },
			})
			if res.Error != nil {
				return res.Error
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&excludes, "exclude", nil, "exclude pattern(s)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "perform a trial run")
	cmd.Flags().StringSliceVar(&flags, "flag", nil, "additional rsync flags")
	return cmd
}

func newSyncJobCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "job",
		Short: "Manage persisted rsync jobs",
	}
	cmd.AddCommand(
		newSyncJobCreateCommand(),
		newSyncJobListCommand(),
		newSyncJobDeleteCommand(),
		newSyncJobRunCommand(),
		newSyncJobHistoryCommand(),
	)
	return cmd
}

func newSyncJobCreateCommand() *cobra.Command {
	var job core.SyncJob
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a persisted rsync job",
		RunE: func(cmd *cobra.Command, args []string) error {
			if job.Name == "" || job.Source == "" || job.Dest == "" {
				return fmt.Errorf("name, source, and dest are required")
			}
			job.ID = uuid.New().String()
			s, err := getStore()
			if err != nil {
				return err
			}
			return s.SaveSyncJob(job)
		},
	}
	cmd.Flags().StringVar(&job.Name, "name", "", "job name")
	cmd.Flags().StringVar(&job.Source, "source", "", "source path")
	cmd.Flags().StringVar(&job.Dest, "dest", "", "destination path")
	cmd.Flags().StringSliceVar(&job.Excludes, "exclude", nil, "exclude patterns")
	cmd.Flags().BoolVar(&job.DryRun, "dry-run", false, "dry run")
	cmd.Flags().StringSliceVar(&job.Flags, "flag", nil, "additional rsync flags")
	return cmd
}

func newSyncJobListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List persisted rsync jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := getStore()
			if err != nil {
				return err
			}
			jobs, err := s.ListSyncJobs()
			if err != nil {
				return err
			}
			for _, j := range jobs {
				fmt.Printf("%s\t%s\t%s -> %s\n", j.ID[:8], j.Name, j.Source, j.Dest)
			}
			return nil
		},
	}
}

func newSyncJobDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a persisted rsync job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := getStore()
			if err != nil {
				return err
			}
			return s.DeleteSyncJob(args[0])
		},
	}
}

func newSyncJobRunCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "run <id>",
		Short: "Run a persisted rsync job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := getStore()
			if err != nil {
				return err
			}
			jobs, err := s.ListSyncJobs()
			if err != nil {
				return err
			}
			var job core.SyncJob
			for _, j := range jobs {
				if j.ID == args[0] || len(j.ID) >= len(args[0]) && j.ID[:len(args[0])] == args[0] {
					job = j
					break
				}
			}
			if job.ID == "" {
				return fmt.Errorf("job %q not found", args[0])
			}

			runner := rsync.Runner{Job: job}
			res := runner.Run(cmd.Context(), rsync.RunOptions{
				Stdout:     os.Stdout,
				Stderr:     os.Stderr,
				CaptureLog: true,
			})

			run := core.JobRun{
				ID:        uuid.New().String(),
				JobID:     job.ID,
				StartedAt: time.Now(),
				EndedAt:   time.Now(),
				Status:    res.Status,
				Log:       res.Log,
			}
			if err := s.SaveRun(run); err != nil {
				return err
			}
			return res.Error
		},
	}
}

func newSyncJobHistoryCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "history",
		Short: "Show job run history",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := getStore()
			if err != nil {
				return err
			}
			runs, err := s.ListRuns()
			if err != nil {
				return err
			}
			for _, r := range runs {
				fmt.Printf("%s\t%s\t%s\t%s\n", r.ID[:8], r.JobID[:8], r.Status, r.StartedAt.Format(time.RFC3339))
			}
			return nil
		},
	}
}

func getStore() (store.Store, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return store.NewJSONStore(dir + "/ssh-helper"), nil
}
