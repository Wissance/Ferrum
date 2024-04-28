package redis

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	sf "github.com/wissance/stringFormatter"
	"testing"
)

const testUser = "ferrum_db"
const testUserPassword = "FeRRuM000"
const testRedisSource = "127.0.0.1:6379"

func TestCreateRealmSuccessfully(t *testing.T) {
	testCases := []struct {
		name    string
		realm   string
		clients []string
	}{
		{name: "realm_without_clients", realm: "app1_test", clients: []string{}},
	}
	manager := createTestRedisDataManager()
	for _, tCase := range testCases {
		realm := data.Realm{
			Name:                   tCase.realm,
			TokenExpiration:        3600,
			RefreshTokenExpiration: 1800,
		}

		err := manager.CreateRealm(realm)
		assert.NoError(t, err)
		r, err := manager.GetRealm(tCase.realm)
		assert.NoError(t, err)
		// TODO(UMV): IMPL FULL COMPARISON
		assert.Equal(t, tCase.realm, r.Name)
		err = manager.DeleteRealm(tCase.realm)
		assert.NoError(t, err)
	}
}

func createTestRedisDataManager() *RedisDataManager {
	rndNamespace := sf.Format("ferrum_test_{0}", uuid.New().String())
	dataSourceCfg := config.DataSourceConfig{
		Type:   config.REDIS,
		Source: testRedisSource,
		Options: map[config.DataSourceConnOption]string{
			config.Namespace: rndNamespace,
			config.DbNumber:  "0",
		},
		Credentials: &config.CredentialsConfig{
			Username: testUser,
			Password: testUserPassword,
		},
	}

	loggerCfg := config.LoggingConfig{}

	logger := logging.CreateLogger(&loggerCfg)
	manager, _ := CreateRedisDataManager(&dataSourceCfg, logger)
	return manager
}
