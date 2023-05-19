package config

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"testing"
)

func TestValidateAppConfigWithRedisDataSourceCfg(t *testing.T) {
	testCases := []struct {
		name                string
		cfgFile             string
		expectedSource      string
		expectedCredentials *CredentialsConfig
		dbNumber            string
	}{
		{
			name:           "MinimalValidAppCfgWithRedis",
			cfgFile:        path.Join("test_configs", "valid_config_w_min_redis.json"),
			expectedSource: "localhost:6380",
			expectedCredentials: &CredentialsConfig{
				Username: "dbmn",
				Password: "123",
			},
			dbNumber: "12",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fileData, err := ioutil.ReadFile(tc.cfgFile)
			assert.NoError(t, err)
			appConfig := AppConfig{}
			err = json.Unmarshal(fileData, &appConfig)
			assert.NoError(t, err)
			// check values itself
			checkDataSourceValues(t, REDIS, tc.expectedSource, tc.expectedCredentials,
				map[DataSourceConnOption]string{
					DbNumber: tc.dbNumber,
				}, appConfig.DataSource)
			err = appConfig.DataSource.Validate()
			assert.NoError(t, err)
		})
	}
}

func checkDataSourceValues(t *testing.T, expectedSourceType DataSourceType, expectedSource string, expectedCredentials *CredentialsConfig,
	expectedOptions map[DataSourceConnOption]string, actualCfg DataSourceConfig) {
	assert.Equal(t, expectedSourceType, actualCfg.Type)
	assert.Equal(t, expectedSource, actualCfg.Source)
	if expectedCredentials == nil {
		assert.Nil(t, actualCfg.Credentials)
	} else {
		assert.NotNil(t, actualCfg.Credentials)
		assert.Equal(t, expectedCredentials.Username, actualCfg.Credentials.Username)
		assert.Equal(t, expectedCredentials.Password, actualCfg.Credentials.Password)
	}
	assert.Equal(t, len(expectedOptions), len(actualCfg.Options))
	for k, v := range expectedOptions {
		av, ok := actualCfg.Options[k]
		assert.True(t, ok)
		assert.Equal(t, v, av)
	}
}
