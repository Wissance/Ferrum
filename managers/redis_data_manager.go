package managers

import (
	"crypto/tls"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	sf "github.com/wissance/stringFormatter"
	"strconv"
)

type RedisDataManager struct {
	redisOption *redis.Options
	redisClient *redis.Client
	logger      *logging.AppLogger
}

func CreateRedisDataManager(dataSourceCfd *config.DataSourceConfig, logger *logging.AppLogger) DataContext {
	opts := buildRedisConfig(dataSourceCfd, logger)
	rClient := redis.NewClient(opts)
	mn := &RedisDataManager{logger: logger, redisOption: opts, redisClient: rClient}
	// todo(umv) think about preload ???
	dc := DataContext(mn)
	return dc
}

func (mn *RedisDataManager) GetRealm(realmName string) *data.Realm {
	return nil
}

func (mn *RedisDataManager) GetClient(realm *data.Realm, name string) *data.Client {
	return nil
}

func (mn *RedisDataManager) GetUser(realm *data.Realm, userName string) *data.User {
	return nil
}

func (mn *RedisDataManager) GetUserById(realm *data.Realm, userId uuid.UUID) *data.User {
	return nil
}

func (mn *RedisDataManager) GetRealmUsers(realmName string) *[]data.User {
	return nil
}

func buildRedisConfig(dataSourceCfd *config.DataSourceConfig, logger *logging.AppLogger) *redis.Options {
	dbNum, err := strconv.Atoi(dataSourceCfd.Options[config.DbNumber])
	if err != nil {
		logger.Error(sf.Format("can't be because we already called Validate(), but in any case: parsing error: {0}", err.Error()))
		return nil
	}
	opts := redis.Options{
		Addr: dataSourceCfd.Source,
		DB:   dbNum,
	}
	// passing credentials if we have it
	if dataSourceCfd.Credentials != nil {
		opts.Username = dataSourceCfd.Credentials.Username
		opts.Password = dataSourceCfd.Credentials.Password
	}
	// passing TLS if we have it
	val, ok := dataSourceCfd.Options[config.UseTls]
	if ok {
		useTls, parseErr := strconv.ParseBool(val)
		if parseErr == nil && useTls {
			opts.TLSConfig = &tls.Config{}
			val, ok = dataSourceCfd.Options[config.InsecureTls]
			if ok {
				inSecTls, parseErr := strconv.ParseBool(val)
				if parseErr == nil {
					opts.TLSConfig.InsecureSkipVerify = inSecTls
				}
			}
		}
	}

	return &opts
}
