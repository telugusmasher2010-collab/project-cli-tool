package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("returns error with correct fields", func(t *testing.T) {
		e := New(ErrTemplateNotFound, "missing template")
		if e.Code != ErrTemplateNotFound {
			t.Errorf("Code = %q, want %q", e.Code, ErrTemplateNotFound)
		}
		if e.Message != "missing template" {
			t.Errorf("Message = %q, want %q", e.Message, "missing template")
		}
		if e.Err != nil {
			t.Errorf("Err = %v, want nil", e.Err)
		}
	})
}

func TestWrap(t *testing.T) {
	t.Run("wraps cause correctly", func(t *testing.T) {
		cause := fmt.Errorf("disk full")
		e := Wrap(ErrFilesystem, "write failed", cause)
		if e.Code != ErrFilesystem {
			t.Errorf("Code = %q, want %q", e.Code, ErrFilesystem)
		}
		if e.Message != "write failed" {
			t.Errorf("Message = %q, want %q", e.Message, "write failed")
		}
		if e.Err != cause {
			t.Errorf("Err = %v, want %v", e.Err, cause)
		}
	})

	t.Run("wraps nil cause", func(t *testing.T) {
		e := Wrap(ErrInternal, "oops", nil)
		if e.Err != nil {
			t.Errorf("Err = %v, want nil", e.Err)
		}
	})

	t.Run("wraps another Error", func(t *testing.T) {
		inner := New(ErrInvalidTemplate, "bad yaml")
		outer := Wrap(ErrGenerationFailed, "could not generate", inner)
		if outer.Err != inner {
			t.Errorf("Err = %v, want %v", outer.Err, inner)
		}
	})
}

func TestError(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "without cause",
			err:  &Error{Code: ErrInternal, Message: "oops"},
			want: "[INTERNAL_ERROR] oops",
		},
		{
			name: "with cause",
			err: &Error{
				Code:    ErrGenerationFailed,
				Message: "build failed",
				Err:     fmt.Errorf("exit code 1"),
			},
			want: "[GENERATION_FAILED] build failed: exit code 1",
		},
		{
			name: "with nested Error cause",
			err: &Error{
				Code:    ErrFilesystem,
				Message: "copy failed",
				Err:     &Error{Code: ErrInternal, Message: "permission denied"},
			},
			want: "[FILESYSTEM_ERROR] copy failed: [INTERNAL_ERROR] permission denied",
		},
		{
			name: "all error codes",
			err:  &Error{Code: ErrTemplateNotFound, Message: "x"},
			want: "[TEMPLATE_NOT_FOUND] x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	t.Run("returns nil when no cause", func(t *testing.T) {
		e := New(ErrInternal, "no cause")
		if e.Unwrap() != nil {
			t.Errorf("Unwrap() = %v, want nil", e.Unwrap())
		}
	})

	t.Run("returns cause", func(t *testing.T) {
		cause := fmt.Errorf("underlying")
		e := Wrap(ErrInternal, "wrapped", cause)
		if e.Unwrap() != cause {
			t.Errorf("Unwrap() = %v, want %v", e.Unwrap(), cause)
		}
	})
}

func TestIsCode(t *testing.T) {
	t.Run("matches own code", func(t *testing.T) {
		e := New(ErrTemplateNotFound, "gone")
		if !IsCode(e, ErrTemplateNotFound) {
			t.Error("IsCode returned false for matching code")
		}
	})

	t.Run("rejects different code", func(t *testing.T) {
		e := New(ErrTemplateNotFound, "gone")
		if IsCode(e, ErrGenerationFailed) {
			t.Error("IsCode returned true for non-matching code")
		}
	})

	t.Run("traverses wrapped errors", func(t *testing.T) {
		cause := New(ErrInvalidTemplate, "bad")
		outer := Wrap(ErrGenerationFailed, "gen failed", cause)
		if !IsCode(outer, ErrInvalidTemplate) {
			t.Error("IsCode should find inner error code through wrap chain")
		}
	})

	t.Run("traverses standard library wraps", func(t *testing.T) {
		inner := New(ErrFilesystem, "disk")
		outer := fmt.Errorf("outer: %w", inner)
		if !IsCode(outer, ErrFilesystem) {
			t.Error("IsCode should find *Error through standard fmt.Errorf wrapping")
		}
	})

	t.Run("returns false for non-Error", func(t *testing.T) {
		if IsCode(fmt.Errorf("plain"), ErrInternal) {
			t.Error("IsCode should return false for non-Error errors")
		}
	})

	t.Run("returns false for nil", func(t *testing.T) {
		if IsCode(nil, ErrInternal) {
			t.Error("IsCode should return false for nil")
		}
	})
}

func TestIsTemplateNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"direct match", New(ErrTemplateNotFound, "nope"), true},
		{"wrapped", Wrap(ErrGenerationFailed, "fail", New(ErrTemplateNotFound, "nope")), true},
		{"wrong code", New(ErrInternal, "x"), false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsTemplateNotFound(tt.err)
			if got != tt.want {
				t.Errorf("IsTemplateNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsGenerationFailed(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"direct match", New(ErrGenerationFailed, "boom"), true},
		{"wrapped", Wrap(ErrInternal, "inner", New(ErrGenerationFailed, "boom")), true},
		{"wrong code", New(ErrTemplateNotFound, "x"), false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsGenerationFailed(tt.err)
			if got != tt.want {
				t.Errorf("IsGenerationFailed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInvalidInput(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"direct match", New(ErrInvalidInput, "bad"), true},
		{"wrapped", Wrap(ErrInternal, "inner", New(ErrInvalidInput, "bad")), true},
		{"wrong code", New(ErrConfigInvalid, "x"), false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsInvalidInput(tt.err)
			if got != tt.want {
				t.Errorf("IsInvalidInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsOutputExists(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"direct match", New(ErrOutputExists, "already here"), true},
		{"wrapped", Wrap(ErrGenerationFailed, "fail", New(ErrOutputExists, "already here")), true},
		{"wrong code", New(ErrInternal, "x"), false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsOutputExists(tt.err)
			if got != tt.want {
				t.Errorf("IsOutputExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserMessage(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"template not found", New(ErrTemplateNotFound, "tauri-llm"), "Template not found: tauri-llm"},
		{"template exists", New(ErrTemplateExists, "tauri-llm"), "Template already exists: tauri-llm"},
		{"invalid template", New(ErrInvalidTemplate, "bad yaml"), "Invalid template: bad yaml"},
		{"generation failed", New(ErrGenerationFailed, "compile error"), "Failed to generate project: compile error"},
		{"invalid input", New(ErrInvalidInput, "empty name"), "Invalid input: empty name"},
		{"config not found", New(ErrConfigNotFound, "missing"), "Config not found: missing"},
		{"config invalid", New(ErrConfigInvalid, "parse error"), "Invalid config: parse error"},
		{"filesystem", New(ErrFilesystem, "no space"), "Filesystem error: no space"},
		{"internal error", New(ErrInternal, "nil pointer"), "Something went wrong: nil pointer"},
		{"unknown code", New(Code("UNKNOWN"), "mystery"), "Something went wrong: mystery"},
		{"plain error", fmt.Errorf("plain error"), "plain error"},
		{"nil", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserMessage(tt.err)
			if got != tt.want {
				t.Errorf("UserMessage() = %q, want %q", got, tt.want)
			}
		})
	}

	t.Run("traverses wrap chain to find Error", func(t *testing.T) {
		inner := New(ErrFilesystem, "disk full")
		outer := fmt.Errorf("outer: %w", inner)
		got := UserMessage(outer)
		want := "Filesystem error: disk full"
		if got != want {
			t.Errorf("UserMessage() = %q, want %q", got, want)
		}
	})
}

func TestErrorInterface(t *testing.T) {
	var err error = New(ErrInternal, "test")
	if err == nil {
		t.Fatal("error should not be nil")
	}
}

func TestErrorsAs(t *testing.T) {
	inner := New(ErrGenerationFailed, "inner fail")
	outer := Wrap(ErrFilesystem, "outer", inner)

	var target *Error
	if !errors.As(outer, &target) {
		t.Fatal("errors.As should find *Error in chain")
	}
	if target.Code != ErrFilesystem {
		t.Errorf("Code = %q, want %q", target.Code, ErrFilesystem)
	}
}

func TestErrorsIs(t *testing.T) {
	inner := New(ErrInternal, "oops")
	outer := Wrap(ErrGenerationFailed, "failed", inner)

	if !errors.Is(outer, inner) {
		t.Error("errors.Is should find inner *Error in chain")
	}

	other := New(ErrInternal, "different")
	if errors.Is(outer, other) {
		t.Error("errors.Is should not match unrelated errors")
	}
}

func TestAllCodes(t *testing.T) {
	codes := []Code{
		ErrTemplateNotFound,
		ErrTemplateExists,
		ErrInvalidTemplate,
		ErrInvalidInput,
		ErrConfigNotFound,
		ErrConfigInvalid,
		ErrGenerationFailed,
		ErrFilesystem,
		ErrInternal,
		ErrOutputExists,
	}

	for _, code := range codes {
		t.Run(string(code), func(t *testing.T) {
			e := New(code, "test")
			if !IsCode(e, code) {
				t.Errorf("IsCode should match %s", code)
			}
			msg := UserMessage(e)
			if msg == "" {
				t.Error("UserMessage should not return empty string")
			}
		})
	}
}
