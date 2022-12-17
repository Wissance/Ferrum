package config

type DataSourceType string

const (
	FILE    DataSourceType = "file"
	MONGODB DataSourceType = "mongodb"
)

// DataSourceConfig represent source where we can get
/*
 * We attempt to provide config that easily could be used with any datasource:
 * - json file (simplest RO mode)
 * - mongodb (but here we have very simple question how to pass parameters)
 */
type DataSourceConfig struct {
	Type   DataSourceType `json:"type"`
	Source string         `json:"source"`
}
