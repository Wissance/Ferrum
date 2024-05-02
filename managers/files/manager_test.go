package files

import (
	"github.com/stretchr/testify/require"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/logging"
	"testing"
)

const testDataFile = "test_data.json"

func TestGetRealmSuccessfully(t *testing.T) {

}

func TestGetClientSuccessfully(t *testing.T) {

}

func TestGetClientsSuccessfully(t *testing.T) {

}

func TestGetClientsSuccessfullyEmptyRealm(t *testing.T) {

}

func TestGetClientsSuccessfullyNonExistingRealm(t *testing.T) {

}

func TestGetUsersSuccessfully(t *testing.T) {

}

func TestGetUserSuccessfully(t *testing.T) {

}

func TestGetUserByIdSuccessfully(t *testing.T) {

}

func TestGetUsersSuccessfullyEmptyRealm(t *testing.T) {

}

func TestGetUsersSuccessfullyNonExistingRealm(t *testing.T) {

}

func TestGetClientSuccessfullyNonExistingRealm(t *testing.T) {

}

func createTestFileDataManager(t *testing.T) *FileDataManager {
	loggerCfg := config.LoggingConfig{}

	logger := logging.CreateLogger(&loggerCfg)

	manager, err := CreateFileDataManager(testDataFile, logger)
	require.NoError(t, err)
	return manager
}
