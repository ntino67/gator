package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("error getting config file path: %w", err)
	}

	res, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("can't read the file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(res, &config); err != nil {
		return Config{}, fmt.Errorf("error while decoding the json: %w", err)
	}

	return config, nil
}

func (cfg *Config) Save() error {
	path, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting config file path: %w", err)
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshalling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error writing to config file: %w", err)
	}

	return nil
}

func (cfg *Config) SetUser(name string) error {
	cfg.CurrentUserName = name
	return cfg.Save()
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}

	path := filepath.Join(homeDir, configFileName)

	return path, nil
}
