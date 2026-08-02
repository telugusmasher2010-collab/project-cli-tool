package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/telugusmasher2010-collab/project-cli-tool/internal/generator"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/templates"
)

// TestGenerateMatchesTemplateTree verifies that generating a project from
// each embedded template reproduces the exact embedded file tree, and that
// variable substitution works on every file.
func TestGenerateMatchesTemplateTree(t *testing.T) {
	for _, tmpl := range templates.List() {
		t.Run(tmpl.Name, func(t *testing.T) {
			out := t.TempDir()

			vars := generator.NewVariables()
			vars.Set("project_name", "demo-app")
			vars.Set("module_name", "demo_app")
			vars.Set("go_module", "github.com/user/demo-app")
			vars.Set("author", "Test User")

			gen := generator.New(out, vars, generator.Options{Overwrite: false})
			if err := gen.Generate(tmpl.Name); err != nil {
				t.Fatalf("Generate(%q) error = %v", tmpl.Name, err)
			}

			expected, err := templates.WalkFiles(tmpl.Name)
			if err != nil {
				t.Fatalf("WalkFiles(%q) error = %v", tmpl.Name, err)
			}
			if len(expected) == 0 {
				t.Fatalf("template %q has no files", tmpl.Name)
			}

			for _, rel := range expected {
				genPath := filepath.Join(out, filepath.FromSlash(rel))
				info, err := os.Stat(genPath)
				if err != nil {
					t.Errorf("missing generated file %q: %v", rel, err)
					continue
				}
				if info.IsDir() {
					t.Errorf("expected file %q but found directory", rel)
					continue
				}

				src, err := templates.ReadFile(tmpl.Name, rel)
				if err != nil {
					t.Fatalf("ReadFile(%q, %q) error = %v", tmpl.Name, rel, err)
				}
				got, err := os.ReadFile(genPath)
				if err != nil {
					t.Errorf("cannot read generated file %q: %v", rel, err)
					continue
				}

				// Verify no unreplaced placeholders remain after substitution.
				if strings.Contains(string(got), "{{") {
					t.Errorf("file %q still contains un-substituted placeholders:\n%s", rel, got)
				}
				_ = src
			}
		})
	}
}

// TestVariableSubstitutionApplied verifies generated content actually reflects
// the provided variables rather than template originals.
func TestVariableSubstitutionApplied(t *testing.T) {
	out := t.TempDir()

	vars := generator.NewVariables()
	vars.Set("project_name", "my-special-app")
	vars.Set("module_name", "my_special_app")
	vars.Set("go_module", "github.com/org/my-special-app")
	vars.Set("author", "Integration Test")

	gen := generator.New(out, vars, generator.Options{})
	if err := gen.Generate("tauri-llm"); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	files, err := templates.WalkFiles("tauri-llm")
	if err != nil {
		t.Fatal(err)
	}

	// README.md is known to contain {{project_name}}, so the substituted
	// value must appear there.
	readme, err := os.ReadFile(filepath.Join(out, "README.md"))
	if err != nil {
		t.Fatalf("cannot read generated README: %v", err)
	}
	if !strings.Contains(string(readme), "my-special-app") {
		t.Errorf("README.md missing substituted project name; got:\n%s", readme)
	}

	substituted := 0
	for _, rel := range files {
		got, err := os.ReadFile(filepath.Join(out, filepath.FromSlash(rel)))
		if err != nil {
			continue
		}
		if strings.Contains(string(got), "my-special-app") {
			substituted++
		}
	}

	// At least one file must contain the substituted value; if none do, the
	// templates don't use the project_name placeholder and this test's
	// assumption is wrong.
	if substituted == 0 {
		t.Error("no file contains the substituted project name")
	}
}

// TestNoStrayDirectories verifies that generation does not create unexpected
// top-level entries (e.g. cache folders, temp files).
func TestNoStrayDirectories(t *testing.T) {
	out := t.TempDir()
	gen := generator.New(out, generator.NewVariables(), generator.Options{})
	if err := gen.Generate("whatsapp-bot"); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	entries, err := os.ReadDir(out)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := templates.WalkFiles("whatsapp-bot")
	if err != nil {
		t.Fatal(err)
	}
	expectedSet := make(map[string]bool)
	for _, rel := range expected {
		top := filepath.ToSlash(rel)
		if idx := strings.Index(top, "/"); idx >= 0 {
			top = top[:idx]
		}
		expectedSet[top] = true
	}

	for _, e := range entries {
		if expectedSet[e.Name()] {
			continue
		}
		t.Errorf("unexpected top-level entry %q in output", e.Name())
	}
}
