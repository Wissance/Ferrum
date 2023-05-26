package config

import "os"

const serverValidationErrExitCode = 567
const dataSourceValidationErrExitCode = 568
const loggingSystemValidationErrExitCode = 569

type AppConfig struct {
	ServerCfg  ServerConfig     `json:"server"`
	DataSource DataSourceConfig `json:"data_source"`
	Logging    LoggingConfig    `json:"logging"`
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
