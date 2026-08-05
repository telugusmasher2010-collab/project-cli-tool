package integration

import (
	"bytes"
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	apperrors "github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/generator"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/templates"
)

// dogfoodVars returns the variable set used for dogfood generation. The
// author value deliberately includes characters (&, spaces) that must not
// break literal substitution.
func dogfoodVars() *generator.Variables {
	v := generator.NewVariables()
	v.Set("project_name", "dogfood-app")
	v.Set("module_name", "dogfood_app")
	v.Set("go_module", "github.com/example/dogfood-app")
	v.Set("author", "Dogfood Tester & Co")
	return v
}

func dogfoodGenerate(t *testing.T, tmpl string, opts generator.Options) string {
	t.Helper()
	out := t.TempDir()
	if err := generator.New(out, dogfoodVars(), opts).Generate(tmpl); err != nil {
		t.Fatalf("Generate(%q) error = %v", tmpl, err)
	}
	return out
}

func dogfoodRead(t *testing.T, out, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(out, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatalf("missing generated file %q: %v", rel, err)
	}
	return string(data)
}

// TestDogfoodGenerateEveryTemplate generates a real project from every
// embedded template and verifies generation succeeds, the file structure
// matches the embedded template exactly, and every placeholder is replaced.
func TestDogfoodGenerateEveryTemplate(t *testing.T) {
	for _, tmpl := range templates.List() {
		t.Run(tmpl.Name, func(t *testing.T) {
			out := dogfoodGenerate(t, tmpl.Name, generator.Options{})

			expected, err := templates.WalkFiles(tmpl.Name)
			if err != nil {
				t.Fatalf("WalkFiles(%q) error = %v", tmpl.Name, err)
			}
			expectedSet := make(map[string]bool, len(expected))
			for _, rel := range expected {
				expectedSet[rel] = true
			}

			// 1. Structure: exactly the embedded tree, nothing more or less.
			got := make(map[string]bool)
			err = filepath.WalkDir(out, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				rel, err := filepath.Rel(out, path)
				if err != nil {
					return err
				}
				got[filepath.ToSlash(rel)] = true
				return nil
			})
			if err != nil {
				t.Fatalf("cannot walk generated output: %v", err)
			}

			for _, rel := range expected {
				if !got[rel] {
					t.Errorf("missing generated file %q", rel)
				}
				if !strings.HasPrefix(rel, ".gitignore") && strings.TrimSpace(dogfoodRead(t, out, rel)) == "" {
					t.Errorf("generated file %q is empty", rel)
				}
			}
			for rel := range got {
				if !expectedSet[rel] {
					t.Errorf("unexpected generated file %q", rel)
				}
			}

			// 2. Placeholders: nothing may survive substitution.
			for _, rel := range expected {
				if content := dogfoodRead(t, out, rel); strings.Contains(content, "{{") {
					t.Errorf("file %q still contains un-substituted placeholders", rel)
				}
			}

			// 3. Dotfile and underscore-file embed regression guards.
			if _, err := os.Stat(filepath.Join(out, ".gitignore")); err != nil {
				t.Errorf("generated project missing .gitignore: %v", err)
			}

			// 4. Substitution correctness on the known manifest locations.
			switch tmpl.Name {
			case "tauri-llm":
				if !strings.Contains(dogfoodRead(t, out, "README.md"), "dogfood-app") {
					t.Error("README.md missing substituted project_name")
				}
				if !strings.Contains(dogfoodRead(t, out, "package.json"), `"name": "dogfood-app"`) {
					t.Error("package.json name not substituted with project_name")
				}
				if !strings.Contains(dogfoodRead(t, out, "src-tauri/Cargo.toml"), `name = "dogfood_app-tauri"`) {
					t.Error("src-tauri/Cargo.toml name not substituted with module_name")
				}
			case "whatsapp-bot":
				if !strings.Contains(dogfoodRead(t, out, "README.md"), "dogfood-app") {
					t.Error("README.md missing substituted project_name")
				}
				if !strings.Contains(dogfoodRead(t, out, "package.json"), `"name": "dogfood_app"`) {
					t.Error("package.json name not substituted with module_name")
				}
				if !strings.Contains(dogfoodRead(t, out, "config.example.json"), `"botName": "dogfood-app"`) {
					t.Error("config.example.json botName not substituted with project_name")
				}
			case "expense-splitter":
				if !strings.Contains(dogfoodRead(t, out, "README.md"), "dogfood-app") {
					t.Error("README.md missing substituted project_name")
				}
				if !strings.Contains(dogfoodRead(t, out, "pubspec.yaml"), "name: dogfood_app") {
					t.Error("pubspec.yaml name not substituted with module_name")
				}
				if !strings.Contains(dogfoodRead(t, out, "lib/core/constants.dart"), "'dogfood-app'") {
					t.Error("constants.dart appName not substituted with project_name")
				}
			case "next-webapp":
				if !strings.Contains(dogfoodRead(t, out, "README.md"), "dogfood-app") {
					t.Error("README.md missing substituted project_name")
				}
				if !strings.Contains(dogfoodRead(t, out, "package.json"), `"name": "dogfood_app"`) {
					t.Error("package.json name not substituted with module_name")
				}
				if !strings.Contains(dogfoodRead(t, out, "app/layout.tsx"), `title: "dogfood-app"`) {
					t.Error("layout.tsx metadata.title not substituted with project_name")
				}
			case "react-native-map":
				if !strings.Contains(dogfoodRead(t, out, "README.md"), "dogfood-app") {
					t.Error("README.md missing substituted project_name")
				}
				if !strings.Contains(dogfoodRead(t, out, "app.json"), `"slug": "dogfood_app"`) {
					t.Error("app.json slug not substituted with module_name")
				}
				if !strings.Contains(dogfoodRead(t, out, "app/index.tsx"), "dogfood-app") {
					t.Error("app/index.tsx missing substituted project_name")
				}
				if _, err := os.Stat(filepath.Join(out, "app", "_layout.tsx")); err != nil {
					t.Errorf("generated project missing app/_layout.tsx: %v", err)
				}
			}
		})
	}
}

