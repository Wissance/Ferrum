package config

type ServerConfig struct {
	Schema  string `json:"schema" example:"http or https"`
	Address string `json:"address" example:"127.0.0.1 or mydomain.com"`
	Port    int    `json:"port" example:"8080"`
}
