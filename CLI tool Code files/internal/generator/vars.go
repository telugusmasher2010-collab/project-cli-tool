// Package generator provides template variable substitution for project scaffolding.
package generator

import (
	"regexp"
	"sort"
	"strings"
	"sync"

	apperrors "github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

// placeholderRe matches {{key}} placeholders in template text.
var placeholderRe = regexp.MustCompile(`\{\{(\w+)\}\}`)

// Variables stores key-value pairs used for template substitution.
// It is safe for concurrent use; all methods are guarded by a sync.RWMutex.
type Variables struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewVariables creates an empty Variables instance.
func NewVariables() *Variables {
	return &Variables{data: make(map[string]string)}
}

// Set stores a variable, overwriting any existing value for the key.
func (v *Variables) Set(key, value string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.data[key] = value
}

// Get retrieves a variable. The second return value reports whether the key exists.
func (v *Variables) Get(key string) (string, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	val, ok := v.data[key]
	return val, ok
}

// Has reports whether a variable with the given key is set.
func (v *Variables) Has(key string) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	_, ok := v.data[key]
	return ok
}

// Clone returns a deep copy of the Variables. The clone is independent
// and safe to modify without affecting the original.
func (v *Variables) Clone() *Variables {
	v.mu.RLock()
	defer v.mu.RUnlock()

	clone := NewVariables()
	for k, val := range v.data {
		clone.data[k] = val
	}
	return clone
}

// Replace substitutes all {{key}} placeholders in text with stored values.
// Placeholders whose keys have no stored value are left unchanged.
// Placeholders with invalid syntax (e.g. unclosed braces) are left unchanged.
func (v *Variables) Replace(text string) string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.replaceLocked(text)
}

// ReplaceStrict is like Replace but returns an error if any placeholder
// references a key that has not been set. Successfully replaced placeholders
// are still substituted in the returned text.
func (v *Variables) ReplaceStrict(text string) (string, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	missing := v.missingKeysLocked(text)
	if len(missing) > 0 {
		return v.replaceLocked(text), missingVarsError(missing)
	}
	return v.replaceLocked(text), nil
}

// Validate checks that every {{key}} placeholder in text has a corresponding
// stored variable. It returns a structured error listing all missing keys,
// or nil if all placeholders are satisfied.
func (v *Variables) Validate(text string) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	missing := v.missingKeysLocked(text)
	if len(missing) > 0 {
		return missingVarsError(missing)
	}
	return nil
}

// Keys returns a sorted slice of all stored variable keys.
func (v *Variables) Keys() []string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	keys := make([]string, 0, len(v.data))
	for k := range v.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// replaceLocked performs substitution without acquiring locks.
// Caller must hold at least v.mu.RLock.
func (v *Variables) replaceLocked(text string) string {
	return placeholderRe.ReplaceAllStringFunc(text, func(match string) string {
		key := match[2 : len(match)-2]
		if val, ok := v.data[key]; ok {
			return val
		}
		return match
	})
}

// missingVarsError builds the structured error reported when one or more
// placeholders in text reference variables that have not been set.
func missingVarsError(missing []string) error {
	sort.Strings(missing)
	return apperrors.New(
		apperrors.ErrInvalidInput,
		"missing required variables: "+strings.Join(missing, ", "),
	)
}

// missingKeysLocked returns sorted keys of placeholders in text that
// have no stored value. Caller must hold at least v.mu.RLock.
func (v *Variables) missingKeysLocked(text string) []string {
	matches := placeholderRe.FindAllStringSubmatch(text, -1)
	var missing []string
	seen := make(map[string]bool)
	for _, m := range matches {
		key := m[1]
		if _, ok := v.data[key]; !ok && !seen[key] {
			missing = append(missing, key)
			seen[key] = true
		}
	}
	return missing
}
