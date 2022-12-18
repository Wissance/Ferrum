package config

import "errors"

type DataSourceType string
type MongoDbOption string
type MongoDbOptionValueType string

const (
	FILE    DataSourceType = "file"
	MONGODB DataSourceType = "mongodb"
)

const (
	OperationTimeout       MongoDbOption = "timeoutMS"
	ConnectionTimeout      MongoDbOption = "connectTimeoutMS"
	ConnectionsPool        MongoDbOption = "maxPoolSize"
	ReplicaSet             MongoDbOption = "replicaSet"
	MaxIdleTime            MongoDbOption = "maxIdleTimeMS"
	SocketTimeout          MongoDbOption = "socketTimeoutMS"
	ServerSelectionTimeout MongoDbOption = "serverSelectionTimeoutMS"
	HeartbeatFrequency     MongoDbOption = "heartbeatFrequencyMS"
	Tls                    MongoDbOption = "tls"
	WriteConcern           MongoDbOption = "w"
	DirectConnection       MongoDbOption = "directConnection"
)

const (
	String MongoDbOptionValueType = "string"
	Integer MongoDbOptionValueType = "integer"
	Boolean MongoDbOptionValueType = "boolean"
	StrOrInt MongoDbOptionValueType = "str or int"
)

var (
	SourceISEmpty error = errors.New("field source (path to file or conn str to db) is empty")

	MongoDbOptionsTypes = map[MongoDbOption]string
)

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
	Type    DataSourceType           `json:"type"`
	Source  string                   `json:"source"`
	Options map[MongoDbOption]string `json:"options"`
}

func (cfg *DataSourceConfig) Validate() error {
	if len(cfg.Source) == 0 {
		return SourceISEmpty

	}
	if cfg.Type == MONGODB {
		// validate options values ...
	}
	return nil
}
