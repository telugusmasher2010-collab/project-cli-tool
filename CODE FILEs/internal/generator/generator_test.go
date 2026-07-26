package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	apperrors "github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

func TestNew(t *testing.T) {
	t.Run("creates generator with variables", func(t *testing.T) {
		vars := NewVariables()
		vars.Set("name", "proj")
		g := New("/tmp/out", vars, Options{})
		if g == nil {
			t.Fatal("New returned nil")
		}
		if g.outputDir != "/tmp/out" {
			t.Errorf("outputDir = %q, want /tmp/out", g.outputDir)
		}
		got, _ := g.vars.Get("name")
		if got != "proj" {
			t.Errorf("vars.name = %q, want proj", got)
		}
	})

	t.Run("creates generator with nil variables", func(t *testing.T) {
		g := New("/tmp/out", nil, Options{})
		if g.vars == nil {
			t.Fatal("vars should be initialized even when nil is passed")
		}
		if len(g.vars.Keys()) != 0 {
			t.Error("new vars should be empty")
		}
	})

	t.Run("copies options", func(t *testing.T) {
		g := New("/tmp/out", nil, Options{Overwrite: true})
		if !g.opts.Overwrite {
			t.Error("Overwrite option not preserved")
		}
	})
}

func TestVariables(t *testing.T) {
	vars := NewVariables()
	vars.Set("key", "val")
	g := New("/tmp/out", vars, Options{})

	got := g.Variables()
	if got != vars {
		t.Error("Variables() should return the same instance")
	}
	v, _ := got.Get("key")
	if v != "val" {
		t.Errorf("Variables().Get(key) = %q, want val", v)
	}
}

func TestGenerate(t *testing.T) {
	t.Run("successful generation from tauri-llm", func(t *testing.T) {
		out := t.TempDir()
		g := New(out, nil, Options{})

		err := g.Generate("tauri-llm")
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		// placeholder.txt should exist in the output
		data, err := os.ReadFile(filepath.Join(out, "tauri-llm", "placeholder.txt"))
		if err != nil {
			t.Fatalf("failed to read generated file: %v", err)
		}
		if !strings.Contains(string(data), "Placeholder template directory") {
			t.Errorf("unexpected file content: %q", string(data))
		}
	})

	t.Run("successful generation from whatsapp-bot", func(t *testing.T) {
		out := t.TempDir()
		g := New(out, nil, Options{})

		err := g.Generate("whatsapp-bot")
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		info, err := os.Stat(filepath.Join(out, "whatsapp-bot", "placeholder.txt"))
		if err != nil {
			t.Fatalf("generated file not found: %v", err)
		}
		if info.Mode().Perm() != defaultFilePerm {
			t.Errorf("file perm = %o, want %o", info.Mode().Perm(), defaultFilePerm)
		}
	})

	t.Run("successful generation from expense-splitter", func(t *testing.T) {
		out := t.TempDir()
		g := New(out, nil, Options{})

		err := g.Generate("expense-splitter")
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		entries, err := os.ReadDir(filepath.Join(out, "expense-splitter"))
		if err != nil {
			t.Fatalf("output dir not created: %v", err)
		}
		if len(entries) == 0 {
			t.Error("generated directory is empty")
		}
	})

	t.Run("template not found", func(t *testing.T) {
		out := t.TempDir()
		g := New(out, nil, Options{})

		err := g.Generate("nonexistent-template")
		if err == nil {
			t.Fatal("expected error for nonexistent template")
		}
		if !apperrors.IsTemplateNotFound(err) {
			t.Errorf("expected ErrTemplateNotFound, got: %v", err)
		}
	})

	t.Run("output directory does not exist yet", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "brand-new-project")
		g := New(out, nil, Options{})

		err := g.Generate("tauri-llm")
		if err != nil {
			t.Fatalf("Generate() should create output dir, error = %v", err)
		}

		info, err := os.Stat(out)
		if err != nil {
			t.Fatalf("output dir should exist: %v", err)
		}
		if !info.IsDir() {
			t.Error("output path should be a directory")
		}
	})

	t.Run("output exists and is empty — allowed", func(t *testing.T) {
		out := t.TempDir()
		g := New(out, nil, Options{})

		err := g.Generate("tauri-llm")
		if err != nil {
			t.Fatalf("first generation error = %v", err)
		}

		// Second generation into same dir — directory is non-empty now.
		g2 := New(out, nil, Options{})
		err = g2.Generate("tauri-llm")
		if err == nil {
			t.Fatal("expected error for non-empty output without overwrite")
		}
		if !apperrors.IsOutputExists(err) {
			t.Errorf("expected ErrOutputExists, got: %v", err)
		}
	})

	t.Run("output exists with overwrite enabled", func(t *testing.T) {
		out := t.TempDir()
		g := New(out, nil, Options{Overwrite: true})

		err := g.Generate("tauri-llm")
		if err != nil {
			t.Fatalf("first generation error = %v", err)
		}

		// Should succeed with overwrite.
		g2 := New(out, nil, Options{Overwrite: true})
		err = g2.Generate("tauri-llm")
		if err != nil {
			t.Fatalf("Generate() with overwrite should succeed, error = %v", err)
		}
	})

	t.Run("output path is a file, not a directory", func(t *testing.T) {
		out := t.TempDir()
		fakePath := filepath.Join(out, "not-a-dir")
		if err := os.WriteFile(fakePath, []byte("hi"), 0644); err != nil {
			t.Fatal(err)
		}

		g := New(fakePath, nil, Options{})
		err := g.Generate("tauri-llm")
		if err == nil {
			t.Fatal("expected error when output path is a file")
		}
		if !apperrors.IsFilesystem(err) {
			t.Errorf("expected ErrFilesystem, got: %v", err)
		}
	})

	t.Run("overwrite does not check emptiness", func(t *testing.T) {
		out := t.TempDir()
		// Put a random file in the output.
		if err := os.WriteFile(filepath.Join(out, "existing.txt"), []byte("old"), 0644); err != nil {
			t.Fatal(err)
		}

		g := New(out, nil, Options{Overwrite: true})
		err := g.Generate("tauri-llm")
		if err != nil {
			t.Fatalf("overwrite should allow non-empty dir, error = %v", err)
		}
	})
}

