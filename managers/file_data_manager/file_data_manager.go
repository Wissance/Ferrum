package file_data_manager

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/Ferrum/managers/errors_managers"
	"github.com/wissance/stringFormatter"
)

// FileDataManager is the simplest Data Storage without any dependencies, it uses single JSON file (it is users and clients RO auth server)
// This context type is extremely useful for simple systems
type FileDataManager struct {
	dataFile   string
	serverData data.ServerData
	logger     *logging.AppLogger
}

// PrepareFileDataContextUsingData initializes instance of FileDataManager and sets loaded data to serverData
/* This factory function creates initialize with data instance of  FileDataManager, error reserved for usage but always nil here
 * Parameters:
 *    serverData already loaded data.ServerData from Json file in memory
 * Returns: context and error (currently is nil)
 */
func CreateFileDataManagerUsingData(serverData *data.ServerData) (*FileDataManager, error) {
	// todo(UMV): todo provide an error handling
	mn := &FileDataManager{serverData: *serverData}
	return mn, nil
}

func CreateFileDataManager(dataFile string, logger *logging.AppLogger) (*FileDataManager, error) {
	// todo(UMV): todo provide an error handling
	mn := &FileDataManager{dataFile: dataFile, logger: logger}
	if err := mn.loadData(); err != nil {
		return nil, fmt.Errorf("loadData failed: %w", err)
	}
	return mn, nil
}

// GetRealm function for getting realm by name, returns the realm without clients and users.
func (mn *FileDataManager) GetRealm(realmName string) (*data.Realm, error) {
	for _, e := range mn.serverData.Realms {
		// case-sensitive comparison, myapp and MyApP are different realms
		if e.Name == realmName {
			e.Clients = nil
			e.Users = nil
			return &e, nil
		}
	}
	return nil, errors_managers.ErrNotFound
}

func (mn *FileDataManager) GetRealmWithClients(realmName string) (*data.Realm, error) {
	for _, e := range mn.serverData.Realms {
		// case-sensitive comparison, myapp and MyApP are different realms
		if e.Name == realmName {
			e.Users = nil
			return &e, nil
		}
	}
	return nil, errors_managers.ErrNotFound
}

func (mn *FileDataManager) GetRealmWithUsers(realmName string) (*data.Realm, error) {
	for _, e := range mn.serverData.Realms {
		// case-sensitive comparison, myapp and MyApP are different realms
		if e.Name == realmName {
			e.Clients = nil
			return &e, nil
		}
	}
	return nil, errors_managers.ErrNotFound
}

func (mn *FileDataManager) GetClient(realmName string, clientName string) (*data.Client, error) {
	realm, err := mn.GetRealmWithClients(realmName)
	if err != nil {
		return nil, fmt.Errorf("GetRealmWithClients failed: %w", err)
	}

	for _, c := range realm.Clients {
		if c.Name == clientName {
			return &c, nil
		}
	}
	return nil, errors_managers.ErrNotFound
}

func (mn *FileDataManager) GetUser(realmName string, userName string) (data.User, error) {
	realm, err := mn.GetRealmWithUsers(realmName)
	if err != nil {
		return nil, fmt.Errorf("GetRealmWithUsers failed: %w", err)
	}

	for _, u := range realm.Users {
		user := data.CreateUser(u)
		if user.GetUsername() == userName {
			return user, nil
		}
	}

	return nil, errors_managers.ErrNotFound
}

func (mn *FileDataManager) GetUserById(realmName string, userId uuid.UUID) (data.User, error) {
	realm, err := mn.GetRealmWithUsers(realmName)
	if err != nil {
		return nil, fmt.Errorf("GetRealmWithUsers failed: %w", err)
	}

	for _, u := range realm.Users {
		user := data.CreateUser(u)
		if user.GetId() == userId {
			return user, nil
		}
	}

	return nil, errors_managers.ErrNotFound
}

func (mn *FileDataManager) GetUsers(realmName string) ([]data.User, error) {
	realm, err := mn.GetRealmWithUsers(realmName)
	if err != nil {
		return nil, fmt.Errorf("GetRealmWithUsers failed: %w", err)
	}

	users := make([]data.User, len(realm.Users))
	for i, u := range realm.Users {
		users[i] = data.CreateUser(u)
	}
	return users, nil
}

// loadData this function loads data from JSON file (dataFile) to serverData
func (mn *FileDataManager) loadData() error {
	rawData, err := os.ReadFile(mn.dataFile)
	if err != nil {
		mn.logger.Error(stringFormatter.Format("An error occurred during config file reading: {0}", err.Error()))
		return fmt.Errorf("os.ReadFile failed: %w", err)
	}
	mn.serverData = data.ServerData{}
	if err = json.Unmarshal(rawData, &mn.serverData); err != nil {
		mn.logger.Error(stringFormatter.Format("An error occurred during data file unmarshal: {0}", err.Error()))
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	return nil
}
