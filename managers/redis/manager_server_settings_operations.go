package redis

import (
	"encoding/json"
	"errors"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	appErrs "github.com/wissance/Ferrum/errors"
	sf "github.com/wissance/stringFormatter"
)

// GetServerSettings function that returns ServerSettings
/* ServerSettings contains main settings that affects whole AuthorizationServer
 * This function reads ServerSettings and return where they are required
 * Arguments: no
 * Returns: *ServerSettings, error
 */
func (mn *RedisDataManager) GetServerSettings() (*data.ServerSettings, error) {
	if !mn.IsAvailable() {
		return nil, appErrs.NewDataProviderNotAvailable(string(config.REDIS), mn.redisOption.Addr)
	}
	return mn.getServerSettingsObject()
}

// SetServerSettings function that updates ServerSettings by full new settings replace
/* This function perform all settings update at once
 * If Admin has no Id it generated in this func
 * Arguments:
 *  - settings - must contain all settings
 * Returns: error (nil if there was no error)
 */
func (mn *RedisDataManager) SetServerSettings(settings *data.ServerSettings) error {
	if settings == nil {
		return appErrs.ErrBadData
	}

	// todo(UMV): check uuid is Empty, if empty -> generate new

	jsonServerSettings, err := json.Marshal(*settings)
	if err != nil {
		return appErrs.ErrBadData
	}

	return mn.upsertServerSettings(string(jsonServerSettings))
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

// upsertRealmObject - create or update a realm without clients and users
/* If such a key exists, the value will be overwritten without error
 * Arguments:
 *    - realmName
 *    - realmJson - string
 * Returns: *Realm, error
 */
func (mn *RedisDataManager) upsertServerSettings(serverSettingsJson string) error {
	realmKey := sf.Format(serverSettingsKeyTemplate, mn.namespace)
	if err := mn.upsertRedisString(ServerSettings, realmKey, serverSettingsJson); err != nil {
		return appErrs.NewUnknownError("upsertRedisString", "RedisDataManager.upsertServerSettings", err)
	}
	return nil
}
