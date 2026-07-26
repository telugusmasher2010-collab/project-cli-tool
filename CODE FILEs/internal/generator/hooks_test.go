package generator

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	apperrors "github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

func TestNewHookRunner(t *testing.T) {
	t.Run("empty hooks", func(t *testing.T) {
		hr := NewHookRunner()
		if hr.Len() != 0 {
			t.Errorf("Len() = %d, want 0", hr.Len())
		}
	})

	t.Run("single hook", func(t *testing.T) {
		h := &FuncHook{Name_: "h1", F: nil}
		hr := NewHookRunner(h)
		if hr.Len() != 1 {
			t.Errorf("Len() = %d, want 1", hr.Len())
		}
		if names := hr.Names(); len(names) != 1 || names[0] != "h1" {
			t.Errorf("Names() = %v, want [h1]", names)
		}
	})

	t.Run("multiple hooks", func(t *testing.T) {
		h1 := &FuncHook{Name_: "h1", F: nil}
		h2 := &FuncHook{Name_: "h2", F: nil}
		h3 := &FuncHook{Name_: "h3", F: nil}
		hr := NewHookRunner(h1, h2, h3)
		if hr.Len() != 3 {
			t.Errorf("Len() = %d, want 3", hr.Len())
		}
		names := hr.Names()
		want := []string{"h1", "h2", "h3"}
		for i, w := range want {
			if names[i] != w {
				t.Errorf("Names()[%d] = %q, want %q", i, names[i], w)
			}
		}
	})

	t.Run("nil hooks are filtered", func(t *testing.T) {
		h1 := &FuncHook{Name_: "h1", F: nil}
		hr := NewHookRunner(nil, h1, nil)
		if hr.Len() != 1 {
			t.Errorf("Len() = %d, want 1", hr.Len())
		}
	})

	t.Run("duplicate hooks are filtered", func(t *testing.T) {
		h1 := &FuncHook{Name_: "h1", F: nil}
		h2 := &FuncHook{Name_: "h1", F: nil}
		hr := NewHookRunner(h1, h2)
		if hr.Len() != 1 {
			t.Errorf("Len() = %d, want 1", hr.Len())
		}
	})

	t.Run("nil interface hook is filtered", func(t *testing.T) {
		hr := NewHookRunner(nil, nil)
		if hr.Len() != 0 {
			t.Errorf("Len() = %d, want 0", hr.Len())
		}
	})
}

func TestHookRunnerRunAll(t *testing.T) {
	t.Run("no hooks succeeds", func(t *testing.T) {
		hr := NewHookRunner()
		err := hr.RunAll(context.Background(), t.TempDir())
		if err != nil {
			t.Errorf("RunAll() error = %v, want nil", err)
		}
	})

	t.Run("single hook succeeds", func(t *testing.T) {
		var called bool
		h := &FuncHook{Name_: "h1", F: func(ctx context.Context, dir string) error {
			called = true
			return nil
		}}
		hr := NewHookRunner(h)
		err := hr.RunAll(context.Background(), t.TempDir())
		if err != nil {
			t.Errorf("RunAll() error = %v, want nil", err)
		}
		if !called {
			t.Error("hook was not called")
		}
	})

	t.Run("hooks run in order", func(t *testing.T) {
		var order []string
		var mu sync.Mutex
		mk := func(name string) *FuncHook {
			return &FuncHook{Name_: name, F: func(ctx context.Context, dir string) error {
				mu.Lock()
				order = append(order, name)
				mu.Unlock()
				return nil
			}}
		}
		hr := NewHookRunner(mk("a"), mk("b"), mk("c"))
		err := hr.RunAll(context.Background(), t.TempDir())
		if err != nil {
			t.Fatalf("RunAll() error = %v", err)
		}
		if len(order) != 3 {
			t.Fatalf("order length = %d, want 3", len(order))
		}
		for i, want := range []string{"a", "b", "c"} {
			if order[i] != want {
				t.Errorf("order[%d] = %q, want %q", i, order[i], want)
			}
		}
	})

	t.Run("stops on first failure", func(t *testing.T) {
		var calls []string
		var mu sync.Mutex
		mk := func(name string, fail bool) *FuncHook {
			return &FuncHook{Name_: name, F: func(ctx context.Context, dir string) error {
				mu.Lock()
				calls = append(calls, name)
				mu.Unlock()
				if fail {
					return errors.New("boom")
				}
				return nil
			}}
		}
		hr := NewHookRunner(mk("ok", false), mk("fail", true), mk("never", false))
		err := hr.RunAll(context.Background(), t.TempDir())
		if err == nil {
			t.Fatal("expected error from failing hook")
		}
		if !apperrors.IsHookFailed(err) {
			t.Errorf("expected ErrHookFailed, got: %v", err)
		}
		if len(calls) != 2 {
			t.Errorf("expected 2 hooks called, got %d: %v", len(calls), calls)
		}
		for _, c := range calls {
			if c == "never" {
				t.Errorf("hook 'never' should not have been called")
			}
		}
	})

	t.Run("error wraps hook name", func(t *testing.T) {
		h := &FuncHook{Name_: "bad-hook", F: func(ctx context.Context, dir string) error {
			return errors.New("kaboom")
		}}
		hr := NewHookRunner(h)
		err := hr.RunAll(context.Background(), t.TempDir())
		if err == nil {
			t.Fatal("expected error")
		}
		if !apperrors.IsHookFailed(err) {
			t.Errorf("error should wrap ErrHookFailed, got: %v", err)
		}
		if msg := err.Error(); !contains(msg, "bad-hook") {
			t.Errorf("error message should contain hook name, got: %s", msg)
		}
	})

	t.Run("cancelled context stops execution", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		var called bool
		h := &FuncHook{Name_: "h1", F: func(ctx context.Context, dir string) error {
			called = true
			return nil
		}}
		hr := NewHookRunner(h)
		err := hr.RunAll(ctx, t.TempDir())
		if err == nil {
			t.Fatal("expected error for cancelled context")
		}
		if called {
			t.Error("hook should not be called when context is cancelled")
		}
		if !apperrors.IsHookFailed(err) {
			t.Errorf("expected ErrHookFailed, got: %v", err)
		}
	})

	t.Run("nil function hook succeeds", func(t *testing.T) {
		h := &FuncHook{Name_: "nil-func", F: nil}
		hr := NewHookRunner(h)
		err := hr.RunAll(context.Background(), t.TempDir())
		if err != nil {
			t.Errorf("nil function hook should succeed, got: %v", err)
		}
	})
}

