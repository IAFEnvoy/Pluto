package util

import (
	"encoding/json"
	"os"
)

func LoadConfig[T any](path string) (T, error) {
	var config T
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return config, err
	}
	return config, nil
}

func SaveConfig(config any, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
