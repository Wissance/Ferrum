package managers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/stringFormatter"
	"io/ioutil"
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
func PrepareFileDataContextUsingData(serverData *data.ServerData) (DataContext, error) {
	// todo(UMV): todo provide an error handling
	mn := &FileDataManager{serverData: *serverData}
	dc := DataContext(mn)
	return dc, nil
}

// GetRealm function for getting Realm by name
func (mn *FileDataManager) GetRealm(realmName string) *data.Realm {
	for _, e := range mn.serverData.Realms {
		// case-sensitive comparison, myapp and MyApP are different realms
		if e.Name == realmName {
			return &e
		}
	}
	return nil
}

// GetClient function for getting Realm Client by name
func (mn *FileDataManager) GetClient(realm *data.Realm, name string) *data.Client {
	for _, c := range realm.Clients {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

func (mn *FileDataManager) GetUser(realm *data.Realm, userName string) *data.User {
	for _, u := range realm.Users {
		user := data.CreateUser(u)
		if user.GetUsername() == userName {
			return &user
		}
	}

	return nil
}

func (mn *FileDataManager) GetUserById(realm *data.Realm, userId uuid.UUID) *data.User {
	for _, u := range realm.Users {
		user := data.CreateUser(u)
		if user.GetId() == userId {
			return &user
		}
	}

	return nil
}

func (mn *FileDataManager) GetRealmUsers(realmName string) *[]data.User {
	realm := mn.GetRealm(realmName)
	if realm == nil {
		return nil
	}
	users := make([]data.User, len(realm.Users))
	for i, u := range realm.Users {
		users[i] = data.CreateUser(u)
	}
	return &users
}

// loadData this function loads data from JSON file (dataFile) to serverData
func (mn *FileDataManager) loadData() error {
	rawData, err := ioutil.ReadFile(mn.dataFile)
	if err != nil {
		mn.logger.Error(stringFormatter.Format("An error occurred during config file reading: {0}", err.Error()))
		return err
	}
	mn.serverData = data.ServerData{}
	if err = json.Unmarshal(rawData, &mn.serverData); err != nil {
		mn.logger.Error(stringFormatter.Format("An error occurred during data file unmarshal: {0}", err.Error()))
		return err
	}

	return nil
}