func TestGenerateVariableSubstitution(t *testing.T) {
	// Create a temp template directory with a file containing placeholders.
	// We can't modify the embedded FS, but we can test substitution indirectly
	// by verifying that variables set on the generator are applied when
	// the generator processes file content through vars.Replace.
	//
	// The embedded placeholder.txt has no {{}} markers, so substitution
	// is a no-op — we verify the content is preserved.
	t.Run("no-op substitution on plain text", func(t *testing.T) {
		out := t.TempDir()
		vars := NewVariables()
		vars.Set("name", "my-project")
		g := New(out, vars, Options{})

		err := g.Generate("tauri-llm")
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		data, err := os.ReadFile(filepath.Join(out, "tauri-llm", "placeholder.txt"))
		if err != nil {
			t.Fatal(err)
		}
		content := string(data)
		if strings.Contains(content, "{{name}}") {
			t.Error("placeholder should not remain in output when no placeholders exist in source")
		}
		if !strings.Contains(content, "Placeholder template directory") {
			t.Errorf("original content lost: %q", content)
		}
	})
}

func TestGenerateExecutablePermissions(t *testing.T) {
	t.Run("non-executable file gets default perm", func(t *testing.T) {
		out := t.TempDir()
		g := New(out, nil, Options{})

		err := g.Generate("tauri-llm")
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		info, err := os.Stat(filepath.Join(out, "tauri-llm", "placeholder.txt"))
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != defaultFilePerm {
			t.Errorf("perm = %o, want %o", info.Mode().Perm(), defaultFilePerm)
		}
	})

	t.Run("isExecutable detects shell scripts", func(t *testing.T) {
		tests := []struct {
			path string
			want bool
		}{
			{"install.sh", true},
			{"setup.bat", true},
			{"run.cmd", true},
			{"build.ps1", true},
			{"main.go", false},
			{"README.md", false},
			{"Dockerfile", false},
			{"script.sh.bak", false},
			{"PATH/TO/install.sh", true},
		}
		for _, tt := range tests {
			got := isExecutable(tt.path)
			if got != tt.want {
				t.Errorf("isExecutable(%q) = %v, want %v", tt.path, got, tt.want)
			}
		}
	})
}

