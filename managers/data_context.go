package managers

import (
	"errors"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/stringFormatter"
)

type DataContext interface {
	GetRealm(realmName string) *data.Realm
	GetClient(realm *data.Realm, name string) *data.Client
	GetUser(realm *data.Realm, userName string) *data.User
	GetUserById(realm *data.Realm, userId uuid.UUID) *data.User
	GetRealmUsers(realmName string) []data.User
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
		mn := &FileDataManager{dataFile: absPath, logger: logger}
		err = mn.loadData()
		if err != nil {
			// at least and think what to do further
			msg := stringFormatter.Format("An error occurred during data loading: {0}", err.Error())
			logger.Error(msg)
		}
		dc = DataContext(mn)

	case config.REDIS:
		if dataSourceCfg.Type == config.REDIS {
			dc, err = CreateRedisDataManager(dataSourceCfg, logger)
		}
		// todo implement other data sources
	}

	return dc, err
}
