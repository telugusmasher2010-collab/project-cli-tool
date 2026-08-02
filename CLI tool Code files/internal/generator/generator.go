package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	apperrors "github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/templates"
)

// defaultFilePerm is used for regular generated files.
const defaultFilePerm os.FileMode = 0644

// executableFilePerm is used for files with executable extensions.
const executableFilePerm os.FileMode = 0755

// executableExtensions lists file extensions that should be marked executable.
var executableExtensions = []string{".sh", ".bat", ".cmd", ".ps1"}

// Options configures the behavior of the Generator.
type Options struct {
	// Overwrite allows generation into a directory that already exists.
	// When false (the default), Generate returns ErrOutputExists if the
	// target directory is non-empty.
	Overwrite bool
	// Hooks are post-generation actions run after all files are written.
	// When non-nil they replace the default template hooks.
	Hooks *HookRunner
	// AutoHooks runs the default hook set computed for the generated
	// template (see HooksForTemplate) after all files are written.
	// It is ignored when Hooks is non-nil.
	AutoHooks bool
}

// Generator scaffolds a project directory from an embedded template.
type Generator struct {
	outputDir string
	vars      *Variables
	opts      Options
}

// New creates a Generator that will write to outputDir using the given
// variables and options.
func New(outputDir string, vars *Variables, opts Options) *Generator {
	if vars == nil {
		vars = NewVariables()
	}
	return &Generator{
		outputDir: filepath.Clean(outputDir),
		vars:      vars,
		opts:      opts,
	}
}

// Generate creates the project from the named template.
// It validates the template exists, checks the output directory, then
// walks every file in the template, applies variable substitution, and
// writes it to the output tree. If hooks are configured in Options,
// they run after all files are written.
func (g *Generator) Generate(templateName string) error {
	if templateName == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "template name must not be empty")
	}

	if _, err := templates.Get(templateName); err != nil {
		return err
	}

	if err := g.checkOutputDir(); err != nil {
		return err
	}

	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return apperrors.Wrap(apperrors.ErrFilesystem,
			fmt.Sprintf("failed to create output directory %q", g.outputDir), err)
	}

	files, err := templates.WalkFiles(templateName)
	if err != nil {
		return err
	}

	for _, relPath := range files {
		if err := g.processFile(templateName, relPath); err != nil {
			return apperrors.Wrap(apperrors.ErrGenerationFailed,
				fmt.Sprintf("failed to process %q", relPath), err)
		}
	}

	if g.opts.Hooks != nil {
		if err := g.opts.Hooks.RunAll(context.Background(), g.outputDir); err != nil {
			return err
		}
	}

	return nil
}

// Variables returns the generator's variable store so callers can
// inspect which placeholders were provided.
func (g *Generator) Variables() *Variables {
	return g.vars
}

// checkOutputDir ensures we don't clobber an existing project.
func (g *Generator) checkOutputDir() error {
	info, err := os.Stat(g.outputDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return apperrors.Wrap(apperrors.ErrFilesystem,
			fmt.Sprintf("failed to stat output directory %q", g.outputDir), err)
	}
	if !info.IsDir() {
		return apperrors.New(apperrors.ErrFilesystem,
			fmt.Sprintf("output path %q exists but is not a directory", g.outputDir))
	}
	if g.opts.Overwrite {
		return nil
	}
	entries, err := os.ReadDir(g.outputDir)
	if err != nil {
		return apperrors.Wrap(apperrors.ErrFilesystem,
			fmt.Sprintf("failed to read output directory %q", g.outputDir), err)
	}
	if len(entries) > 0 {
		return apperrors.New(apperrors.ErrOutputExists,
			fmt.Sprintf("directory %q already exists and is not empty", g.outputDir))
	}
	return nil
}

// processFile reads a single file from the embedded template, applies
// variable substitution, and writes it to the output directory.
func (g *Generator) processFile(templateName, relPath string) error {
	data, err := templates.ReadFile(templateName, relPath)
	if err != nil {
		return err
	}

	substituted := g.vars.Replace(string(data))

	outPath := filepath.Join(g.outputDir, filepath.FromSlash(relPath))
	outDir := filepath.Dir(outPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return apperrors.Wrap(apperrors.ErrFilesystem,
			fmt.Sprintf("failed to create directory %q", outDir), err)
	}

	perm := defaultFilePerm
	if isExecutable(relPath) {
		perm = executableFilePerm
	}

	if err := os.WriteFile(outPath, []byte(substituted), perm); err != nil {
		return apperrors.Wrap(apperrors.ErrFilesystem,
			fmt.Sprintf("failed to write file %q", outPath), err)
	}

	return nil
}

// isExecutable reports whether the file should be marked executable
// based on its extension.
func isExecutable(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, e := range executableExtensions {
		if ext == e {
			return true
		}
	}
	return false
}
