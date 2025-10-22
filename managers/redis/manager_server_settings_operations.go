package redis

import (
	"errors"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	appErrs "github.com/wissance/Ferrum/errors"
	sf "github.com/wissance/stringFormatter"
)

// GetServerSettings function that returns ServerSettings
func (mn *RedisDataManager) GetServerSettings() (*data.ServerSettings, error) {
	if !mn.IsAvailable() {
		return nil, appErrs.NewDataProviderNotAvailable(string(config.REDIS), mn.redisOption.Addr)
	}
	return mn.getServerSettingsObject()
}

// SetServerSettings function that updates ServerSettings by full new settings replace
func (mn *RedisDataManager) SetServerSettings(settings *data.ServerSettings) error {

	return nil
}

// getServerSettingsObject - returns server settings
/* Server settings is common settings for the whole Authorization Server
 * Arguments: no
 * Returns: *ServerSettings, error
 */
func (mn *RedisDataManager) getServerSettingsObject() (*data.ServerSettings, error) {
	serverSettingsKey := sf.Format(serverSettingsKeyTemplate, mn.namespace)
	serverSettings, err := getSingleRedisObject[data.ServerSettings](mn.redisClient, mn.ctx, mn.logger, ServerSettings,
		serverSettingsKey)
	if err != nil {
		if errors.As(err, &appErrs.EmptyNotFoundErr) {
			mn.logger.Debug(sf.Format("Redis does not have ServerSettings"))
		}
		return nil, err
	}
	return serverSettings, nil
}
