package generator

import (
	"errors"
	"strings"
	"sync"
	"testing"

	apperrors "github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

func TestNewVariables(t *testing.T) {
	v := NewVariables()
	if v == nil {
		t.Fatal("NewVariables returned nil")
	}
	if len(v.Keys()) != 0 {
		t.Errorf("new Variables should be empty, got %d keys", len(v.Keys()))
	}
}

func TestSetAndGet(t *testing.T) {
	t.Run("set and get single key", func(t *testing.T) {
		v := NewVariables()
		v.Set("name", "proj-init")

		got, ok := v.Get("name")
		if !ok {
			t.Fatal("Get returned false for existing key")
		}
		if got != "proj-init" {
			t.Errorf("Get = %q, want %q", got, "proj-init")
		}
	})

	t.Run("get nonexistent key", func(t *testing.T) {
		v := NewVariables()
		got, ok := v.Get("nope")
		if ok {
			t.Error("Get returned true for nonexistent key")
		}
		if got != "" {
			t.Errorf("Get = %q, want empty string", got)
		}
	})

	t.Run("overwrite existing key", func(t *testing.T) {
		v := NewVariables()
		v.Set("x", "old")
		v.Set("x", "new")

		got, _ := v.Get("x")
		if got != "new" {
			t.Errorf("Get after overwrite = %q, want %q", got, "new")
		}
	})

	t.Run("empty value is valid", func(t *testing.T) {
		v := NewVariables()
		v.Set("key", "")

		_, ok := v.Get("key")
		if !ok {
			t.Error("empty value should still be present")
		}
	})
}

func TestHas(t *testing.T) {
	v := NewVariables()

	if v.Has("missing") {
		t.Error("Has returned true for nonexistent key")
	}

	v.Set("exists", "val")
	if !v.Has("exists") {
		t.Error("Has returned false for existing key")
	}
}

func TestClone(t *testing.T) {
	t.Run("clone contains same data", func(t *testing.T) {
		v := NewVariables()
		v.Set("a", "1")
		v.Set("b", "2")

		c := v.Clone()
		for _, key := range []string{"a", "b"} {
			got, _ := c.Get(key)
			orig, _ := v.Get(key)
			if got != orig {
				t.Errorf("Clone key %q = %q, want %q", key, got, orig)
			}
		}
	})

	t.Run("clone is independent", func(t *testing.T) {
		v := NewVariables()
		v.Set("shared", "original")

		c := v.Clone()
		c.Set("shared", "modified")
		c.Set("clone-only", "value")

		orig, _ := v.Get("shared")
		if orig != "original" {
			t.Errorf("modifying clone changed original: got %q", orig)
		}

		if v.Has("clone-only") {
			t.Error("clone-only key should not exist in original")
		}
	})

	t.Run("modifying original does not affect clone", func(t *testing.T) {
		v := NewVariables()
		v.Set("key", "original")

		c := v.Clone()
		v.Set("key", "changed")
		v.Set("new-key", "value")

		got, _ := c.Get("key")
		if got != "original" {
			t.Errorf("modifying original affected clone: got %q", got)
		}

		if c.Has("new-key") {
			t.Error("new key in original should not appear in clone")
		}
	})
}

func TestKeys(t *testing.T) {
	v := NewVariables()
	v.Set("charlie", "3")
	v.Set("alpha", "1")
	v.Set("bravo", "2")

	keys := v.Keys()
	expected := []string{"alpha", "bravo", "charlie"}
	if len(keys) != len(expected) {
		t.Fatalf("Keys() returned %d items, want %d", len(keys), len(expected))
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("Keys()[%d] = %q, want %q", i, k, expected[i])
		}
	}
}

