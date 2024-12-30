package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type AppConfig struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func ReadConfig() (AppConfig, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return AppConfig{}, err
	}
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return AppConfig{}, err
	}

	defer jsonFile.Close()

	var config AppConfig
	byteValue, _ := io.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return AppConfig{}, err
	}

	return config, nil
}

func (config *AppConfig) SetUser(username string) error {
	config.CurrentUserName = username
	return WriteConfig(*config)
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".gatorconfig.json"), nil
}

func WriteConfig(config AppConfig) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}
	jsonString, err := json.MarshalIndent(config, "", "	")
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, jsonString, 0600)
	if err != nil {
		return err
	}
	return nil
}
