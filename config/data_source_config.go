package config

import (
	"errors"
)

type DataSourceType string
type DataSourceConnOption string
type DataSourceConnOptionValueType string

const (
	FILE DataSourceType = "file"
	// MONGODB TODO (UMV): Mongo won't be using for sometime, maybe it will be removed completely
	MONGODB DataSourceType = "mongodb"
	// REDIS : redis server should be running with dump every write on a disk (AOF)
	REDIS DataSourceType = "redis"
)

const (
	OperationTimeout       DataSourceConnOption = "timeoutMS"
	ConnectionTimeout      DataSourceConnOption = "connectTimeoutMS"
	ConnectionsPool        DataSourceConnOption = "maxPoolSize"
	ReplicaSet             DataSourceConnOption = "replicaSet"
	MaxIdleTime            DataSourceConnOption = "maxIdleTimeMS"
	SocketTimeout          DataSourceConnOption = "socketTimeoutMS"
	ServerSelectionTimeout DataSourceConnOption = "serverSelectionTimeoutMS"
	HeartbeatFrequency     DataSourceConnOption = "heartbeatFrequencyMS"
	Tls                    DataSourceConnOption = "tls"
	WriteConcern           DataSourceConnOption = "w"
	DirectConnection       DataSourceConnOption = "directConnection"
)

const (
	String   DataSourceConnOptionValueType = "string"
	Integer  DataSourceConnOptionValueType = "integer"
	Boolean  DataSourceConnOptionValueType = "boolean"
	StrOrInt DataSourceConnOptionValueType = "str or int"
)

var (
	SourceISEmpty error = errors.New("field source (path to file or conn str to db) is empty")

	/*MongoDbOptionsTypes = map[DataSourceConnOption]DataSourceConnOptionValueType{
		OperationTimeout: Integer,
	}*/
)

// DataSourceConfig represent source where we can get
/*
 * We attempt to provide config that easily could be used with any datasource:
 * - json file (simplest RO mode)
 * - mongodb (but here we have very simple question how to pass parameters)
 * Source contains:
 * 1) if Type is FILE - full path to Json File
 * 2) if Type is REDIS - connection string without options, which looks like mongodb://user:password@host:port/
 * Options are connection options, see - https://www.mongodb.com/docs/drivers/go/current/fundamentals/connection/#std-label-golang-connection-guide
 * Here we should have Validator too
 */
type DataSourceConfig struct {
	Type        DataSourceType                  `json:"type"`
	Source      string                          `json:"source"`
	Credentials *CredentialsConfig              `json:"credentials"`
	Options     map[DataSourceConnOption]string `json:"options"`
}

func (cfg *DataSourceConfig) Validate() error {
	if len(cfg.Source) == 0 {
		return SourceISEmpty

	}
	if cfg.Type == MONGODB {
		return errors.New("mongodb is not supported ")
	}
	if cfg.Type == REDIS {

	}
	/*if cfg.Type == MONGODB {
		// validate options values ...
		allParamValidation := map[string]string{}
		for k, v := range cfg.Options {
			keyType := MongoDbOptionsTypes[k]
			err := cfg.validateParam(&keyType, &v)
			if err != nil {
				explanation := stringFormatter.Format("Error at MongoDb parameter \"{0}\" validation, reason: {1}", k, err.Error())
				allParamValidation[string(k)] = explanation
			}
		}
		if len(allParamValidation) > 0 {
			// todo(UMV): combine && return

		}
	}*/
	return nil
}

/*func (cfg *DataSourceConfig) validateParam(keyType *DataSourceConnOptionValueType, value *string) error {
	switch *keyType {
	case Integer:
		_, e := strconv.Atoi(*value)
		return e
	case Boolean:
		_, e := strconv.ParseBool(*value)
		return e
	default:
		return nil

	}
}*/
