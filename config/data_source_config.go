package config

import (
	"errors"
	"github.com/wissance/Ferrum/utils/validators"
	sf "github.com/wissance/stringFormatter"
	"strconv"
	"strings"
)

type DataSourceType string
type DataSourceConnOption string

const (
	FILE DataSourceType = "file"
	// MONGODB TODO (UMV): Mongo won't be using for sometime, maybe it will be removed completely
	MONGODB DataSourceType = "mongodb"
	// REDIS : redis server should be running with dump every write on a disk (AOF)
	REDIS DataSourceType = "redis"
)

const (
	// DbNumber is a REDIS connection options, here we expect to receive int in a string
	DbNumber DataSourceConnOption = "db_number"
	// UseTls is a REDIS option to set up tls.Config and use TLS connection, here we expect to receive bool value in a string
	UseTls DataSourceConnOption = "use_tls"
	// InsecureTls is a REDIS option to set up TLSConfig: &tls.Config{InsecureSkipVerify: true}, here we expect to receive bool value in a string
	InsecureTls DataSourceConnOption = "allow_insecure_tls"
	// Namespace is a prefix before any key
	Namespace DataSourceConnOption = "namespace"
)

var (
	SourceISEmpty error = errors.New("field source (path to file or conn str to db) is empty")
)

// DataSourceConfig represent source where we can get
/*
 * We attempt to provide config that easily could be used with any datasource:
 * - json file (simplest RO mode)
 * - mongodb (but here we have very simple question how to pass parameters)
 * Source contains:
 * 1) if Type is FILE - full path to Json File
 * 2) if Type is REDIS - redis server address i.e. localhost:6739
 * Options are connection options, see - https://www.mongodb.com/docs/drivers/go/current/fundamentals/connection/#std-label-golang-connection-guide
 * Here we should have Validator too
 * Credentials contains Username && Password could be null id authorization is not required:
 *
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
		// 1. Check whether source contains address or not
		// 1.1 Check format
		parts := strings.Split(cfg.Source, ":")
		if len(parts) != 2 {
			return errors.New("field source for Redis datasource must contain pair IP Address/Domain Name:Port, i.e. 127.0.0.1:6379")
		}
		// 1.2
		// todo(UMV): check IP Address / Domain Name is Valid
		_, err := strconv.Atoi(parts[1])
		if err != nil {
			return errors.New(sf.Format("second part must be integer value, got parsing error: {0}", err.Error()))
		}
		// 1.3 Check we have following required fields: dbNumber
		dbNumber, ok := cfg.Options[DbNumber]
		if !ok {
			return errors.New("config must contain \"db_number\" in options")
		}
		checkResult := validators.IsStrValueOfRequiredType(validators.Integer, &dbNumber)
		if !checkResult {
			return errors.New("\"db_number\" redis config options must be int value")
		}
		return nil
	}
	if cfg.Type == MONGODB {
		// validate options values ...
		/*allParamValidation := map[string]string{}
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

		}*/
		return errors.New("MongoDB is not supported")
	}
	return nil
}
