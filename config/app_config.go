package config

type AppConfig struct {
	ServerCfg ServerConfig  `json:"server"`
	Logging   LoggingConfig `json:"logging"`
}