func TestReplace(t *testing.T) {
	tests := []struct {
		name    string
		vars    map[string]string
		input   string
		want    string
	}{
		{
			name:  "single placeholder",
			vars:  map[string]string{"name": "proj-init"},
			input: "project: {{name}}",
			want:  "project: proj-init",
		},
		{
			name:  "multiple placeholders",
			vars:  map[string]string{"name": "proj-init", "author": "Suhrit"},
			input: "{{name}} by {{author}}",
			want:  "proj-init by Suhrit",
		},
		{
			name:  "repeated placeholder",
			vars:  map[string]string{"name": "proj-init"},
			input: "{{name}} and {{name}}",
			want:  "proj-init and proj-init",
		},
		{
			name:  "unknown placeholder preserved",
			vars:  map[string]string{"name": "proj-init"},
			input: "{{name}} and {{unknown}}",
			want:  "proj-init and {{unknown}}",
		},
		{
			name:  "no placeholders",
			vars:  map[string]string{"name": "proj-init"},
			input: "no placeholders here",
			want:  "no placeholders here",
		},
		{
			name:  "empty input",
			vars:  map[string]string{"name": "proj-init"},
			input: "",
			want:  "",
		},
		{
			name:  "no variables set",
			vars:  map[string]string{},
			input: "{{a}} and {{b}}",
			want:  "{{a}} and {{b}}",
		},
		{
			name:  "empty value replacement",
			vars:  map[string]string{"name": ""},
			input: "hello {{name}}!",
			want:  "hello !",
		},
		{
			name:  "adjacent placeholders",
			vars:  map[string]string{"a": "X", "b": "Y"},
			input: "{{a}}{{b}}",
			want:  "XY",
		},
		{
			name:  "placeholder in multiline text",
			vars:  map[string]string{"project": "my-app"},
			input: "name: {{project}}\nversion: 1.0",
			want:  "name: my-app\nversion: 1.0",
		},
		{
			name:  "partial match not replaced",
			vars:  map[string]string{"name": "proj-init"},
			input: "{{name-variant}} and {{name}}",
			want:  "{{name-variant}} and proj-init",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVariables()
			for k, val := range tt.vars {
				v.Set(k, val)
			}
			got := v.Replace(tt.input)
			if got != tt.want {
				t.Errorf("Replace(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestReplaceStrict(t *testing.T) {
	t.Run("all variables present", func(t *testing.T) {
		v := NewVariables()
		v.Set("name", "proj-init")
		v.Set("author", "Suhrit")

		got, err := v.ReplaceStrict("{{name}} by {{author}}")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "proj-init by Suhrit" {
			t.Errorf("ReplaceStrict = %q, want %q", got, "proj-init by Suhrit")
		}
	})

	t.Run("missing variable returns error", func(t *testing.T) {
		v := NewVariables()
		v.Set("name", "proj-init")

		got, err := v.ReplaceStrict("{{name}} and {{missing}}")
		if err == nil {
			t.Fatal("expected error for missing variable")
		}

		if !apperrors.IsInvalidInput(err) {
			t.Errorf("error code should be ErrInvalidInput, got: %v", err)
		}

		if !strings.Contains(err.Error(), "missing") {
			t.Errorf("error should mention missing, got: %v", err)
		}

		if got != "proj-init and {{missing}}" {
			t.Errorf("ReplaceStrict should still substitute known vars, got %q", got)
		}
	})

	t.Run("multiple missing variables", func(t *testing.T) {
		v := NewVariables()

		_, err := v.ReplaceStrict("{{a}} and {{b}} and {{c}}")
		if err == nil {
			t.Fatal("expected error")
		}

		errMsg := err.Error()
		for _, key := range []string{"a", "b", "c"} {
			if !strings.Contains(errMsg, key) {
				t.Errorf("error should mention key %q, got: %s", key, errMsg)
			}
		}
	})

	t.Run("no placeholders", func(t *testing.T) {
		v := NewVariables()
		got, err := v.ReplaceStrict("no placeholders")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "no placeholders" {
			t.Errorf("got %q, want %q", got, "no placeholders")
		}
	})
}

func TestValidate(t *testing.T) {
	t.Run("all variables present", func(t *testing.T) {
		v := NewVariables()
		v.Set("name", "proj-init")
		v.Set("author", "Suhrit")

		if err := v.Validate("{{name}} by {{author}}"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("missing variable", func(t *testing.T) {
		v := NewVariables()
		v.Set("name", "proj-init")

		err := v.Validate("{{name}} and {{missing}}")
		if err == nil {
			t.Fatal("expected error for missing variable")
		}
		if !apperrors.IsInvalidInput(err) {
			t.Errorf("error code should be ErrInvalidInput, got: %v", err)
		}
	})

	t.Run("no placeholders", func(t *testing.T) {
		v := NewVariables()
		if err := v.Validate("plain text"); err != nil {
			t.Errorf("expected nil, got: %v", err)
		}
	})

	t.Run("empty text", func(t *testing.T) {
		v := NewVariables()
		if err := v.Validate(""); err != nil {
			t.Errorf("expected nil, got: %v", err)
		}
	})

	t.Run("duplicate missing key reported once", func(t *testing.T) {
		v := NewVariables()
		err := v.Validate("{{x}} and {{x}}")
		if err == nil {
			t.Fatal("expected error")
		}
		errMsg := err.Error()
		count := strings.Count(errMsg, "x")
		if count > 1 {
			t.Errorf("missing key 'x' reported multiple times: %s", errMsg)
		}
	})
}

func TestReplaceAndValidate(t *testing.T) {
	v := NewVariables()
	v.Set("name", "proj")

	// Validate first, then replace — common workflow pattern
	if err := v.Validate("{{name}} and {{other}}"); err != nil {
		// validation caught the issue
		if !apperrors.IsInvalidInput(err) {
			t.Errorf("wrong error code: %v", err)
		}
	}

	// Replace still works for known vars
	got := v.Replace("{{name}} and {{other}}")
	if got != "proj and {{other}}" {
		t.Errorf("Replace = %q", got)
	}
}

func TestConcurrencySafety(t *testing.T) {
	v := NewVariables()
	v.Set("shared", "value")

	var wg sync.WaitGroup
	errs := make(chan error, 200)

	// Concurrent readers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = v.Get("shared")
			_ = v.Has("shared")
			_ = v.Replace("{{shared}}")
			_, _ = v.ReplaceStrict("{{shared}}")
			_ = v.Validate("{{shared}}")
			_ = v.Keys()
			_ = v.Clone()
		}()
	}

	// Concurrent cloners
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := v.Clone()
			got, _ := c.Get("shared")
			if got != "value" {
				errs <- errors.New("clone returned wrong value")
			}
		}()
	}

	wg.Wait()
	close(errs)

	for e := range errs {
		t.Errorf("concurrent access error: %v", e)
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("placeholder-like text without full syntax", func(t *testing.T) {
		v := NewVariables()
		v.Set("name", "proj")
		tests := []struct{ input, want string }{
			{"{name}", "{name}"},
			{"{{name", "{{name"},
			{"name}}", "name}}"},
			{"{{{name}}}", "{proj}"},
			{"{{ name }}", "{{ name }}"},
		}
		for _, tt := range tests {
			got := v.Replace(tt.input)
			if got != tt.want {
				t.Errorf("Replace(%q) = %q, want %q", tt.input, got, tt.want)
			}
		}
	})

	t.Run("empty key in placeholder", func(t *testing.T) {
		v := NewVariables()
		got := v.Replace("{{}}")
		if got != "{{}}" {
			t.Errorf("Replace({{}}) = %q, want {{}}", got)
		}
	})

	t.Run("numeric key", func(t *testing.T) {
		v := NewVariables()
		v.Set("123", "numeric")
		got := v.Replace("{{123}}")
		if got != "numeric" {
			t.Errorf("Replace({{123}}) = %q", got)
		}
	})

	t.Run("underscore key", func(t *testing.T) {
		v := NewVariables()
		v.Set("my_var", "val")
		got := v.Replace("{{my_var}}")
		if got != "val" {
			t.Errorf("Replace({{my_var}}) = %q", got)
		}
	})
}
