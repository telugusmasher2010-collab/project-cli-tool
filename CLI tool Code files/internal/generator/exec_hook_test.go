package generator

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	apperrors "github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

// TestHelperProcess is not a real test. It is re-invoked as a subprocess by
// CommandHook tests so process behavior can be exercised without depending
// on a particular external binary being installed on the host.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	switch os.Getenv("GO_HELPER_BEHAVIOR") {
	case "stream":
		_, _ = os.Stdout.WriteString("stdout-line\n")
		_, _ = os.Stderr.WriteString("stderr-line\n")
		os.Exit(0)
	case "write-cwd":
		dir, _ := os.Getwd()
		_ = os.WriteFile("helper-cwd.txt", []byte(dir), 0644)
		os.Exit(0)
	case "fail":
		os.Exit(1)
	case "sleep":
		_ = os.WriteFile("helper-started.txt", []byte("started"), 0644)
		time.Sleep(30 * time.Second)
		os.Exit(0)
	}
	os.Exit(0)
}

func helperExecutable(t *testing.T) string {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error = %v", err)
	}
	return exe
}

func helperProcessArgs() []string {
	return []string{"-test.run", "^TestHelperProcess$"}
}

func waitForFile(t *testing.T, path string) bool {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(path); err == nil {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestCommandHookName(t *testing.T) {
	h := &CommandHook{Name_: "my hook", Command: "go"}
	if got := h.Name(); got != "my hook" {
		t.Errorf("Name() = %q, want %q", got, "my hook")
	}
}

func TestCommandHookRunSuccess(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")
	t.Setenv("GO_HELPER_BEHAVIOR", "stream")

	var stdout, stderr bytes.Buffer
	h := &CommandHook{
		Name_:   "streamer",
		Command: helperExecutable(t),
		Args:    helperProcessArgs(),
		Stdout:  &stdout,
		Stderr:  &stderr,
	}

	if err := h.Run(context.Background(), t.TempDir()); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}
	if got := stdout.String(); !strings.Contains(got, "stdout-line") {
		t.Errorf("stdout = %q, want it to contain %q", got, "stdout-line")
	}
	if got := stderr.String(); !strings.Contains(got, "stderr-line") {
		t.Errorf("stderr = %q, want it to contain %q", got, "stderr-line")
	}
}

func TestCommandHookRunsInProjectDir(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")
	t.Setenv("GO_HELPER_BEHAVIOR", "write-cwd")

	dir := t.TempDir()
	h := &CommandHook{
		Name_:   "cwd",
		Command: helperExecutable(t),
		Args:    helperProcessArgs(),
	}

	if err := h.Run(context.Background(), dir); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}

	// The helper writes its working directory into helper-cwd.txt; if the
	// file exists inside the project dir, the command ran there.
	data, err := os.ReadFile(filepath.Join(dir, "helper-cwd.txt"))
	if err != nil {
		t.Fatalf("marker file not written in project dir: %v", err)
	}
	got, err := filepath.EvalSymlinks(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatalf("cannot resolve recorded cwd: %v", err)
	}
	want, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("cannot resolve project dir: %v", err)
	}
	if !strings.EqualFold(got, want) {
		t.Errorf("command cwd = %q, want %q", got, want)
	}
}

func TestCommandHookRunFailure(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")
	t.Setenv("GO_HELPER_BEHAVIOR", "fail")

	h := &CommandHook{
		Name_:   "failing",
		Command: helperExecutable(t),
		Args:    helperProcessArgs(),
	}

	err := h.Run(context.Background(), t.TempDir())
	if err == nil {
		t.Fatal("Run() error = nil, want non-nil")
	}
	if !apperrors.IsHookFailed(err) {
		t.Errorf("error should wrap ErrHookFailed, got: %v", err)
	}
	if msg := err.Error(); !strings.Contains(msg, "failing") {
		t.Errorf("error message should mention hook name, got: %s", msg)
	}
}

func TestCommandHookInvalidExecutable(t *testing.T) {
	h := &CommandHook{
		Name_:   "ghost",
		Command: "definitely-not-a-real-command-9f2c1a",
	}

	err := h.Run(context.Background(), t.TempDir())
	if err == nil {
		t.Fatal("expected error for missing executable")
	}
	if !apperrors.IsHookFailed(err) {
		t.Errorf("error should wrap ErrHookFailed, got: %v", err)
	}
	if msg := err.Error(); !strings.Contains(msg, "ghost") {
		t.Errorf("error message should mention hook name, got: %s", msg)
	}
}

func TestCommandHookCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	h := &CommandHook{
		Name_:   "cancelled",
		Command: helperExecutable(t),
		Args:    helperProcessArgs(),
	}

	err := h.Run(ctx, t.TempDir())
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
	if !apperrors.IsHookFailed(err) {
		t.Errorf("error should wrap ErrHookFailed, got: %v", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("error should wrap context.Canceled, got: %v", err)
	}
}

func TestCommandHookKillsRunningProcess(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")
	t.Setenv("GO_HELPER_BEHAVIOR", "sleep")

	dir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h := &CommandHook{
		Name_:   "sleepy",
		Command: helperExecutable(t),
		Args:    helperProcessArgs(),
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- h.Run(ctx, dir)
	}()

	if !waitForFile(t, filepath.Join(dir, "helper-started.txt")) {
		t.Fatal("helper process did not start in time")
	}
	cancel()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected error after context cancellation")
		}
		if !apperrors.IsHookFailed(err) {
			t.Errorf("error should wrap ErrHookFailed, got: %v", err)
		}
		// Note: on Windows a hard-killed process surfaces as an exit error
		// (e.g. "exit status 1") rather than context.Canceled, so we only
		// assert that cancellation aborted the run and produced a wrapped
		// ErrHookFailed.
	case <-time.After(10 * time.Second):
		t.Fatal("hook did not stop after context cancellation")
	}
}
