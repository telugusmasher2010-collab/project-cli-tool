package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadExplicitPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	content := "default_template: tauri-llm\nauthor_name: 'Abhi Test'\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	SetConfigPath(path)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.DefaultTemplate != "tauri-llm" {
		t.Errorf("DefaultTemplate = %q, want tauri-llm", cfg.DefaultTemplate)
	}
	if cfg.AuthorName != "Abhi Test" {
		t.Errorf("AuthorName = %q, want Abhi Test", cfg.AuthorName)
	}
}
