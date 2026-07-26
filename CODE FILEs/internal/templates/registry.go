package templates

import (
	"fmt"
	"sort"
	"sync"

	"github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

// Category represents a project template category.
type Category string

const (
	CategoryCLI     Category = "CLI"
	CategoryWeb     Category = "Web"
	CategoryMobile  Category = "Mobile"
	CategoryDesktop Category = "Desktop"
	CategoryAI      Category = "AI"
)

// Template describes a scaffolding template available for project generation.
type Template struct {
	Name               string
	Description        string
	Directory          string
	Category           Category
	SupportedLanguages []string
	Version            string
}

// Registry stores and manages available project templates. It is safe for
// concurrent use.
type Registry struct {
	mu        sync.RWMutex
	templates map[string]Template
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		templates: make(map[string]Template),
	}
}

// Register adds a template to the registry. It returns ErrTemplateExists
// if a template with the same name is already registered. Fields are
// validated and ErrInvalidInput is returned for empty required fields.
func (r *Registry) Register(t Template) error {
	if t.Name == "" {
		return errors.New(errors.ErrInvalidInput, "template name must not be empty")
	}
	if t.Description == "" {
		return errors.New(errors.ErrInvalidInput, fmt.Sprintf("description must not be empty for template %q", t.Name))
	}
	if t.Directory == "" {
		return errors.New(errors.ErrInvalidInput, fmt.Sprintf("directory must not be empty for template %q", t.Name))
	}
	if t.Version == "" {
		return errors.New(errors.ErrInvalidInput, fmt.Sprintf("version must not be empty for template %q", t.Name))
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.templates[t.Name]; exists {
		return errors.New(errors.ErrTemplateExists, fmt.Sprintf("template %q is already registered", t.Name))
	}

	r.templates[t.Name] = t
	return nil
}

// Get returns the template with the given name. Returns
// ErrTemplateNotFound if no such template exists.
func (r *Registry) Get(name string) (Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.templates[name]
	if !ok {
		return Template{}, errors.New(errors.ErrTemplateNotFound, fmt.Sprintf("template %q not found", name))
	}
	return t, nil
}

// List returns all registered templates sorted by name.
func (r *Registry) List() []Template {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Template, 0, len(r.templates))
	for _, t := range r.templates {
		result = append(result, t)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

// Exists reports whether a template with the given name is registered.
func (r *Registry) Exists(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.templates[name]
	return ok
}