func TestFuncHook(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		h := &FuncHook{Name_: "my-hook", F: nil}
		if h.Name() != "my-hook" {
			t.Errorf("Name() = %q, want my-hook", h.Name())
		}
	})

	t.Run("nil function is no-op", func(t *testing.T) {
		h := &FuncHook{Name_: "nil", F: nil}
		err := h.Run(context.Background(), t.TempDir())
		if err != nil {
			t.Errorf("nil F should return nil, got: %v", err)
		}
	})

	t.Run("function is invoked", func(t *testing.T) {
		var called bool
		h := &FuncHook{Name_: "fn", F: func(ctx context.Context, dir string) error {
			called = true
			return nil
		}}
		err := h.Run(context.Background(), "/tmp/test")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !called {
			t.Error("function was not called")
		}
	})

	t.Run("function returns error", func(t *testing.T) {
		h := &FuncHook{Name_: "fn", F: func(ctx context.Context, dir string) error {
			return errors.New("fail")
		}}
		err := h.Run(context.Background(), "/tmp/test")
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "fail" {
			t.Errorf("error = %q, want fail", err.Error())
		}
	})

	t.Run("receives correct arguments", func(t *testing.T) {
		var gotDir string
		h := &FuncHook{Name_: "fn", F: func(ctx context.Context, dir string) error {
			gotDir = dir
			return nil
		}}
		dir := t.TempDir()
		h.Run(context.Background(), dir)
		if gotDir != dir {
			t.Errorf("dir = %q, want %q", gotDir, dir)
		}
	})

	t.Run("receives context", func(t *testing.T) {
		type key struct{}
		ctx := context.WithValue(context.Background(), key{}, "test-value")
		var got interface{}
		h := &FuncHook{Name_: "fn", F: func(ctx context.Context, dir string) error {
			got = ctx.Value(key{})
			return nil
		}}
		h.Run(ctx, "/tmp")
		if got != "test-value" {
			t.Errorf("context value = %v, want test-value", got)
		}
	})
}

func TestHookRunnerConcurrency(t *testing.T) {
	t.Run("hooks are not run concurrently", func(t *testing.T) {
		var count int64
		var maxConcurrent int64
		hooks := make([]Hook, 10)
		for i := 0; i < 10; i++ {
			hooks[i] = &FuncHook{
				Name_: "hook",
				F: func(ctx context.Context, dir string) error {
					cur := atomic.AddInt64(&count, 1)
					if cur > atomic.LoadInt64(&maxConcurrent) {
						atomic.StoreInt64(&maxConcurrent, cur)
					}
					atomic.AddInt64(&count, -1)
					return nil
				},
			}
		}
		hr := NewHookRunner(hooks...)
		err := hr.RunAll(context.Background(), t.TempDir())
		if err != nil {
			t.Fatalf("RunAll() error = %v", err)
		}
		// Hooks run sequentially, so maxConcurrent should be 1.
		if atomic.LoadInt64(&maxConcurrent) != 1 {
			t.Errorf("max concurrent = %d, want 1", atomic.LoadInt64(&maxConcurrent))
		}
	})
}

func TestHookRunnerErrorWrapping(t *testing.T) {
	t.Run("nested error chain", func(t *testing.T) {
		inner := errors.New("inner error")
		h := &FuncHook{Name_: "nested", F: func(ctx context.Context, dir string) error {
			return apperrors.Wrap(apperrors.ErrGenerationFailed, "gen failed", inner)
		}}
		hr := NewHookRunner(h)
		err := hr.RunAll(context.Background(), t.TempDir())
		if err == nil {
			t.Fatal("expected error")
		}
		if !apperrors.IsHookFailed(err) {
			t.Errorf("outermost should be ErrHookFailed, got: %v", err)
		}
		if !apperrors.IsGenerationFailed(err) {
			t.Errorf("chain should contain ErrGenerationFailed, got: %v", err)
		}
	})

	t.Run("multiple hook failures preserve first error", func(t *testing.T) {
		h1 := &FuncHook{Name_: "first", F: func(ctx context.Context, dir string) error {
			return errors.New("first-fail")
		}}
		h2 := &FuncHook{Name_: "second", F: func(ctx context.Context, dir string) error {
			return errors.New("second-fail")
		}}
		hr := NewHookRunner(h1, h2)
		err := hr.RunAll(context.Background(), t.TempDir())
		if err == nil {
			t.Fatal("expected error")
		}
		if msg := err.Error(); !contains(msg, "first") {
			t.Errorf("error should mention first failing hook, got: %s", msg)
		}
	})
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && searchSubstring(s, sub)
}

func searchSubstring(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
