package managers

import (
	"Ferrum/data"
	"github.com/google/uuid"
	"path/filepath"
)

type FileDataManager struct {
	dataFile string
}

func Create(dataFile string) DataContext {
	absPath, err := filepath.Abs(dataFile)
	if err != nil {

	}
	// init, load data in memory ...
	mn := &FileDataManager{dataFile: absPath}

	dc := DataContext(mn)
	return dc
}

func (mn *FileDataManager) GetRealm(realmId *string) *data.Realm {
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
