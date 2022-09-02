package config

type Schema string

const (
	HTTP  Schema = "http"
	HTTPS Schema = "https"
	GRPC  Schema = "grpc"
	GRPCS Schema = "grpcs"
)

type ServerConfig struct {
	Schema  Schema `json:"schema" example:"http or https"`
	Address string `json:"address" example:"127.0.0.1 or mydomain.com"`
	Port    int    `json:"port" example:"8080"`
}
