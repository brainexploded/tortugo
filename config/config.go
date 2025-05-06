package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const (
	CONFIG_FILENAME = ".tortugo"
)

type Config struct {
	LibraryPath   string
	IndexFilename string
}

func Load(path string) (*Config, error) {
	if path == "" {
		homepath, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("Can't detect user's home dir %w", err)
		}
		path = homepath
	}

	path = filepath.Join(path, CONFIG_FILENAME)

	fi, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return CreateNewConfig(path)
	} else if err != nil {
		return nil, fmt.Errorf("Error checking config existence %w", err)
	}

	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Can't create config file %w", err)
	}

	data := make([]byte, fi.Size())

	_, err = fd.Read(data)
	if err != nil {
		return nil, fmt.Errorf("Can't read config file %w", err)
	}

	var cfg Config

	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("Can't decode config file %w", err)
	}

	cfg.LibraryPath = os.ExpandEnv(cfg.LibraryPath)

	return &cfg, nil
}

func CreateNewConfig(path string) (*Config, error) {
	fd, createErr := os.Create(path)
	if createErr != nil {
		return nil, fmt.Errorf("Can't create config file %w", createErr)
	}
	defer fd.Close()

	defaultConfig := Config{}

	b, err := toml.Marshal(defaultConfig)
	if err != nil {
		return nil, fmt.Errorf("Can't write to config file %w", err)
	}

	fd.Write(b)

	return &defaultConfig, nil
}
