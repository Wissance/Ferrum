package main

import (
	"fmt"

	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/Ferrum/managers/redis"
)

type dataContext interface {
	GetRealm(realmName string) (*data.Realm, error)
	GetClient(realmName string, clientName string) (*data.Client, error)
	GetUser(realmName string, userName string) (data.User, error)

	CreateRealm(newRealm data.Realm) error
	CreateClient(realmName string, clientNew data.Client) error
	CreateUser(realmName string, userNew data.User) error

	DeleteRealm(realmName string) error
	DeleteClient(realmName string, clientName string) error
	DeleteUser(realmName string, userName string) error

	UpdateRealm(realmName string, realmNew data.Realm) error
	UpdateClient(realmName string, clientName string, clientNew data.Client) error
	UpdateUser(realmName string, userName string, userNew data.User) error

	SetPassword(realmName string, userName string, password string) error
}

var _ dataContext = (*redis.RedisDataManager)(nil)

func prepareContext(dataSourceCfg *config.DataSourceConfig, logger *logging.AppLogger) (dataContext, error) {
	var manager dataContext
	var err error
	switch dataSourceCfg.Type {
	case config.FILE:
		err = fmt.Errorf("not supported")
	case config.REDIS:
		manager, err = redis.CreateRedisDataManager(dataSourceCfg, logger)
	default:
		err = fmt.Errorf("not supported")
	}
	return manager, err
}
