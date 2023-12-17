package managers

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/wissance/Ferrum/managers/files"
	"github.com/wissance/Ferrum/managers/redis"

	"github.com/google/uuid"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/stringFormatter"
)

// DataContext is a common interface to implement operations with authorization server entities (data.Realm, data.Client, data.User)
// now contains only set of Get methods, during implementation admin CLI should be expanded to create && update entities
type DataContext interface {
	// GetReal returns realm by name (unique) returns realm with clients but no users
	GetRealm(realmName string) (*data.Realm, error)
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

// PrepareContextUsingData is a factory function that creates instance of DataContext
/* This function creates instance of appropriate DataContext according to input arguments values, if dataSourceConfig is config.FILE function
 * creates instance of FileDataManager.
 * loads all data (realms, clients and users) in a memory.
 * Parameters:
 *     - dataSourceCfg configuration section related to DataSource
 *     - data - ServerData
 *     - logger - logger instance
 * Return: new instance of DataContext and error (nil if there are no errors)
 */
func PrepareContextUsingData(dataSourceCfg *config.DataSourceConfig, data *data.ServerData, logger *logging.AppLogger) (DataContext, error) {
	var dc DataContext
	var err error
	switch dataSourceCfg.Type {
	case config.FILE:
		dc, err = files.CreateFileDataManagerWithInitData(data)

	case config.REDIS:
		return nil, fmt.Errorf("not supported initialization with init data")

	default:
		return nil, fmt.Errorf("not supported")
	}

	return dc, err
}

// PrepareContextUsingFile is a factory function that creates instance of DataContext
/* This function creates instance of appropriate DataContext according to input arguments values, if dataSourceConfig is config.FILE function
 * creates instance of FileDataManager. For this type of context if dataFile is not nil and exists this function also provides data initialization:
 * loads all data (realms, clients and users) in a memory.
 * Parameters:
 *     - dataSourceCfg configuration section related to DataSource
 *     - dataFile - data for initialization (this is using only when dataSourceCfg is config.FILE)
 *     - logger - logger instance
 * Return: new instance of DataContext and error (nil if there are no errors)
 */
func PrepareContextUsingFile(dataSourceCfg *config.DataSourceConfig, dataFile *string, logger *logging.AppLogger) (DataContext, error) {
	if dataFile == nil {
		return nil, errors.New("data file is nil")
	}
	var dc DataContext
	var err error
	switch dataSourceCfg.Type {
	case config.FILE:
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
		return nil, fmt.Errorf("not supported initialization with init data")

	default:
		return nil, fmt.Errorf("not supported")
	}

	return dc, err
}

// PrepareContext is a factory function that creates instance of DataContext
/* If dataSourceCfg is config.REDIS this function creates instance of RedisDataManager by calling CreateRedisDataManager function
 * Parameters:
 *     - dataSourceCfg configuration section related to DataSource
 *     - logger - logger instance
 * Return: new instance of DataContext and error (nil if there are no errors)
 */
func PrepareContext(dataSourceCfg *config.DataSourceConfig, logger *logging.AppLogger) (DataContext, error) {
	var dc DataContext
	var err error
	switch dataSourceCfg.Type {
	case config.FILE:
		return nil, fmt.Errorf("not supported initialization without init data")

	case config.REDIS:
		dc, err = redis.CreateRedisDataManager(dataSourceCfg, logger)

	default:
		return nil, fmt.Errorf("not supported")
	}
	// todo implement other data sources

	return dc, err
}
