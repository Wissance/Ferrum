package managers

import (
	"errors"
	"github.com/google/uuid"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/stringFormatter"
	"path/filepath"
)

// DataContext is a common interface to implement operations with authorization server entities (data.Realm, data.Client, data.User)
// now contains only set of Get methods, during implementation admin CLI should be expanded to create && update entities
type DataContext interface {
	// GetRealm returns realm by name (unique)
	GetRealm(realmName string) *data.Realm
	// GetClient returns realm client by name (client name is also unique in a realm)
	GetClient(realm *data.Realm, name string) *data.Client
	// GetUser return realm user (consider what to do with Federated users) by name
	GetUser(realm *data.Realm, userName string) *data.User
	// GetUserById return realm user by id
	GetUserById(realm *data.Realm, userId uuid.UUID) *data.User
	// GetRealmUsers return all realm Users
	// TODO(UMV): when we deal with a lot of Users we should query portion of Users instead of all
	GetRealmUsers(realmName string) *[]data.User
}

// PrepareContext is a factory function that creates instance of DataContext
/* This function creates instance of appropriate DataContext according to input arguments values, if dataSourceConfig is config.FILE function
 * creates instance of FileDataManager. For this type of context if dataFile is not nil and exists this function also provides data initialization:
 * loads all data (realms, clients and users) in a memory. If dataSourceCfg is config.REDIS this function creates instance of RedisDataManager
 * by calling CreateRedisDataManager function
 * Parameters:
 *     - dataSourceCfg configuration section related to DataSource
 *     - dataFile - data for initialization (this is using only when dataSourceCfg is config.FILE)
 *     - logger - logger instance
 * Return: new instance of DataContext and error (nil if there are no errors)
 */
func PrepareContext(dataSourceCfg *config.DataSourceConfig, dataFile *string, logger *logging.AppLogger) (DataContext, error) {
	var dc DataContext
	var err error
	if dataSourceCfg.Type == config.FILE {
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
		mn := &FileDataManager{dataFile: absPath, logger: logger}
		err = mn.loadData()
		if err != nil {
			// at least and think what to do further
			msg := stringFormatter.Format("An error occurred during data loading: {0}", err.Error())
			logger.Error(msg)
		}
		dc = DataContext(mn)

	} else {
		if dataSourceCfg.Type == config.REDIS {
			dc, err = CreateRedisDataManager(dataSourceCfg, logger)
		}
		// todo implement other data sources
	}
	return dc, err
}