// TestDogfoodAuthorSubstitution verifies values containing shell and regex
// metacharacters are substituted literally in the README license line.
func TestDogfoodAuthorSubstitution(t *testing.T) {
	for _, tmpl := range templates.List() {
		t.Run(tmpl.Name, func(t *testing.T) {
			out := dogfoodGenerate(t, tmpl.Name, generator.Options{})
			readme := dogfoodRead(t, out, "README.md")
			if tmpl.Name == "tauri-llm" {
				// tauri-llm README does not carry a license line; author lives
				// in Cargo.toml instead.
				cargo := dogfoodRead(t, out, "Cargo.toml")
				if !strings.Contains(cargo, `authors = ["Dogfood Tester & Co"]`) {
					t.Errorf("Cargo.toml authors not substituted literally; got:\n%s", cargo)
				}
				return
			}
			if !strings.Contains(readme, "MIT — Dogfood Tester & Co") {
				t.Errorf("README license line not substituted literally; got:\n%s", readme)
			}
		})
	}
}

// TestDogfoodHookSelection verifies the default post-generation hook set is
// derived from the template manifests as expected.
func TestDogfoodHookSelection(t *testing.T) {
	expected := map[string][]string{
		"tauri-llm":        {"git init", "npm install"},
		"whatsapp-bot":     {"git init", "npm install"},
		"expense-splitter": {"git init", "flutter pub get"},
		"next-webapp":      {"git init", "npm install"},
		"react-native-map": {"git init", "npm install"},
	}
	for _, tmpl := range templates.List() {
		t.Run(tmpl.Name, func(t *testing.T) {
			got := generator.HooksForTemplate(tmpl.Name).Names()
			if !reflect.DeepEqual(got, expected[tmpl.Name]) {
				t.Errorf("HooksForTemplate(%q) = %v, want %v", tmpl.Name, got, expected[tmpl.Name])
			}
		})
	}
}

// TestDogfoodHooksRunInGeneratedProject verifies hooks execute against the
// generated directory and that git init actually initializes a repository.
func TestDogfoodHooksRunInGeneratedProject(t *testing.T) {
	out := dogfoodGenerate(t, "whatsapp-bot", generator.Options{})

	hr := generator.NewHookRunner(
		generator.GitInitHook(),
		&generator.FuncHook{
			Name_: "marker",
			F: func(ctx context.Context, dir string) error {
				return os.WriteFile(filepath.Join(dir, "hook-ran.txt"), []byte("ran"), 0o644)
			},
		},
	)
	if err := hr.RunAll(context.Background(), out); err != nil {
		t.Fatalf("RunAll error = %v", err)
	}

	if info, err := os.Stat(filepath.Join(out, ".git")); err != nil || !info.IsDir() {
		t.Errorf("git init hook did not create .git directory: %v", err)
	}
	if _, err := os.Stat(filepath.Join(out, "hook-ran.txt")); err != nil {
		t.Errorf("marker hook did not run in project dir: %v", err)
	}
}

