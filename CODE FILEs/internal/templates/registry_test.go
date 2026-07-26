package templates

import (
	"sync"
	"testing"

	"github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry returned nil")
	}
	if len(r.templates) != 0 {
		t.Errorf("expected empty registry, got %d templates", len(r.templates))
	}
}

func TestRegister(t *testing.T) {
	t.Run("registers a valid template", func(t *testing.T) {
		r := NewRegistry()
		tmpl := Template{
			Name:               "tauri-llm",
			Description:        "Tauri v2 + Rust + React + local LLM sidecar",
			Directory:          "tauri-llm",
			Category:           CategoryDesktop,
			SupportedLanguages: []string{"Rust", "TypeScript"},
			Version:            "1.0.0",
		}
		err := r.Register(tmpl)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !r.Exists("tauri-llm") {
			t.Error("template should exist after registration")
		}
	})

	t.Run("rejects duplicate registration", func(t *testing.T) {
		r := NewRegistry()
		tmpl := Template{
			Name:        "whatsapp-bot",
			Description: "Node.js Baileys bot",
			Directory:   "whatsapp-bot",
			Version:     "1.0.0",
		}
		if err := r.Register(tmpl); err != nil {
			t.Fatalf("first register failed: %v", err)
		}
		err := r.Register(tmpl)
		if err == nil {
			t.Fatal("expected error for duplicate registration")
		}
		if !errors.IsCode(err, errors.ErrTemplateExists) {
			t.Errorf("expected ErrTemplateExists, got: %v", err)
		}
	})

	t.Run("rejects empty name", func(t *testing.T) {
		r := NewRegistry()
		err := r.Register(Template{
			Name:        "",
			Description: "desc",
			Directory:   "dir",
			Version:     "1.0.0",
		})
		if err == nil {
			t.Fatal("expected error for empty name")
		}
		if !errors.IsCode(err, errors.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got: %v", err)
		}
	})

	t.Run("rejects empty description", func(t *testing.T) {
		r := NewRegistry()
		err := r.Register(Template{
			Name:        "valid",
			Description: "",
			Directory:   "dir",
			Version:     "1.0.0",
		})
		if err == nil {
			t.Fatal("expected error for empty description")
		}
		if !errors.IsCode(err, errors.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got: %v", err)
		}
	})

	t.Run("rejects empty directory", func(t *testing.T) {
		r := NewRegistry()
		err := r.Register(Template{
			Name:        "valid",
			Description: "desc",
			Directory:   "",
			Version:     "1.0.0",
		})
		if err == nil {
			t.Fatal("expected error for empty directory")
		}
		if !errors.IsCode(err, errors.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got: %v", err)
		}
	})

	t.Run("rejects empty version", func(t *testing.T) {
		r := NewRegistry()
		err := r.Register(Template{
			Name:        "valid",
			Description: "desc",
			Directory:   "dir",
			Version:     "",
		})
		if err == nil {
			t.Fatal("expected error for empty version")
		}
		if !errors.IsCode(err, errors.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got: %v", err)
		}
	})

	t.Run("allows empty supported languages", func(t *testing.T) {
		r := NewRegistry()
		err := r.Register(Template{
			Name:        "minimal",
			Description: "minimal template",
			Directory:   "minimal",
			Version:     "0.1.0",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("registers multiple distinct templates", func(t *testing.T) {
		r := NewRegistry()
		templates := []Template{
			{Name: "tauri-llm", Description: "Tauri desktop", Directory: "tauri-llm", Version: "1.0.0"},
			{Name: "whatsapp-bot", Description: "WhatsApp bot", Directory: "whatsapp-bot", Version: "1.0.0"},
			{Name: "expense-splitter", Description: "Flutter expense app", Directory: "expense-splitter", Version: "1.0.0"},
		}
		for _, tmpl := range templates {
			if err := r.Register(tmpl); err != nil {
				t.Fatalf("failed to register %q: %v", tmpl.Name, err)
			}
		}
		if len(r.templates) != 3 {
			t.Errorf("expected 3 templates, got %d", len(r.templates))
		}
	})

	t.Run("all categories work", func(t *testing.T) {
		r := NewRegistry()
		categories := []Category{CategoryCLI, CategoryWeb, CategoryMobile, CategoryDesktop, CategoryAI}
		for i, cat := range categories {
			name := string(rune('a' + i))
			err := r.Register(Template{
				Name:        name,
				Description: "test",
				Directory:   name,
				Category:    cat,
				Version:     "1.0.0",
			})
			if err != nil {
				t.Fatalf("failed to register with category %q: %v", cat, err)
			}
		}
		if len(r.templates) != 5 {
			t.Errorf("expected 5 templates, got %d", len(r.templates))
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("returns registered template", func(t *testing.T) {
		r := NewRegistry()
		want := Template{
			Name:               "tauri-llm",
			Description:        "Tauri v2 + Rust + React",
			Directory:          "tauri-llm",
			Category:           CategoryDesktop,
			SupportedLanguages: []string{"Rust", "TypeScript"},
			Version:            "1.0.0",
		}
		if err := r.Register(want); err != nil {
			t.Fatalf("register failed: %v", err)
		}
		got, err := r.Get("tauri-llm")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != want.Name {
			t.Errorf("Name = %q, want %q", got.Name, want.Name)
		}
		if got.Description != want.Description {
			t.Errorf("Description = %q, want %q", got.Description, want.Description)
		}
		if got.Directory != want.Directory {
			t.Errorf("Directory = %q, want %q", got.Directory, want.Directory)
		}
		if got.Category != want.Category {
			t.Errorf("Category = %q, want %q", got.Category, want.Category)
		}
		if got.Version != want.Version {
			t.Errorf("Version = %q, want %q", got.Version, want.Version)
		}
		if len(got.SupportedLanguages) != len(want.SupportedLanguages) {
			t.Errorf("SupportedLanguages length = %d, want %d", len(got.SupportedLanguages), len(want.SupportedLanguages))
		}
	})

	t.Run("returns error for missing template", func(t *testing.T) {
		r := NewRegistry()
		_, err := r.Get("nonexistent")
		if err == nil {
			t.Fatal("expected error for missing template")
		}
		if !errors.IsCode(err, errors.ErrTemplateNotFound) {
			t.Errorf("expected ErrTemplateNotFound, got: %v", err)
		}
	})

	t.Run("error message includes template name", func(t *testing.T) {
		r := NewRegistry()
		_, err := r.Get("my-missing-template")
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() == "" {
			t.Error("error message should not be empty")
		}
	})
}

func TestList(t *testing.T) {
	t.Run("returns empty list for empty registry", func(t *testing.T) {
		r := NewRegistry()
		list := r.List()
		if len(list) != 0 {
			t.Errorf("expected empty list, got %d items", len(list))
		}
	})

	t.Run("returns all registered templates sorted by name", func(t *testing.T) {
		r := NewRegistry()
		names := []string{"expense-splitter", "tauri-llm", "whatsapp-bot"}
		for _, name := range names {
			if err := r.Register(Template{
				Name:        name,
				Description: name + " template",
				Directory:   name,
				Version:     "1.0.0",
			}); err != nil {
				t.Fatalf("failed to register %q: %v", name, err)
			}
		}
		list := r.List()
		if len(list) != 3 {
			t.Fatalf("expected 3 templates, got %d", len(list))
		}
		// Should be sorted: expense-splitter, tauri-llm, whatsapp-bot
		sorted := []string{"expense-splitter", "tauri-llm", "whatsapp-bot"}
		for i, want := range sorted {
			if list[i].Name != want {
				t.Errorf("list[%d].Name = %q, want %q", i, list[i].Name, want)
			}
		}
	})

	t.Run("returns copies not references", func(t *testing.T) {
		r := NewRegistry()
		if err := r.Register(Template{
			Name:        "original",
			Description: "original desc",
			Directory:   "original",
			Version:     "1.0.0",
		}); err != nil {
			t.Fatalf("register failed: %v", err)
		}
		list1 := r.List()
		list2 := r.List()
		// Mutate the first list's element
		list1[0].Description = "mutated"
		// Second list should be unaffected
		if list2[0].Description != "original desc" {
			t.Error("List() returned a reference, not a copy")
		}
	})
}

func TestExists(t *testing.T) {
	t.Run("returns true for registered template", func(t *testing.T) {
		r := NewRegistry()
		if err := r.Register(Template{
			Name:        "test-template",
			Description: "test",
			Directory:   "test-template",
			Version:     "1.0.0",
		}); err != nil {
			t.Fatalf("register failed: %v", err)
		}
		if !r.Exists("test-template") {
			t.Error("Exists should return true for registered template")
		}
	})

	t.Run("returns false for missing template", func(t *testing.T) {
		r := NewRegistry()
		if r.Exists("nonexistent") {
			t.Error("Exists should return false for missing template")
		}
	})

	t.Run("returns false after unregistering by creating new registry", func(t *testing.T) {
		r := NewRegistry()
		if err := r.Register(Template{
			Name:        "temp",
			Description: "temporary",
			Directory:   "temp",
			Version:     "1.0.0",
		}); err != nil {
			t.Fatalf("register failed: %v", err)
		}
		r2 := NewRegistry()
		if r2.Exists("temp") {
			t.Error("Exists should return false on different registry instance")
		}
	})
}

func TestConcurrency(t *testing.T) {
	t.Run("concurrent Register calls", func(t *testing.T) {
		r := NewRegistry()
		var wg sync.WaitGroup
		errs := make(chan error, 100)
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				name := string(rune('a' + n%26))
				err := r.Register(Template{
					Name:        name,
					Description: "desc",
					Directory:   name,
					Version:     "1.0.0",
				})
				if err != nil {
					errs <- err
				}
			}(i)
		}
		wg.Wait()
		close(errs)
		// Some registrations should fail with ErrTemplateExists, which is expected.
		for err := range errs {
			if !errors.IsCode(err, errors.ErrTemplateExists) {
				t.Errorf("unexpected error: %v", err)
			}
		}
		// Registry should be in a valid state.
		list := r.List()
		if len(list) == 0 {
			t.Error("expected at least one template after concurrent registration")
		}
	})

	t.Run("concurrent Get calls", func(t *testing.T) {
		r := NewRegistry()
		for i := 0; i < 10; i++ {
			name := string(rune('a' + i))
			if err := r.Register(Template{
				Name:        name,
				Description: "desc",
				Directory:   name,
				Version:     "1.0.0",
			}); err != nil {
				t.Fatalf("register failed: %v", err)
			}
		}
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				name := string(rune('a' + n%10))
				_, _ = r.Get(name)
			}(i)
		}
		wg.Wait()
	})

	t.Run("concurrent List calls", func(t *testing.T) {
		r := NewRegistry()
		for i := 0; i < 5; i++ {
			name := string(rune('a' + i))
			if err := r.Register(Template{
				Name:        name,
				Description: "desc",
				Directory:   name,
				Version:     "1.0.0",
			}); err != nil {
				t.Fatalf("register failed: %v", err)
			}
		}
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				list := r.List()
				if len(list) != 5 {
					t.Errorf("expected 5 templates, got %d", len(list))
				}
			}()
		}
		wg.Wait()
	})

	t.Run("concurrent Exists calls", func(t *testing.T) {
		r := NewRegistry()
		if err := r.Register(Template{
			Name:        "exists",
			Description: "desc",
			Directory:   "exists",
			Version:     "1.0.0",
		}); err != nil {
			t.Fatalf("register failed: %v", err)
		}
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = r.Exists("exists")
			}()
		}
		wg.Wait()
	})

	t.Run("concurrent mixed operations", func(t *testing.T) {
		r := NewRegistry()
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				name := string(rune('a' + n%10))
				tmpl := Template{
					Name:        name,
					Description: "desc",
					Directory:   name,
					Version:     "1.0.0",
				}
				_ = r.Register(tmpl)
				_, _ = r.Get(name)
				_ = r.List()
				_ = r.Exists(name)
			}(i)
		}
		wg.Wait()
	})
}

