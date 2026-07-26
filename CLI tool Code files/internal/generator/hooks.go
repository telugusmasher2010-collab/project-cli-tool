package generator

import (
	"context"
	"fmt"

	apperrors "github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

// Hook defines a post-generation action that runs after files are created.
// Implementations must be safe for concurrent use.
type Hook interface {
	// Name returns a human-readable identifier for this hook.
	Name() string
	// Run executes the hook against the generated project directory.
	// If the hook fails, it must return an error wrapping ErrHookFailed.
	Run(ctx context.Context, projectDir string) error
}

// HookRunner executes a sequence of hooks against a project directory.
// Hooks run in registration order; execution stops on the first failure.
type HookRunner struct {
	hooks []Hook
}

// NewHookRunner creates a HookRunner with the given hooks.
// Nil entries and duplicate hooks are silently ignored.
func NewHookRunner(hooks ...Hook) *HookRunner {
	seen := make(map[string]bool)
	filtered := make([]Hook, 0, len(hooks))
	for _, h := range hooks {
		if h == nil {
			continue
		}
		name := h.Name()
		if seen[name] {
			continue
		}
		seen[name] = true
		filtered = append(filtered, h)
	}
	return &HookRunner{hooks: filtered}
}

// RunAll executes every registered hook in order. If any hook returns an
// error, execution stops immediately and that error (wrapped with
// ErrHookFailed and the hook name) is returned.
func (hr *HookRunner) RunAll(ctx context.Context, projectDir string) error {
	for _, h := range hr.hooks {
		if ctx.Err() != nil {
			return apperrors.Wrap(apperrors.ErrHookFailed,
				fmt.Sprintf("hook %q: context cancelled", h.Name()), ctx.Err())
		}
		if err := h.Run(ctx, projectDir); err != nil {
			return apperrors.Wrap(apperrors.ErrHookFailed,
				fmt.Sprintf("hook %q failed", h.Name()), err)
		}
	}
	return nil
}

// Len returns the number of registered hooks.
func (hr *HookRunner) Len() int {
	return len(hr.hooks)
}

// Names returns the names of all registered hooks in order.
func (hr *HookRunner) Names() []string {
	names := make([]string, len(hr.hooks))
	for i, h := range hr.hooks {
		names[i] = h.Name()
	}
	return names
}

// FuncHook is an adapter to allow the use of ordinary functions as Hook
// implementations. If f is nil, Run is a no-op.
type FuncHook struct {
	Name_ string
	F     func(ctx context.Context, projectDir string) error
}

// Name returns the hook's name.
func (fh *FuncHook) Name() string { return fh.Name_ }

// Run invokes the underlying function. A nil function is treated as a
// successful no-op.
func (fh *FuncHook) Run(ctx context.Context, projectDir string) error {
	if fh.F == nil {
		return nil
	}
	return fh.F(ctx, projectDir)
}
