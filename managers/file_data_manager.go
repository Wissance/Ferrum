package managers

import (
	"Ferrum/data"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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

func (mn *FileDataManager) GetRealm(realmId string) *data.Realm {
	for _, e := range mn.serverData.Realms {
		// todo(umv): should we compare in lower case
		if e.Name == realmId {
			return &e
		}
	}
	return nil
}

func (mn *FileDataManager) GetClient(clientId uuid.UUID) *data.Client {
	return nil
}

func (mn *FileDataManager) GetUser(userId uuid.UUID) *data.User {
	return nil
}

func (mn *FileDataManager) GetClientUsers(clientId uuid.UUID) *[]data.User {
	return nil
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
