package managers

import (
	"Ferrum/data"
	"encoding/json"
	"fmt"
	"github.com/wissance/stringFormatter"
	"io/ioutil"
	"path/filepath"
)

type FileDataManager struct {
	dataFile   string
	serverData data.ServerData
}

func Create(dataFile string) DataContext {
	absPath, err := filepath.Abs(dataFile)
	if err != nil {
		fmt.Println(stringFormatter.Format("An error occurred during attempt to get abs path of data file: {0}", err.Error()))
	}
	// init, load data in memory ...
	mn := &FileDataManager{dataFile: absPath}
	err = mn.loadData()
	if err != nil {
		// at least and think what to do further
		fmt.Println(stringFormatter.Format("An error occurred during data loading: {0}", err.Error()))
	}
	dc := DataContext(mn)
	return dc
}

func (mn *FileDataManager) GetRealm(realmName string) *data.Realm {
	for _, e := range mn.serverData.Realms {
		// case-sensitive comparison, myapp and MyApP are different realms
		if e.Name == realmName {
			return &e
		}
	}
	return nil
}

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

func (mn *FileDataManager) loadData() error {
	rawData, err := ioutil.ReadFile(mn.dataFile)
	if err != nil {
		fmt.Println(stringFormatter.Format("An error occurred during config file reading: {0}", err.Error()))
		return err
	}
	mn.serverData = data.ServerData{}
	if err = json.Unmarshal(rawData, &mn.serverData); err != nil {
		fmt.Println(stringFormatter.Format("An error occurred during data file unmarshal: {0}", err.Error()))
		return err
	}

	return nil
}
