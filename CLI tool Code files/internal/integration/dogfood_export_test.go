package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/telugusmasher2010-collab/project-cli-tool/internal/generator"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/templates"
)

// TestDogfoodExport writes generated projects to a fixed directory so the
// dogfooding session can exercise the real toolchains (npm, node) against
// them. It is a manual aid and not part of CI.
func TestDogfoodExport(t *testing.T) {
	root := os.Getenv("DOGFOOD_EXPORT")
	if root == "" {
		t.Skip("set DOGFOOD_EXPORT to a directory to export generated projects")
	}
	for _, tmpl := range templates.List() {
		out := filepath.Join(root, tmpl.Name)
		if err := os.RemoveAll(out); err != nil {
			t.Fatal(err)
		}
		vars := dogfoodVars()
		if err := generator.New(out, vars, generator.Options{}).Generate(tmpl.Name); err != nil {
			t.Fatalf("Generate(%q) error = %v", tmpl.Name, err)
		}
	}
}
