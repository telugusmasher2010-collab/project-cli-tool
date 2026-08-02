package generator

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	apperrors "github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

// CommandHook runs a single executable against the generated project
// directory. The executable is resolved through PATH and invoked directly
// without a shell, so the same hook works on Windows, macOS and Linux.
//
// CommandHook is safe for concurrent use: each Run call spawns its own
// process and shares no mutable state.
type CommandHook struct {
	// Name_ is the human-readable hook name used in errors and output.
	Name_ string
	// Command is the executable to run. It may be a bare name such as
	// "git" (resolved through PATH) or an absolute path.
	Command string
	// Args are the arguments passed to the executable.
	Args []string
	// Stdout receives the process's standard output. When nil, output is
	// streamed to os.Stdout.
	Stdout io.Writer
	// Stderr receives the process's standard error. When nil, output is
	// streamed to os.Stderr.
	Stderr io.Writer
}

// Name returns the hook's name.
func (h *CommandHook) Name() string { return h.Name_ }

// Run executes the command with projectDir as its working directory and
// streams its output. A non-zero exit code or an execution error is returned
// as an error wrapping ErrHookFailed.
func (h *CommandHook) Run(ctx context.Context, projectDir string) error {
	stdout := h.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := h.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	cmd := exec.CommandContext(ctx, h.Command, h.Args...)
	cmd.Dir = projectDir
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return apperrors.Wrap(apperrors.ErrHookFailed,
			fmt.Sprintf("hook %q failed in %q", h.Name(), projectDir), err)
	}
	return nil
}
