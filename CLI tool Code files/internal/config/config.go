package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	DefaultTemplate string `mapstructure:"default_template"`
	OutputDir       string `mapstructure:"output_dir"`
	AuthorName      string `mapstructure:"author_name"`
}

// explicitPath, when set via SetConfigPath, overrides the default
// config file location (~/.proj-init/config.yml).
var explicitPath string

// SetConfigPath overrides the default config file location (~/.proj-init/config.yml).
// It must be called before Load.
func SetConfigPath(path string) {
	explicitPath = path
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(home, ".proj-init"))
	viper.SetEnvPrefix("PROJ_INIT")
	viper.AutomaticEnv()

	if explicitPath != "" {
		viper.SetConfigFile(explicitPath)
	} else {
		viper.SetConfigName("config")
	}

	viper.SetDefault("default_template", "")
	viper.SetDefault("output_dir", ".")
	viper.SetDefault("author_name", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func EnsureConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".proj-init")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}
