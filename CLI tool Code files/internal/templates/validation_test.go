package templates

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"go.yaml.in/yaml/v3"
)

// manifestKind identifies the serialization format of a template manifest.
type manifestKind int

const (
	manifestJSON manifestKind = iota
	manifestYAML
	manifestTOML
)

// templateManifests lists the manifest files that must exist and parse for
// each registered template, keyed by template name. Every new template must
// be added here or TestTemplateTablesCoverAllRegisteredTemplates fails.
var templateManifests = map[string][]struct {
	path string
	kind manifestKind
}{
	"tauri-llm": {
		{"package.json", manifestJSON},
		{"src-tauri/tauri.conf.json", manifestJSON},
		{"Cargo.toml", manifestTOML},
		{"src-tauri/Cargo.toml", manifestTOML},
	},
	"whatsapp-bot": {
		{"package.json", manifestJSON},
		{"config.example.json", manifestJSON},
	},
	"expense-splitter": {
		{"pubspec.yaml", manifestYAML},
	},
	"next-webapp": {
		{"package.json", manifestJSON},
		{"tsconfig.json", manifestJSON},
	},
	"react-native-map": {
		{"package.json", manifestJSON},
		{"app.json", manifestJSON},
		{"tsconfig.json", manifestJSON},
	},
}

// templateRequiredFiles lists the files every generated project of a given
// stack must contain. It is the structural contract for a usable template.
// Note: go:embed excludes dotfiles (.gitignore) and underscore-prefixed files
// (_layout.tsx) by design, so those cannot be asserted here.
var templateRequiredFiles = map[string][]string{
	"tauri-llm":        {"README.md", "package.json", "Cargo.toml", "src-tauri/Cargo.toml", "src-tauri/tauri.conf.json", "src/main.rs"},
	"whatsapp-bot":     {"README.md", "package.json", "config.example.json", "index.js", "database/schema.sql"},
	"expense-splitter": {"README.md", "pubspec.yaml", "lib/main.dart", "lib/app.dart", "supabase/schema.sql"},
	"next-webapp":      {"README.md", "package.json", "tsconfig.json", "app/layout.tsx", "app/page.tsx"},
	"react-native-map": {"README.md", "package.json", "app.json", "tsconfig.json", "app/index.tsx"},
}

// allowedPlaceholders is the set of variable keys the CLI engine is
// guaranteed to populate. Templates must not introduce new keys, otherwise
// generated projects would contain un-substituted placeholders.
var allowedPlaceholders = map[string]bool{
	"project_name": true,
	"module_name":  true,
	"author":       true,
	"go_module":    true,
}

var placeholderRe = regexp.MustCompile(`\{\{(\w+)\}\}`)

func templateFileSet(t *testing.T, name string) []string {
	t.Helper()
	files, err := WalkFiles(name)
	if err != nil {
		t.Fatalf("WalkFiles(%q) error = %v", name, err)
	}
	return files
}

// TestTemplateTablesCoverAllRegisteredTemplates ensures the validation tables
// below stay in sync with the registry: a new template must be documented
// before it is considered validated.
func TestTemplateTablesCoverAllRegisteredTemplates(t *testing.T) {
	for _, info := range List() {
		if _, ok := templateRequiredFiles[info.Name]; !ok {
			t.Errorf("template %q has no required-files entry; add one to templateRequiredFiles", info.Name)
		}
		if _, ok := templateManifests[info.Name]; !ok {
			t.Errorf("template %q has no manifest entry; add one to templateManifests", info.Name)
		}
	}
}

