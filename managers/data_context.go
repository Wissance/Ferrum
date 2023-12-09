package managers

import (
	"errors"
	"fmt"
	"github.com/wissance/Ferrum/managers/files"
	"github.com/wissance/Ferrum/managers/redis"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/stringFormatter"
)

// DataContext is a common interface to implement operations with authorization server entities (data.Realm, data.Client, data.User)
// now contains only set of Get methods, during implementation admin CLI should be expanded to create && update entities
type DataContext interface {
	// GetRealmWithClients returns realm by name (unique) with all clients
	GetRealmWithClients(realmName string) (*data.Realm, error)
	// GetClient returns realm client by name (client name is also unique in a realm)
	GetClient(realmName string, name string) (*data.Client, error)
	// GetUser return realm user (consider what to do with Federated users) by name
	GetUser(realmName string, userName string) (data.User, error)
	// GetUserById return realm user by id
	GetUserById(realmName string, userId uuid.UUID) (data.User, error)
	// GetRealmUsers return all realm Users
	// TODO(UMV): when we deal with a lot of Users we should query portion of Users instead of all
	// GetRealmUsers(realmName string) ([]data.User, error)
}

func PrepareContext(dataSourceCfg *config.DataSourceConfig, dataFile *string, logger *logging.AppLogger) (DataContext, error) {
	var dc DataContext
	var err error
	switch dataSourceCfg.Type {
	case config.FILE:
		if dataFile == nil {
			err = errors.New("data file is nil")
		}
		absPath, pathErr := filepath.Abs(*dataFile)
		if pathErr != nil {
			// todo: umv: think what to do on error
			msg := stringFormatter.Format("An error occurred during attempt to get abs path of data file: {0}", err.Error())
			logger.Error(msg)
			err = pathErr
		}
		// init, load data in memory ...
		mn, err := files.CreateFileDataManager(absPath, logger)
		if err != nil {
			// at least and think what to do further
			msg := stringFormatter.Format("An error occurred during data loading: {0}", err.Error())
			logger.Error(msg)
		}
		dc = DataContext(mn)

	case config.REDIS:
		if dataSourceCfg.Type == config.REDIS {
			dc, err = redis.CreateRedisDataManager(dataSourceCfg, logger)
		}
		// todo implement other data sources
	}

	return dc, err
}

func PrepareContextUsingData(dataSourceCfgType config.DataSourceType, data *data.ServerData) (DataContext, error) {
	var dc DataContext
	var err error
	switch dataSourceCfgType {
	case config.FILE:
		mn, err := files.CreateFileDataManagerWithInitData(data)
		if err != nil {
			return nil, fmt.Errorf("CreateFileDataManagerWithInitData failed: %w", err)
		}
		dc = DataContext(mn)

	case config.REDIS:
		return nil, fmt.Errorf("Not supported")
	default:
		return nil, fmt.Errorf("Not supported")
	}

	return dc, err
}
