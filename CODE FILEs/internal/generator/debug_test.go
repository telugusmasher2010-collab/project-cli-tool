package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/telugusmasher2010-collab/project-cli-tool/internal/templates"
)

func TestDebugWalkFiles(t *testing.T) {
	files, err := templates.WalkFiles("tauri-llm")
	if err != nil {
		t.Fatalf("WalkFiles error: %v", err)
	}
	fmt.Println("WalkFiles returned:")
	for _, f := range files {
		fmt.Printf("  %q\n", f)
	}

	out := t.TempDir()
	g := New(out, nil, Options{})
	err = g.Generate("tauri-llm")
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	fmt.Println("\nGenerated files:")
	filepath.Walk(out, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(out, path)
		fmt.Printf("  %s (dir=%v, size=%d, mode=%o)\n", rel, info.IsDir(), info.Size(), info.Mode().Perm())
		return nil
	})
}