// TestTemplatesHaveRequiredFiles verifies every template ships its mandatory
// files and that none of them are empty.
func TestTemplatesHaveRequiredFiles(t *testing.T) {
	for name, required := range templateRequiredFiles {
		t.Run(name, func(t *testing.T) {
			have := make(map[string]bool)
			for _, f := range templateFileSet(t, name) {
				have[f] = true
			}
			for _, req := range required {
				if !have[req] {
					t.Errorf("missing required file %q", req)
					continue
				}
				data, err := ReadFile(name, req)
				if err != nil {
					t.Errorf("cannot read required file %q: %v", req, err)
					continue
				}
				if len(strings.TrimSpace(string(data))) == 0 {
					t.Errorf("required file %q is empty", req)
				}
			}
		})
	}
}

// TestEmbeddedManifestsAreStructurallyValid parses every manifest file with
// the appropriate parser (JSON, YAML or TOML) to prove the embedded templates
// are well-formed and would survive downstream tooling.
func TestEmbeddedManifestsAreStructurallyValid(t *testing.T) {
	for name, manifests := range templateManifests {
		t.Run(name, func(t *testing.T) {
			have := make(map[string]bool)
			for _, f := range templateFileSet(t, name) {
				have[f] = true
			}
			for _, m := range manifests {
				if !have[m.path] {
					t.Errorf("expected manifest %q to exist", m.path)
					continue
				}
				data, err := ReadFile(name, m.path)
				if err != nil {
					t.Errorf("cannot read manifest %q: %v", m.path, err)
					continue
				}
				// Placeholders make YAML flow syntax ambiguous: yaml.v3 parses
				// "{{module_name}}" as a nested flow mapping. Normalize them so
				// the rest of the document is parsed as it will exist after
				// substitution by the generator.
				normalized := placeholderRe.ReplaceAllString(string(data), "x")
				var out map[string]any
				switch m.kind {
				case manifestJSON:
					err = json.Unmarshal([]byte(normalized), &out)
				case manifestYAML:
					err = yaml.Unmarshal([]byte(normalized), &out)
				case manifestTOML:
					err = toml.Unmarshal([]byte(normalized), &out)
				default:
					t.Errorf("manifest %q: unknown kind %d", m.path, m.kind)
					continue
				}
				if err != nil {
					t.Errorf("manifest %q is not structurally valid: %v", m.path, err)
					continue
				}
				if len(out) == 0 {
					t.Errorf("manifest %q parsed to an empty document", m.path)
				}
			}
		})
	}
}

// TestTemplatesUseKnownPlaceholders verifies the placeholder contract: every
// {{key}} used by a template must be a key the CLI engine populates, so no
// un-substituted placeholders can survive generation.
func TestTemplatesUseKnownPlaceholders(t *testing.T) {
	for _, info := range List() {
		t.Run(info.Name, func(t *testing.T) {
			seen := make(map[string]bool)
			for _, rel := range templateFileSet(t, info.Name) {
				data, err := ReadFile(info.Name, rel)
				if err != nil {
					t.Errorf("cannot read %q: %v", rel, err)
					continue
				}
				for _, match := range placeholderRe.FindAllStringSubmatch(string(data), -1) {
					key := match[1]
					seen[key] = true
					if !allowedPlaceholders[key] {
						t.Errorf("file %q uses unsupported placeholder {{%s}}", rel, key)
					}
				}
			}
			if len(seen) == 0 {
				t.Error("template uses no placeholders; expected the placeholder contract to be exercised")
			}
		})
	}
}

// TestEveryTemplateUsesProjectName verifies each template consumes the core
// variable that cmd/init.go always supplies, so generation always has an
// observable effect on the scaffolded project.
func TestEveryTemplateUsesProjectName(t *testing.T) {
	for _, info := range List() {
		t.Run(info.Name, func(t *testing.T) {
			found := false
			for _, rel := range templateFileSet(t, info.Name) {
				data, err := ReadFile(info.Name, rel)
				if err != nil {
					t.Errorf("cannot read %q: %v", rel, err)
					continue
				}
				if strings.Contains(string(data), "{{project_name}}") {
					found = true
					break
				}
			}
			if !found {
				t.Error("template never uses {{project_name}}; the CLI always supplies it")
			}
		})
	}
}
