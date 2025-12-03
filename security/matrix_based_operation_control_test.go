package security

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/testUtils"
	"github.com/wissance/Ferrum/utils/encoding"
	sf "github.com/wissance/stringFormatter"
	"testing"
)

const (
	testUser         = "ferrum_db"
	testUserPassword = "FeRRuM000"
	testRedisSource  = "127.0.0.1:6379"
)

func TestIsOperationAllowedWithRedisDataManager(t *testing.T) {
	manager, err := testUtils.CreateTestRedisDataManager(testRedisSource, testUser, testUserPassword)
	require.NoError(t, err)
	// Init data
	// 1. Create ServerSettings
	// 2. Create 2 Realms with 2 Client, and 3-4 users
	adminUuid, err := uuid.Parse("00000000-0000-0000-0000-000000000001")
	passwordSalt := "123"
	encoder := encoding.NewPasswordJsonEncoder(passwordSalt)
	serverSettings := data.ServerSettings{
		AllowedHosts:      []string{"*"},
		AdminApiUrlPrefix: "1234567890",
		Admin: data.AdminUser{
			Id:           adminUuid,
			Username:     "admin",
			PasswordSalt: passwordSalt,
			PasswordHash: encoder.GetB64PasswordHash(passwordSalt),
		},
	}
	err = manager.SetServerSettings(&serverSettings)
	require.NoError(t, err)

	realm1 := data.Realm{
		Name:                   sf.Format("realm1_{0}", uuid.New().String()),
		TokenExpiration:        3600,
		RefreshTokenExpiration: 1800,
	}

	realm2 := data.Realm{
		Name:                   sf.Format("realm2_{0}", uuid.New().String()),
		TokenExpiration:        3600,
		RefreshTokenExpiration: 1800,
	}

	// 1. Create an empty Realms
	err = manager.CreateRealm(realm1)
	require.NoError(t, err)
	err = manager.CreateRealm(realm2)
	require.NoError(t, err)
	// 2. Create Users
	realm1UsersId := []uuid.UUID{
		getUuidFromStr("00000000-0000-0000-0000-100000000001"),
		getUuidFromStr("00000000-0000-0000-0000-100000000002"),
		getUuidFromStr("00000000-0000-0000-0000-100000000003"),
	}

	realm1Users := make([]data.User, 3)
	for i, v := range realm1UsersId {
		userId := v
		userName := sf.Format("test_user_{0}_{1}", i, userId)
		jsonTemplate := `{"info":{"sub":"{2}", "name":"{0}", "preferred_username": "{1}"}, "credentials":{"password": "123"}}`
		jsonStr := sf.Format(jsonTemplate, userName, userName, userId)
		var rawUser interface{}
		err = json.Unmarshal([]byte(jsonStr), &rawUser)
		assert.NoError(t, err)
		user := data.CreateUser(rawUser, nil)
		err = manager.CreateUser(realm1.Name, user)
		realm1Users[i] = user
		assert.NoError(t, err)
	}
	realm1.Owner = realm1UsersId[0]
	realm1.Admins = []uuid.UUID{realm1UsersId[0], realm1UsersId[1]}
	err = manager.UpdateRealm(realm1.Name, realm1)
	require.NoError(t, err)

	realm1Client1 := data.Client{
		Name: "r1_app_client_1",
		Type: data.Public,
		ID:   getUuidFromStr("00000000-0000-0000-0000-110000000001"),
		Auth: data.Authentication{
			Type:  data.ClientIdAndSecrets,
			Value: uuid.New().String(),
		},
	}
	err = manager.CreateClient(realm1.Name, realm1Client1)
	require.NoError(t, err)

	realm1Client2 := data.Client{
		Name: "r1_app_client_2",
		Type: data.Public,
		ID:   getUuidFromStr("00000000-0000-0000-0000-110000000002"),
		Auth: data.Authentication{
			Type:  data.ClientIdAndSecrets,
			Value: uuid.New().String(),
		},
	}
	err = manager.CreateClient(realm1.Name, realm1Client2)
	require.NoError(t, err)

	realm2UsersId := []uuid.UUID{
		getUuidFromStr("00000000-0000-0000-0000-100000000004"),
		getUuidFromStr("00000000-0000-0000-0000-100000000005"),
		getUuidFromStr("00000000-0000-0000-0000-100000000006"),
	}

	realm2Users := make([]data.User, 3)
	for i, v := range realm2UsersId {
		userId := v
		userName := sf.Format("test_user_{0}_{1}", i, userId)
		jsonTemplate := `{"info":{"sub":"{2}", "name":"{0}", "preferred_username": "{1}"}, "credentials":{"password": "123"}}`
		jsonStr := sf.Format(jsonTemplate, userName, userName, userId)
		var rawUser interface{}
		err = json.Unmarshal([]byte(jsonStr), &rawUser)
		assert.NoError(t, err)
		user := data.CreateUser(rawUser, nil)
		err = manager.CreateUser(realm2.Name, user)
		realm2Users[i] = user
		assert.NoError(t, err)
	}
	realm2.Owner = realm2UsersId[0]
	realm2.Admins = []uuid.UUID{realm2UsersId[0], realm2UsersId[1]}
	err = manager.UpdateRealm(realm2.Name, realm2)
	require.NoError(t, err)

	realm2Client1 := data.Client{
		Name: "r2_app_client_1",
		Type: data.Public,
		ID:   getUuidFromStr("00000000-0000-0000-0000-120000000001"),
		Auth: data.Authentication{
			Type:  data.ClientIdAndSecrets,
			Value: uuid.New().String(),
		},
	}
	err = manager.CreateClient(realm2.Name, realm2Client1)
	require.NoError(t, err)

	realm2Client2 := data.Client{
		Name: "r2_app_client_2",
		Type: data.Public,
		ID:   getUuidFromStr("00000000-0000-0000-0000-120000000002"),
		Auth: data.Authentication{
			Type:  data.ClientIdAndSecrets,
			Value: uuid.New().String(),
		},
	}
	err = manager.CreateClient(realm2.Name, realm2Client2)
	require.NoError(t, err)

	// test data was created - 2 Realms with 2 Clients each and 3 Users

	err = manager.DeleteRealm(realm1.Name)
	assert.NoError(t, err)
	err = manager.DeleteRealm(realm2.Name)
	assert.NoError(t, err)
}

func getUuidFromStr(str string) uuid.UUID {
	id, _ := uuid.Parse(str)
	return id
}
