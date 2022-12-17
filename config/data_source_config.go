package config

type DataSourceType string
type MongoDbOptions string

const (
	FILE    DataSourceType = "file"
	MONGODB DataSourceType = "mongodb"
)

const ()

// DataSourceConfig represent source where we can get
/*
 * We attempt to provide config that easily could be used with any datasource:
 * - json file (simplest RO mode)
 * - mongodb (but here we have very simple question how to pass parameters)
 * Source contains:
 * 1) if Type is FILE - full path to Json File
 * 2) if Type is MONGODB - connection string without options, which looks like mongodb://user:password@host:port/
 * Options are connection options, see - https://www.mongodb.com/docs/drivers/go/current/fundamentals/connection/#std-label-golang-connection-guide
 * Here we should have Validator too
 */
type DataSourceConfig struct {
	Type    DataSourceType    `json:"type"`
	Source  string            `json:"source"`
	Options map[string]string `json:"options"`
}

func (cfg *DataSourceConfig) Validate() error {
	return nil
}