// TestDogfoodHooksSkippedWhenNotRequested verifies that with AutoHooks off
// and no explicit hooks, no post-generation action runs.
func TestDogfoodHooksSkippedWhenNotRequested(t *testing.T) {
	out := dogfoodGenerate(t, "tauri-llm", generator.Options{})

	for _, name := range []string{".git", "node_modules", "package-lock.json"} {
		if _, err := os.Stat(filepath.Join(out, name)); !os.IsNotExist(err) {
			t.Errorf("expected %q to be absent without hooks, err = %v", name, err)
		}
	}
}

// TestDogfoodMissingToolHookFailsLoudly verifies a hook whose executable is
// unavailable returns an ErrHookFailed error instead of being silently
// skipped, so users learn why a post-generation step did not run.
func TestDogfoodMissingToolHookFailsLoudly(t *testing.T) {
	out := dogfoodGenerate(t, "whatsapp-bot", generator.Options{})

	h := &generator.CommandHook{
		Name_:   "missing tool",
		Command: "projinit-tool-that-does-not-exist-xyz",
	}
	err := h.Run(context.Background(), out)
	if err == nil {
		t.Fatal("expected error for missing executable")
	}
	if !apperrors.IsCode(err, apperrors.ErrHookFailed) {
		t.Errorf("expected ErrHookFailed, got: %v", err)
	}
}

// TestDogfoodRealHooks generates each template and runs the default hook set
// for real. It is opt-in because npm install and flutter pub get require
// network access and the corresponding toolchains. git init must always
// succeed; tool-dependent hooks are reported via logs and their failure is
// expected to be loud (ErrHookFailed) rather than silent. Run with
// DOGFOOD_REAL_HOOKS=1.
func TestDogfoodRealHooks(t *testing.T) {
	if os.Getenv("DOGFOOD_REAL_HOOKS") != "1" {
		t.Skip("set DOGFOOD_REAL_HOOKS=1 to run real post-generation hooks")
	}

	factories := map[string]func() generator.Hook{
		"git init":        generator.GitInitHook,
		"npm install":     generator.NpmInstallHook,
		"flutter pub get": generator.FlutterPubGetHook,
		"go mod tidy":     generator.GoModTidyHook,
	}

	for _, tmpl := range templates.List() {
		t.Run(tmpl.Name, func(t *testing.T) {
			out := dogfoodGenerate(t, tmpl.Name, generator.Options{})

			for _, name := range generator.HooksForTemplate(tmpl.Name).Names() {
				t.Run(name, func(t *testing.T) {
					var stdout, stderr bytes.Buffer
					h := factories[name]()
					if c, ok := h.(*generator.CommandHook); ok {
						c.Stdout = &stdout
						c.Stderr = &stderr
					}
					err := h.Run(context.Background(), out)
					if name == "git init" {
						if err != nil {
							t.Fatalf("git init must always succeed: %v", err)
						}
						if info, e := os.Stat(filepath.Join(out, ".git")); e != nil || !info.IsDir() {
							t.Fatalf("git init hook did not create .git: %v", e)
						}
						return
					}
					t.Logf("%s stdout:\n%s\nstderr:\n%s", name, stdout.String(), stderr.String())
					if err != nil {
						if !apperrors.IsCode(err, apperrors.ErrHookFailed) {
							t.Errorf("hook %q failed with a non-hook error: %v", name, err)
						} else {
							t.Logf("hook %q failed loudly (expected if the tool is unavailable or network/registry access is blocked): %v", name, err)
						}
						return
					}
					if name == "npm install" {
						if _, e := os.Stat(filepath.Join(out, "node_modules")); e != nil {
							t.Errorf("npm install succeeded but node_modules is missing: %v", e)
						}
					}
				})
			}
		})
	}
}
