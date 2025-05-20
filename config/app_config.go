package config

import (
	gce "github.com/wissance/go-config-extender"
	"os"
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
	cfg, err := gce.LoadJSONConfigWithEnvOverride[AppConfig](pathToConfig)
	if err != nil {
		return nil, err
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
