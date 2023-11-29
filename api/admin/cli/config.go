package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wissance/Ferrum/config"
)

type CliConfig struct {
	DataSource config.DataSourceConfig `json:"data_source"`
	Logging    config.LoggingConfig    `json:"logging"`
}

func readConfig(pathToConfigFile string) (*CliConfig, error) {
	fileData, err := os.ReadFile(pathToConfigFile)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile failed: %w", err)
	}
	var cfg CliConfig
	if err = json.Unmarshal(fileData, &cfg); err != nil {
		return nil, fmt.Errorf("json.Unmarshal failed: %w", err)
	}
	if err := cfg.DataSource.Validate(); err != nil {
		return nil, fmt.Errorf("DataSourceConfig.Validate failed: %w", err)
	}
	if err := cfg.Logging.Validate(); err != nil {
		return nil, fmt.Errorf("LoggingConfig.Validate failed: %w", err)
	}
	return &cfg, nil
}
