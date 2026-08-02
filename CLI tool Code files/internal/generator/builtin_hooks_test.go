package generator

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	apperrors "github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

func TestBuiltinHookNames(t *testing.T) {
	tests := []struct {
		hook Hook
		want string
	}{
		{GitInitHook(), "git init"},
		{NpmInstallHook(), "npm install"},
		{FlutterPubGetHook(), "flutter pub get"},
		{GoModTidyHook(), "go mod tidy"},
	}
	for _, tt := range tests {
		if got := tt.hook.Name(); got != tt.want {
			t.Errorf("Name() = %q, want %q", got, tt.want)
		}
	}
}

func TestBuiltinHooksConfiguredCommands(t *testing.T) {
	tests := []struct {
		hook    Hook
		command string
		args    []string
	}{
		{GitInitHook(), "git", []string{"init"}},
		{NpmInstallHook(), "npm", []string{"install"}},
		{FlutterPubGetHook(), "flutter", []string{"pub", "get"}},
		{GoModTidyHook(), "go", []string{"mod", "tidy"}},
	}
	for _, tt := range tests {
		t.Run(tt.hook.Name(), func(t *testing.T) {
			ch, ok := tt.hook.(*CommandHook)
			if !ok {
				t.Fatalf("hook is %T, want *CommandHook", tt.hook)
			}
			if ch.Command != tt.command {
				t.Errorf("Command = %q, want %q", ch.Command, tt.command)
			}
			if len(ch.Args) != len(tt.args) {
				t.Fatalf("Args = %v, want %v", ch.Args, tt.args)
			}
			for i := range tt.args {
				if ch.Args[i] != tt.args[i] {
					t.Errorf("Args[%d] = %q, want %q", i, ch.Args[i], tt.args[i])
				}
			}
		})
	}
}

func TestHooksForTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		want     []string
	}{
		{"tauri-llm", "tauri-llm", []string{"git init", "npm install"}},
		{"whatsapp-bot", "whatsapp-bot", []string{"git init", "npm install"}},
		{"next-webapp", "next-webapp", []string{"git init", "npm install"}},
		{"react-native-map", "react-native-map", []string{"git init", "npm install"}},
		{"expense-splitter", "expense-splitter", []string{"git init", "flutter pub get"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HooksForTemplate(tt.template).Names()
			if len(got) != len(tt.want) {
				t.Fatalf("Names() = %v, want %v", got, tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("Names()[%d] = %q, want %q (full: %v)", i, got[i], tt.want[i], got)
				}
			}
		})
	}
}

func TestHooksForTemplateUnknown(t *testing.T) {
	got := HooksForTemplate("no-such-template").Names()
	want := []string{"git init"}
	if len(got) != 1 || got[0] != want[0] {
		t.Errorf("Names() = %v, want %v", got, want)
	}
}

func TestSelectHooks(t *testing.T) {
	t.Run("no hooks configured", func(t *testing.T) {
		if got := selectHooks(Options{}, "tauri-llm"); got != nil {
			t.Errorf("selectHooks() = %v, want nil", got)
		}
	})

	t.Run("auto hooks selected from template", func(t *testing.T) {
		got := selectHooks(Options{AutoHooks: true}, "expense-splitter")
		if got == nil {
			t.Fatal("selectHooks() = nil, want hook runner")
		}
		want := []string{"git init", "flutter pub get"}
		names := got.Names()
		for i := range want {
			if i >= len(names) || names[i] != want[i] {
				t.Fatalf("Names() = %v, want %v", names, want)
			}
		}
	})

	t.Run("explicit hooks take precedence", func(t *testing.T) {
		explicit := NewHookRunner(&FuncHook{Name_: "custom", F: nil})
		got := selectHooks(Options{Hooks: explicit, AutoHooks: true}, "tauri-llm")
		if got != explicit {
			t.Errorf("selectHooks() = %v, want explicit runner", got)
		}
	})
}

