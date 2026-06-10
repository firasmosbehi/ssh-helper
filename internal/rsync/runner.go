package rsync

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/firasmosbehi/ssh-helper/internal/core"
)

// Runner builds and executes rsync commands.
type Runner struct {
	Job core.SyncJob
}

// BuildArgs produces the rsync argument slice.
func (r *Runner) BuildArgs() []string {
	args := []string{"-avz"}
	if r.Job.DryRun {
		args = append(args, "-n")
	}
	for _, e := range r.Job.Excludes {
		args = append(args, "--exclude", e)
	}
	args = append(args, r.Job.Flags...)
	args = append(args, "--")
	args = append(args, r.Job.Source, r.Job.Dest)
	return args
}

// Progress holds parsed progress information.
type Progress struct {
	BytesTransferred string
	Percent          int
	Speed            string
	ETA              string
}

// RunOptions configures how Run behaves.
type RunOptions struct {
	Stdout     io.Writer
	Stderr     io.Writer
	OnProgress func(Progress)
	CaptureLog bool
}

// RunResult contains the output and metadata of a run.
type RunResult struct {
	Log    string
	Error  error
	Status string
}

// Run executes the rsync job.
func (r *Runner) Run(ctx context.Context, opts RunOptions) RunResult {
	args := r.BuildArgs()
	cmd := exec.CommandContext(ctx, "rsync", args...)

	var logBuf strings.Builder
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return RunResult{Error: fmt.Errorf("start rsync: %w", err), Status: "failed"}
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		scanOutput(stdoutPipe, opts.Stdout, opts.OnProgress, &logBuf, opts.CaptureLog)
	}()
	go func() {
		defer wg.Done()
		scanOutput(stderrPipe, opts.Stderr, nil, &logBuf, opts.CaptureLog)
	}()

	wg.Wait()
	err := cmd.Wait()
	status := "success"
	if err != nil {
		status = "failed"
	}
	return RunResult{Log: logBuf.String(), Error: err, Status: status}
}

func scanOutput(r io.Reader, out io.Writer, onProgress func(Progress), logBuf *strings.Builder, capture bool) {
	scanner := bufio.NewScanner(r)
	scanner.Split(scanLinesOrCR)
	for scanner.Scan() {
		line := scanner.Text()
		if out != nil {
			fmt.Fprintln(out, line)
		}
		if capture && logBuf != nil {
			fmt.Fprintln(logBuf, line)
		}
		if onProgress != nil {
			if p, ok := ParseProgress(line); ok {
				onProgress(p)
			}
		}
	}
}

func scanLinesOrCR(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' || data[i] == '\r' {
			return i + 1, data[:i], nil
		}
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

var progressRe = regexp.MustCompile(`(\d+)%`)
var speedRe = regexp.MustCompile(`([\d\.]+[KMGT]?B/s)`)

// ParseProgress attempts to extract progress from an rsync output line.
func ParseProgress(line string) (Progress, bool) {
	if !strings.Contains(line, "%") {
		return Progress{}, false
	}
	m := progressRe.FindStringSubmatch(line)
	if m == nil {
		return Progress{}, false
	}
	pct, _ := strconv.Atoi(m[1])
	p := Progress{Percent: pct}
	fields := strings.Fields(line)
	if len(fields) > 0 {
		p.BytesTransferred = fields[0]
	}
	if s := speedRe.FindString(line); s != "" {
		p.Speed = s
	}
	// Try to find ETA field (e.g. 0:00:05)
	for _, f := range fields {
		if strings.Contains(f, ":") && len(f) >= 5 {
			p.ETA = f
			break
		}
	}
	return p, true
}
