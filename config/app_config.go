package config

type AppConfig struct {
	ServerCfg  ServerConfig     `json:"server"`
	DataSource DataSourceConfig `json:"data_source"`
	Logging    LoggingConfig    `json:"logging"`
}
