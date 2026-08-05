package templates

import (
	"io/fs"
	"sort"
	"strings"
	"testing"

	apperrors "github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

// TestRegistryMatchesEmbeddedFilesystem verifies that the static registry in
// embed.go agrees exactly with the directories embedded via TemplateFS. A
// template that is registered but not embedded (or vice versa) would break
// generation at runtime.
func TestRegistryMatchesEmbeddedFilesystem(t *testing.T) {
	embedded, err := fs.ReadDir(TemplateFS, ".")
	if err != nil {
		t.Fatalf("cannot list embedded templates: %v", err)
	}

	registered := make(map[string]bool, len(registry))
	for _, info := range registry {
		if registered[info.Name] {
			t.Errorf("duplicate registry entry for template %q", info.Name)
		}
		registered[info.Name] = true
	}

	// Every registered template must exist as a real directory in the embed
	// filesystem, and Name must equal Directory.
	for _, info := range registry {
		entry, err := fs.Stat(TemplateFS, info.Directory)
		if err != nil {
			t.Errorf("template %q: embedded directory %q missing: %v", info.Name, info.Directory, err)
			continue
		}
		if !entry.IsDir() {
			t.Errorf("template %q: embedded path %q is not a directory", info.Name, info.Directory)
		}
		if info.Name != info.Directory {
			t.Errorf("template %q: Name and Directory must match, got %q", info.Name, info.Directory)
		}
	}

	// Every embedded top-level directory must be registered.
	for _, entry := range embedded {
		if !entry.IsDir() {
			t.Errorf("unexpected non-directory entry %q at embedded filesystem root", entry.Name())
			continue
		}
		if !registered[entry.Name()] {
			t.Errorf("embedded directory %q is not registered", entry.Name())
		}
	}

	if len(registry) == 0 {
		t.Fatal("registry is empty")
	}
}

func TestList(t *testing.T) {
	t.Run("returns all registered templates sorted by name", func(t *testing.T) {
		got := List()
		if len(got) != len(registry) {
			t.Fatalf("List() returned %d templates, want %d", len(got), len(registry))
		}
		if !sort.SliceIsSorted(got, func(i, j int) bool { return got[i].Name < got[j].Name }) {
			t.Error("List() result is not sorted by name")
		}
		for _, info := range registry {
			found := false
			for _, item := range got {
				if item.Name == info.Name {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("List() missing template %q", info.Name)
			}
		}
	})

	t.Run("returns copies, not references", func(t *testing.T) {
		first := List()
		second := List()
		if len(first) == 0 {
			t.Fatal("List() is empty")
		}
		first[0].Description = "mutated"
		if second[0].Description == "mutated" {
			t.Error("List() returned a reference, not a copy")
		}
	})
}

func TestGet(t *testing.T) {
	for _, info := range registry {
		t.Run(info.Name, func(t *testing.T) {
			got, err := Get(info.Name)
			if err != nil {
				t.Fatalf("Get(%q) error = %v", info.Name, err)
			}
			if got.Name != info.Name {
				t.Errorf("Get(%q).Name = %q, want %q", info.Name, got.Name, info.Name)
			}
			if got.Directory != info.Directory {
				t.Errorf("Get(%q).Directory = %q, want %q", info.Name, got.Directory, info.Directory)
			}
			if got.Description == "" {
				t.Errorf("Get(%q).Description must not be empty", info.Name)
			}
		})
	}

	t.Run("unknown template", func(t *testing.T) {
		_, err := Get("no-such-template")
		if err == nil {
			t.Fatal("expected error for unknown template")
		}
		if !apperrors.IsCode(err, apperrors.ErrTemplateNotFound) {
			t.Errorf("expected ErrTemplateNotFound, got: %v", err)
		}
	})
}

func TestReadFile(t *testing.T) {
	t.Run("reads a known file", func(t *testing.T) {
		data, err := ReadFile("tauri-llm", "README.md")
		if err != nil {
			t.Fatalf("ReadFile() error = %v", err)
		}
		if len(data) == 0 {
			t.Error("expected non-empty content")
		}
	})

	t.Run("unknown file in known template", func(t *testing.T) {
		_, err := ReadFile("tauri-llm", "missing.txt")
		if err == nil {
			t.Fatal("expected error for missing file")
		}
		if !apperrors.IsCode(err, apperrors.ErrTemplateNotFound) {
			t.Errorf("expected ErrTemplateNotFound, got: %v", err)
		}
	})

	t.Run("unknown template", func(t *testing.T) {
		_, err := ReadFile("no-such-template", "README.md")
		if err == nil {
			t.Fatal("expected error for unknown template")
		}
		if !apperrors.IsCode(err, apperrors.ErrTemplateNotFound) {
			t.Errorf("expected ErrTemplateNotFound, got: %v", err)
		}
	})
}

func TestWalkFiles(t *testing.T) {
	for _, info := range registry {
		t.Run(info.Name, func(t *testing.T) {
			files, err := WalkFiles(info.Name)
			if err != nil {
				t.Fatalf("WalkFiles(%q) error = %v", info.Name, err)
			}
			if len(files) == 0 {
				t.Fatal("template has no files")
			}
			if !sort.StringsAreSorted(files) {
				t.Error("WalkFiles() result is not sorted")
			}
			for _, rel := range files {
				if rel == "" || strings.HasPrefix(rel, "/") {
					t.Errorf("unexpected relative path %q", rel)
				}
				entry, err := fs.Stat(TemplateFS, info.Directory+"/"+rel)
				if err != nil {
					t.Errorf("path %q not found in embedded filesystem: %v", rel, err)
					continue
				}
				if entry.IsDir() {
					t.Errorf("WalkFiles returned directory %q", rel)
					continue
				}
				if _, err := ReadFile(info.Name, rel); err != nil {
					t.Errorf("ReadFile(%q) failed for walked file: %v", rel, err)
				}
			}

			// WalkFiles must enumerate exactly the files under the template
			// directory (excluding directories), in the same order.
			var expected []string
			err = fs.WalkDir(TemplateFS, info.Directory, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				expected = append(expected, strings.TrimPrefix(path, info.Directory+"/"))
				return nil
			})
			if err != nil {
				t.Fatalf("cannot walk embedded filesystem: %v", err)
			}
			sort.Strings(expected)
			if strings.Join(files, ",") != strings.Join(expected, ",") {
				t.Errorf("WalkFiles mismatch:\n got %v\nwant %v", files, expected)
			}
		})
	}

	t.Run("unknown template", func(t *testing.T) {
		_, err := WalkFiles("no-such-template")
		if err == nil {
			t.Fatal("expected error for unknown template")
		}
		if !apperrors.IsCode(err, apperrors.ErrTemplateNotFound) {
			t.Errorf("expected ErrTemplateNotFound, got: %v", err)
		}
	})
}

func TestExists(t *testing.T) {
	for _, info := range registry {
		if !Exists(info.Name) {
			t.Errorf("Exists(%q) = false, want true", info.Name)
		}
	}
	if Exists("no-such-template") {
		t.Error("Exists(unknown) = true, want false")
	}
}
