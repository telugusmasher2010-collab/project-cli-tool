package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

//go:embed all:tauri-llm all:whatsapp-bot all:expense-splitter all:next-webapp all:react-native-map
var TemplateFS embed.FS // the all: prefix embeds dotfiles (.gitignore) and underscore-prefixed files (app/_layout.tsx)

type TemplateInfo struct {
	Name        string
	Description string
	Directory   string
}

var registry = []TemplateInfo{
	{
		Name:        "tauri-llm",
		Description: "Tauri v2 + Rust + React + local LLM sidecar",
		Directory:   "tauri-llm",
	},
	{
		Name:        "whatsapp-bot",
		Description: "Node.js + Baileys + SQLite + Fastify",
		Directory:   "whatsapp-bot",
	},
	{
		Name:        "expense-splitter",
		Description: "Flutter + Dart + Supabase + UPI",
		Directory:   "expense-splitter",
	},
	{
		Name:        "next-webapp",
		Description: "Next.js 15 + React 19 + TypeScript (App Router)",
		Directory:   "next-webapp",
	},
	{
		Name:        "react-native-map",
		Description: "React Native + Expo + Expo Router + Maps",
		Directory:   "react-native-map",
	},
}

func List() []TemplateInfo {
	result := make([]TemplateInfo, len(registry))
	copy(result, registry)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func Get(name string) (TemplateInfo, error) {
	for _, t := range registry {
		if t.Name == name {
			return t, nil
		}
	}
	return TemplateInfo{}, errors.New(errors.ErrTemplateNotFound, fmt.Sprintf("template %q not found", name))
}

func ReadFile(templateName, filePath string) ([]byte, error) {
	fullPath := templateName + "/" + filePath
	data, err := TemplateFS.ReadFile(fullPath)
	if err != nil {
		return nil, errors.Wrap(errors.ErrTemplateNotFound, fmt.Sprintf("file %q not found in template %q", filePath, templateName), err)
	}
	return data, nil
}

func WalkFiles(templateName string) ([]string, error) {
	info, err := Get(templateName)
	if err != nil {
		return nil, err
	}

	var files []string
	err = fs.WalkDir(TemplateFS, info.Directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		relPath := strings.TrimPrefix(path, info.Directory+"/")
		files = append(files, relPath)
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(errors.ErrGenerationFailed, fmt.Sprintf("failed to walk template %q", templateName), err)
	}

	sort.Strings(files)
	return files, nil
}

func Exists(name string) bool {
	_, err := Get(name)
	return err == nil
}
