package redis

import (
	"encoding/json"
	"errors"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	appErrs "github.com/wissance/Ferrum/errors"
	sf "github.com/wissance/stringFormatter"
)

// GetUserFederationConfig return data.UserFederationServiceConfig of configured Federation service
/* This function constructs Redis key by pattern combines namespace and realm name and config name (realmUserFederationService)
 * Parameters:
 *     - realmName - name of a Realm
 *     - configName - name of a User Federation Service config
 * Returns: client and error
 */
func (mn *RedisDataManager) GetUserFederationConfig(realmName string, configName string) (*data.UserFederationServiceConfig, error) {
	if !mn.IsAvailable() {
		return nil, appErrs.NewDataProviderNotAvailable(string(config.REDIS), mn.redisOption.Addr)
	}

	userFederationServiceConfigKey := sf.Format(realmUserFederationService, mn.namespace, realmName, configName)
	userFederationConfig, err := getSingleRedisObject[data.UserFederationServiceConfig](mn.redisClient, mn.ctx, mn.logger, RealmUserFederationConfig,
		userFederationServiceConfigKey)
	if err != nil {
		return nil, err
	}
	return userFederationConfig, nil
}

// CreateUserFederationConfig creates new data.UserFederationServiceConfig related to data.Realm by name
/* This function constructs Redis key by pattern combines namespace and realm name and config name (realmUserFederationService)
 * and creates config, unlike Users or Clients number of UserFederationConfig is not big, therefore we don't create a new sub-storage
 * Parameters:
 *     - realmName - name of a Realm
 *     - userFederationConfig - newly creating object data.UserFederationServiceConfig
 * Returns: client and error
 */
func (mn *RedisDataManager) CreateUserFederationConfig(realmName string, userFederationConfig data.UserFederationServiceConfig) error {
	if !mn.IsAvailable() {
		return appErrs.NewDataProviderNotAvailable(string(config.REDIS), mn.redisOption.Addr)
	}
	_, err := mn.getRealmObject(realmName)
	if err != nil {
		return err
	}
	// TODO(UMV): use function isExists
	_, err = mn.GetUserFederationConfig(realmName, userFederationConfig.Name)
	if err == nil {
		return appErrs.NewObjectExistsError(RealmUserFederationConfig, userFederationConfig.Name, sf.Format("realm: {0}", realmName))
	}
	if !errors.As(err, &appErrs.ObjectNotFoundError{}) {
		return err
	}

	userFederationConfigBytes, err := json.Marshal(userFederationConfig)
	if err != nil {
		mn.logger.Error(sf.Format("An error occurred during Marshal Client: {0}", err.Error()))
		return appErrs.NewUnknownError("json.Marshal", "RedisDataManager.CreateUserFederationConfig", err)
	}
	err = mn.upsertUserFederationConfigObject(realmName, userFederationConfig.Name, string(userFederationConfigBytes))
	if err != nil {
		return appErrs.NewUnknownError("upsertClientObject", "RedisDataManager.CreateUserFederationConfig", err)
	}

	return nil
}

func (mn *RedisDataManager) UpdateUserFederationConfig(realmName string, configName string, userFederationConfig data.UserFederationServiceConfig) error {
	return appErrs.ErrOperationNotImplemented
}

// DeleteUserFederationConfig removes data.UserFederationServiceConfig from storage
/* It simply removes data.UserFederationServiceConfig by key based on realmName + configName
 * Arguments:
 *    - realmName - name of a data.Realm
 *    - configName - name of a data.UserFederationServiceConfig
 * Returns: error
 */
func (mn *RedisDataManager) DeleteUserFederationConfig(realmName string, configName string) error {
	if !mn.IsAvailable() {
		return appErrs.NewDataProviderNotAvailable(string(config.REDIS), mn.redisOption.Addr)
	}

	if err := mn.deleteUserFederationConfigObject(realmName, configName); err != nil {
		if errors.As(err, &appErrs.EmptyNotFoundErr) {
			return err
		}
		return appErrs.NewUnknownError("deleteUserFederationConfigObject", "RedisDataManager.DeleteUserFederationConfig", err)
	}

	return nil
}

// upsertUserFederationConfigObject - create or update a data.UserFederationServiceConfig
/* If such a key exists, the value will be overwritten without error
 * Arguments:
 *    - realmName name of a data.Realm
 *    - userFederationConfigName name of a data.UserFederationServiceConfig that is using as a Unique Identifier among Realm
 *    - userFederationJson - string with serialized (Marshalled object)
 * Returns: error
 */
func (mn *RedisDataManager) upsertUserFederationConfigObject(realmName string, userFederationConfigName string, userFederationJson string) error {
	configKey := sf.Format(realmUserFederationService, mn.namespace, realmName, userFederationConfigName)
	if err := mn.upsertRedisString(Client, configKey, userFederationJson); err != nil {
		return appErrs.NewUnknownError("upsertRedisString", "RedisDataManager.upsertUserFederationConfigObject", err)
	}
	return nil
}

// deleteUserFederationConfigObject - deleting a data.UserFederationServiceConfig
/* Inside uses realmUserFederationService
 * Arguments:
 *    - realmName - name of data.Realm
 *    - configName - name of data.UserFederationServiceConfig
 * Returns: error
 */
func (mn *RedisDataManager) deleteUserFederationConfigObject(realmName string, configName string) error {
	configKey := sf.Format(realmUserFederationService, mn.namespace, realmName, configName)
	if err := mn.deleteRedisObject(RealmUserFederationConfig, configKey); err != nil {
		if errors.As(err, &appErrs.EmptyNotFoundErr) {
			return err
		}
		return appErrs.NewUnknownError("deleteRedisObject", "RedisDataManager.deleteUserFederationConfigObject", err)
	}
	return nil
}
