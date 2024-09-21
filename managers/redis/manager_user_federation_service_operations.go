package redis

import (
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/errors"
)

func (mn *RedisDataManager) GetUserFederationConfig(realmName string, configName string) (*data.UserFederationServiceConfig, error) {
	return nil, errors.ErrOperationNotImplemented
}

func (mn *RedisDataManager) CreateUserFederationConfig(realmName string, userFederationConfig data.UserFederationServiceConfig) error {
	return errors.ErrOperationNotImplemented
}

func (mn *RedisDataManager) UpdateUserFederationConfig(realmName string, configName string, userFederationConfig data.UserFederationServiceConfig) error {
	return errors.ErrOperationNotImplemented
}

func (mn *RedisDataManager) DeleteUserFederationConfig(realmName string, configName string) error {
	return errors.ErrOperationNotImplemented
}