func TestTemplateFields(t *testing.T) {
	t.Run("all fields are stored correctly", func(t *testing.T) {
		r := NewRegistry()
		tmpl := Template{
			Name:               "full-template",
			Description:        "A full template description",
			Directory:          "full-template",
			Category:           CategoryAI,
			SupportedLanguages: []string{"Python", "Go", "Rust"},
			Version:            "2.5.3",
		}
		if err := r.Register(tmpl); err != nil {
			t.Fatalf("register failed: %v", err)
		}
		got, err := r.Get("full-template")
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if got.Name != "full-template" {
			t.Errorf("Name = %q, want %q", got.Name, "full-template")
		}
		if got.Description != "A full template description" {
			t.Errorf("Description = %q, want %q", got.Description, "A full template description")
		}
		if got.Directory != "full-template" {
			t.Errorf("Directory = %q, want %q", got.Directory, "full-template")
		}
		if got.Category != CategoryAI {
			t.Errorf("Category = %q, want %q", got.Category, CategoryAI)
		}
		if got.Version != "2.5.3" {
			t.Errorf("Version = %q, want %q", got.Version, "2.5.3")
		}
		if len(got.SupportedLanguages) != 3 {
			t.Fatalf("SupportedLanguages length = %d, want 3", len(got.SupportedLanguages))
		}
		for i, want := range []string{"Python", "Go", "Rust"} {
			if got.SupportedLanguages[i] != want {
				t.Errorf("SupportedLanguages[%d] = %q, want %q", i, got.SupportedLanguages[i], want)
			}
		}
	})

	t.Run("category constants have expected values", func(t *testing.T) {
		tests := []struct {
			cat  Category
			want string
		}{
			{CategoryCLI, "CLI"},
			{CategoryWeb, "Web"},
			{CategoryMobile, "Mobile"},
			{CategoryDesktop, "Desktop"},
			{CategoryAI, "AI"},
		}
		for _, tt := range tests {
			if string(tt.cat) != tt.want {
				t.Errorf("Category %v = %q, want %q", tt.cat, string(tt.cat), tt.want)
			}
		}
	})
}