func TestGenerateDirectoryStructure(t *testing.T) {
	t.Run("creates nested directories", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "nested", "output")
		g := New(out, nil, Options{})

		err := g.Generate("tauri-llm")
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		// The template puts files under tauri-llm/ inside the output dir.
		info, err := os.Stat(filepath.Join(out, "tauri-llm"))
		if err != nil {
			t.Fatalf("subdirectory not created: %v", err)
		}
		if !info.IsDir() {
			t.Error("expected subdirectory, got file")
		}
	})

	t.Run("output directory has correct permissions", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "perm-check")
		g := New(out, nil, Options{})

		err := g.Generate("whatsapp-bot")
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		info, err := os.Stat(out)
		if err != nil {
			t.Fatal(err)
		}
		if !info.IsDir() {
			t.Error("output should be a directory")
		}
		// Directory should be executable (0755 allows traversal).
		if info.Mode().Perm()&0100 == 0 {
			t.Error("output directory should have execute bit set")
		}
	})
}

func TestGenerateAllTemplates(t *testing.T) {
	templates := []string{"tauri-llm", "whatsapp-bot", "expense-splitter"}

	for _, tmpl := range templates {
		t.Run(tmpl, func(t *testing.T) {
			out := t.TempDir()
			g := New(out, nil, Options{})

			err := g.Generate(tmpl)
			if err != nil {
				t.Fatalf("Generate(%q) error = %v", tmpl, err)
			}

			// Verify the output directory contains the template subdirectory.
			subDir := filepath.Join(out, tmpl)
			info, err := os.Stat(subDir)
			if err != nil {
				t.Fatalf("template subdirectory %q not found: %v", subDir, err)
			}
			if !info.IsDir() {
				t.Errorf("%q should be a directory", subDir)
			}

			// Verify at least one file was generated.
			entries, err := os.ReadDir(subDir)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) == 0 {
				t.Errorf("template %q produced no files", tmpl)
			}
		})
	}
}

func TestGenerateErrorWrapping(t *testing.T) {
	t.Run("template not found error wraps correctly", func(t *testing.T) {
		out := t.TempDir()
		g := New(out, nil, Options{})

		err := g.Generate("no-such-template")
		if err == nil {
			t.Fatal("expected error")
		}

		// The error chain should contain ErrTemplateNotFound.
		if !apperrors.IsTemplateNotFound(err) {
			t.Errorf("error chain should contain ErrTemplateNotFound, got: %v", err)
		}

		// The error message should mention the template name.
		if !strings.Contains(err.Error(), "no-such-template") {
			t.Errorf("error should mention template name, got: %v", err)
		}
	})

	t.Run("output exists error wraps correctly", func(t *testing.T) {
		out := t.TempDir()
		// Create a file so the dir is non-empty.
		os.WriteFile(filepath.Join(out, "x.txt"), []byte("x"), 0644)

		g := New(out, nil, Options{})
		err := g.Generate("tauri-llm")
		if err == nil {
			t.Fatal("expected error")
		}
		if !apperrors.IsOutputExists(err) {
			t.Errorf("expected ErrOutputExists, got: %v", err)
		}
		if !strings.Contains(err.Error(), out) {
			t.Errorf("error should mention output path, got: %v", err)
		}
	})
}

func TestIsExecutable(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"shell script", "script.sh", true},
		{"batch file", "run.bat", true},
		{"cmd file", "start.cmd", true},
		{"powershell", "deploy.ps1", true},
		{"go file", "main.go", false},
		{"markdown", "README.md", false},
		{"yaml", "config.yaml", false},
		{"no extension", "Makefile", false},
		{"nested shell", "scripts/install.sh", true},
		{"uppercase extension", "SCRIPT.SH", true},
		{"mixed case", "Run.Bat", true},
		{"sh inside path", "/usr/local/bin/app.sh", true},
		{"dotfile", ".bashrc", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isExecutable(tt.path)
			if got != tt.want {
				t.Errorf("isExecutable(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
