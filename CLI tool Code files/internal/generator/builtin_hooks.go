package generator

import (
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/templates"
)

// GitInitHook initializes a git repository in the generated project.
func GitInitHook() Hook {
	return &CommandHook{
		Name_:   "git init",
		Command: "git",
		Args:    []string{"init"},
	}
}

// NpmInstallHook installs Node.js dependencies in the generated project.
func NpmInstallHook() Hook {
	return &CommandHook{
		Name_:   "npm install",
		Command: "npm",
		Args:    []string{"install"},
	}
}

// FlutterPubGetHook fetches Flutter/Dart package dependencies in the
// generated project.
func FlutterPubGetHook() Hook {
	return &CommandHook{
		Name_:   "flutter pub get",
		Command: "flutter",
		Args:    []string{"pub", "get"},
	}
}

// GoModTidyHook runs `go mod tidy` in the generated project to prune and
// normalize its module dependencies.
func GoModTidyHook() Hook {
	return &CommandHook{
		Name_:   "go mod tidy",
		Command: "go",
		Args:    []string{"mod", "tidy"},
	}
}

// HooksForTemplate returns the default post-generation hook set for a
// template. It always starts with git init, then selects additional hooks
// by inspecting the embedded template for well-known manifest files:
//
//   - package.json -> npm install
//   - pubspec.yaml -> flutter pub get
//   - go.mod       -> go mod tidy
//
// Templates that cannot be inspected (or that contain none of the known
// manifests) fall back to git init alone.
func HooksForTemplate(templateName string) *HookRunner {
	files, err := templates.WalkFiles(templateName)
	if err != nil {
		return NewHookRunner(GitInitHook())
	}

	seen := make(map[string]bool, len(files))
	for _, f := range files {
		seen[f] = true
	}

	hooks := []Hook{GitInitHook()}
	if seen["package.json"] {
		hooks = append(hooks, NpmInstallHook())
	}
	if seen["pubspec.yaml"] {
		hooks = append(hooks, FlutterPubGetHook())
	}
	if seen["go.mod"] {
		hooks = append(hooks, GoModTidyHook())
	}
	return NewHookRunner(hooks...)
}
