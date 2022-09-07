package config

type Schema string

const (
	HTTP  Schema = "http"
	HTTPS Schema = "https"
)

type SecurityConfig struct {
	CertificateFile string `json:"certificate_file" example:"./certificates/server.crt"`
	KeyFile         string `json:"key_file" example:"./certificates/server.key"`
}

type ServerConfig struct {
	Schema   Schema         `json:"schema" example:"http or https"`
	Address  string         `json:"address" example:"127.0.0.1 or mydomain.com"`
	Port     int            `json:"port" example:"8080"`
	Security SecurityConfig `json:"security"`
}
