package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	ParseAppConfigFailed = errors.New("parse app config failed")
)

func ParseAppConfigFromEnv() (*AppConfig, error) {
	path := "config/" + os.Getenv("CONFIG") + ".yaml"
	return ParseAppConfig(path)
}

func ParseAppConfig(path string) (*AppConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ParseAppConfigFailed, err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	var config AppConfig
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ParseAppConfigFailed, err)
	}
	return &config, nil
}
