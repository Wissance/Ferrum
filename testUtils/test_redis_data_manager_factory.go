package testUtils

import (
	"github.com/google/uuid"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/Ferrum/managers/redis"
	sf "github.com/wissance/stringFormatter"
)

func CreateTestRedisDataManager(redisUri string, user string, password string) (*redis.RedisDataManager, error) {
	rndNamespace := sf.Format("ferrum_test_{0}", uuid.New().String())
	dataSourceCfg := config.DataSourceConfig{
		Type:   config.REDIS,
		Source: redisUri,
		Options: map[config.DataSourceConnOption]string{
			config.Namespace: rndNamespace,
			config.DbNumber:  "0",
		},
		Credentials: &config.CredentialsConfig{
			Username: user,
			Password: password,
		},
	}

	loggerCfg := config.LoggingConfig{}

	logger := logging.CreateLogger(&loggerCfg)
	manager, err := redis.CreateRedisDataManager(&dataSourceCfg, logger)
	return manager, err
}