// fakeExecutable returns the file name and content for a cross-platform stub
// executable. A successful stub writes "<name>.txt" into its working
// directory; a failing stub exits non-zero.
func fakeExecutable(name string, fail bool) (fileName, content string) {
	if runtime.GOOS == "windows" {
		fileName = name + ".bat"
		if fail {
			content = "@echo off\r\nexit /b 1\r\n"
		} else {
			content = "@echo off\r\necho ran > " + name + ".txt\r\n"
		}
		return fileName, content
	}
	fileName = name
	if fail {
		content = "#!/bin/sh\nexit 1\n"
	} else {
		content = "#!/bin/sh\necho ran > " + name + ".txt\n"
	}
	return fileName, content
}

// withFakeBins writes stub executables into a temp directory, prepends it to
// PATH, and returns the stub directory so tests can exercise hooks end-to-end
// without invoking real git/npm/flutter binaries.
func withFakeBins(t *testing.T, bins ...fakeBin) {
	t.Helper()
	binDir := t.TempDir()
	for _, b := range bins {
		fileName, content := fakeExecutable(b.name, b.fail)
		if err := os.WriteFile(filepath.Join(binDir, fileName), []byte(content), 0755); err != nil {
			t.Fatalf("failed to write fake executable %q: %v", b.name, err)
		}
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

type fakeBin struct {
	name string
	fail bool
}

func TestGenerateRunsAutoHooks(t *testing.T) {
	withFakeBins(t, fakeBin{name: "git"}, fakeBin{name: "npm"})

	out := t.TempDir()
	g := New(out, nil, Options{AutoHooks: true})
	if err := g.Generate("tauri-llm"); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	for _, marker := range []string{"git.txt", "npm.txt"} {
		if _, err := os.Stat(filepath.Join(out, marker)); err != nil {
			t.Errorf("expected marker %q written by hook in project dir: %v", marker, err)
		}
	}
}

func TestGenerateRunsAutoHooksFlutter(t *testing.T) {
	withFakeBins(t, fakeBin{name: "git"}, fakeBin{name: "flutter"})

	out := t.TempDir()
	g := New(out, nil, Options{AutoHooks: true})
	if err := g.Generate("expense-splitter"); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	for _, marker := range []string{"git.txt", "flutter.txt"} {
		if _, err := os.Stat(filepath.Join(out, marker)); err != nil {
			t.Errorf("expected marker %q written by hook in project dir: %v", marker, err)
		}
	}
}

func TestGenerateAutoHooksStopOnFirstFailure(t *testing.T) {
	withFakeBins(t, fakeBin{name: "git", fail: true}, fakeBin{name: "npm"})

	out := t.TempDir()
	g := New(out, nil, Options{AutoHooks: true})
	err := g.Generate("tauri-llm")
	if err == nil {
		t.Fatal("expected error when git init fails")
	}
	if !apperrors.IsHookFailed(err) {
		t.Errorf("expected ErrHookFailed, got: %v", err)
	}
	if _, err := os.Stat(filepath.Join(out, "npm.txt")); err == nil {
		t.Error("npm hook should not run after git init fails")
	}
}

func TestGenerateExplicitHooksOverrideAutoHooks(t *testing.T) {
	withFakeBins(t, fakeBin{name: "git"}, fakeBin{name: "npm"})

	var calls []string
	var mu sync.Mutex
	explicit := NewHookRunner(&FuncHook{
		Name_: "explicit",
		F: func(_ context.Context, _ string) error {
			mu.Lock()
			calls = append(calls, "explicit")
			mu.Unlock()
			return nil
		},
	})

	out := t.TempDir()
	g := New(out, nil, Options{AutoHooks: true, Hooks: explicit})
	if err := g.Generate("tauri-llm"); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	mu.Lock()
	if len(calls) != 1 || calls[0] != "explicit" {
		t.Errorf("calls = %v, want [explicit]", calls)
	}
	mu.Unlock()

	for _, marker := range []string{"git.txt", "npm.txt"} {
		if _, err := os.Stat(filepath.Join(out, marker)); err == nil {
			t.Errorf("auto hook marker %q should not exist when explicit hooks override", marker)
		}
	}
}
