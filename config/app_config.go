package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wissance/stringFormatter"
)

const (
	serverValidationErrExitCode        = 567
	dataSourceValidationErrExitCode    = 568
	loggingSystemValidationErrExitCode = 569
)

type AppConfig struct {
	ServerCfg  ServerConfig     `json:"server"`
	DataSource DataSourceConfig `json:"data_source"`
	Logging    LoggingConfig    `json:"logging"`
}

func ReadAppConfig(pathToConfig string) (*AppConfig, error) {
	absPath, err := filepath.Abs(pathToConfig)
	if err != nil {
		return nil, fmt.Errorf(stringFormatter.Format("An error occurred during getting config file abs path: {0}", err.Error()))
	}
	fileData, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf(stringFormatter.Format("An error occurred during config file reading: {0}", err.Error()))
	}
	var cfg AppConfig
	if err = json.Unmarshal(fileData, &cfg); err != nil {
		return nil, fmt.Errorf(stringFormatter.Format("An error occurred during config file unmarshal: {0}", err.Error()))
	}
	cfg.Validate()
	return &cfg, nil
}

func (cfg *AppConfig) Validate() {
	serverCfgValidationErr := cfg.ServerCfg.Validate()
	if serverCfgValidationErr != nil {
		println(serverCfgValidationErr.Error())
		os.Exit(serverValidationErrExitCode)
	}

	dataSourceCfgValidationErr := cfg.DataSource.Validate()
	if dataSourceCfgValidationErr != nil {
		println(dataSourceCfgValidationErr.Error())
		os.Exit(dataSourceValidationErrExitCode)
	}

	loggingSystemCfgValidationErr := cfg.Logging.Validate()
	if loggingSystemCfgValidationErr != nil {
		println(loggingSystemCfgValidationErr.Error())
		os.Exit(loggingSystemValidationErrExitCode)
	}
}
