package generator

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	apperrors "github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

func TestGenerateRunsHooksAfterGeneration(t *testing.T) {
	var ran bool
	var ranDir string
	runner := NewHookRunner(&FuncHook{
		Name_: "post-gen",
		F: func(_ context.Context, dir string) error {
			ran = true
			ranDir = dir
			return nil
		},
	})

	out := t.TempDir()
	g := New(out, nil, Options{Hooks: runner})
	if err := g.Generate("tauri-llm"); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !ran {
		t.Fatal("hook was not run after generation")
	}

	gotClean, err := filepath.EvalSymlinks(ranDir)
	if err != nil {
		t.Fatalf("cannot resolve hook dir: %v", err)
	}
	wantClean, err := filepath.EvalSymlinks(out)
	if err != nil {
		t.Fatalf("cannot resolve output dir: %v", err)
	}
	if !strings.EqualFold(gotClean, wantClean) {
		t.Errorf("hook ran in %q, want %q", ranDir, out)
	}
}

func TestGenerateRunsHooksInOrder(t *testing.T) {
	var order []string
	var mu sync.Mutex
	mk := func(name string) *FuncHook {
		return &FuncHook{Name_: name, F: func(_ context.Context, _ string) error {
			mu.Lock()
			order = append(order, name)
			mu.Unlock()
			return nil
		}}
	}
	runner := NewHookRunner(mk("first"), mk("second"), mk("third"))

	out := t.TempDir()
	g := New(out, nil, Options{Hooks: runner})
	if err := g.Generate("tauri-llm"); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	want := []string{"first", "second", "third"}
	mu.Lock()
	defer mu.Unlock()
	if len(order) != len(want) {
		t.Fatalf("order = %v, want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Errorf("order[%d] = %q, want %q", i, order[i], want[i])
		}
	}
}

func TestGenerateStopsOnFirstHookFailure(t *testing.T) {
	var calls []string
	var mu sync.Mutex
	mk := func(name string, fail bool) *FuncHook {
		return &FuncHook{Name_: name, F: func(_ context.Context, _ string) error {
			mu.Lock()
			calls = append(calls, name)
			mu.Unlock()
			if fail {
				return errors.New("boom")
			}
			return nil
		}}
	}
	runner := NewHookRunner(mk("ok", false), mk("fail", true), mk("never", false))

	out := t.TempDir()
	g := New(out, nil, Options{Hooks: runner})
	err := g.Generate("tauri-llm")
	if err == nil {
		t.Fatal("expected error from failing hook")
	}
	if !apperrors.IsHookFailed(err) {
		t.Errorf("error should wrap ErrHookFailed, got: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(calls) != 2 {
		t.Errorf("expected 2 hooks called, got %v", calls)
	}
	for _, c := range calls {
		if c == "never" {
			t.Error("hook 'never' should not run after failure")
		}
	}
}

func TestGenerateHookErrorWrapping(t *testing.T) {
	runner := NewHookRunner(&FuncHook{
		Name_: "boom",
		F: func(_ context.Context, _ string) error {
			return errors.New("kaboom")
		},
	})

	out := t.TempDir()
	g := New(out, nil, Options{Hooks: runner})
	err := g.Generate("tauri-llm")
	if err == nil {
		t.Fatal("expected error")
	}
	if !apperrors.IsHookFailed(err) {
		t.Errorf("error should wrap ErrHookFailed, got: %v", err)
	}
	if msg := err.Error(); !strings.Contains(msg, "boom") {
		t.Errorf("error message should mention hook name, got: %s", msg)
	}
}

func TestGenerateNoHooksWhenNoneConfigured(t *testing.T) {
	out := t.TempDir()
	g := New(out, nil, Options{})
	if err := g.Generate("tauri-llm"); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
}
