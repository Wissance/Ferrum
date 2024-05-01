package redis

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	appErrs "github.com/wissance/Ferrum/errors"
	"github.com/wissance/Ferrum/logging"
	sf "github.com/wissance/stringFormatter"
	"testing"
)

const testUser = "ferrum_db"
const testUserPassword = "FeRRuM000"
const testRedisSource = "127.0.0.1:6379"

func TestCreateRealmSuccessfully(t *testing.T) {
	testCases := []struct {
		name              string
		realmNameTemplate string
		clients           []string
		users             []string
	}{
		{name: "realm_without_clients", realmNameTemplate: "app1_test_{0}", clients: []string{}, users: []string{}},
		{name: "realm_with_one_client", realmNameTemplate: "app2_test_{0}", clients: []string{"app_client2"}, users: []string{}},
		{name: "realm_with_one_client_and_one_user", realmNameTemplate: "app3_test_{0}", clients: []string{"app_client3"}, users: []string{"app_user3"}},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()
			manager := createTestRedisDataManager()
			realm := data.Realm{
				Name:                   sf.Format(tCase.realmNameTemplate, uuid.New().String()),
				TokenExpiration:        3600,
				RefreshTokenExpiration: 1800,
			}

			for _, c := range tCase.clients {
				client := data.Client{
					Name: c,
					Type: data.Public,
					ID:   uuid.New(),
					Auth: data.Authentication{
						Type:  data.ClientIdAndSecrets,
						Value: uuid.New().String(),
					},
				}
				realm.Clients = append([]data.Client{client})
			}

			for _, u := range tCase.users {
				userJson := sf.Format(`{"info":{"preferred_username":"{0}"}}`, u)
				var rawUser interface{}
				err := json.Unmarshal([]byte(userJson), &rawUser)
				assert.NoError(t, err)
				realm.Users = append([]interface{}{rawUser})
			}

			err := manager.CreateRealm(realm)
			assert.NoError(t, err)
			r, err := manager.GetRealm(realm.Name)
			assert.NoError(t, err)
			checkRealm(t, &realm, r)
			assert.Equal(t, len(tCase.clients), len(r.Clients))
			users, err := manager.GetUsers(realm.Name)
			assert.NoError(t, err)
			assert.Equal(t, len(realm.Users), len(users))
			expectedUsers := make([]data.User, len(realm.Users))
			if len(realm.Users) > 0 {

				for i, _ := range realm.Users {
					expectedUsers[i] = data.CreateUser(realm.Users[i])
				}
			}
			checkUsers(t, &expectedUsers, &users)
			err = manager.DeleteRealm(realm.Name)
			assert.NoError(t, err)
		})
	}
}

func TestCreateRealmFailsDuplicateRealm(t *testing.T) {
	manager := createTestRedisDataManager()
	realm := data.Realm{
		Name:                   "realm_for_duplicate_check",
		TokenExpiration:        3600,
		RefreshTokenExpiration: 1800,
	}

	err := manager.CreateRealm(realm)
	assert.NoError(t, err)

	err = manager.CreateRealm(realm)
	assert.Error(t, err)
	assert.True(t, errors.As(err, &appErrs.ErrExists))

	err = manager.DeleteRealm(realm.Name)
	assert.NoError(t, err)
}

func TestUpdateRealmSuccessfully(t *testing.T) {
	// 1. Create Realm
	manager := createTestRedisDataManager()
	realm := data.Realm{
		Name:                   sf.Format("realm_4_update_check_{0}", uuid.New().String()),
		TokenExpiration:        3600,
		RefreshTokenExpiration: 1800,
	}
	err := manager.CreateRealm(realm)
	assert.NoError(t, err)
	r, err := manager.GetRealm(realm.Name)
	assert.NoError(t, err)
	checkRealm(t, &realm, r)
	prevRealmName := realm.Name
	// 2. Update Realm
	realm.Name = sf.Format("realm_4_update_check_{0}_new_realm_name", uuid.New().String())
	realm.TokenExpiration = 5400
	client := data.Client{
		Name: "app_4_update_realm",
		Type: data.Public,
		ID:   uuid.New(),
		Auth: data.Authentication{
			Type:  data.ClientIdAndSecrets,
			Value: uuid.New().String(),
		},
	}
	realm.Clients = append([]data.Client{client})

	userJson := sf.Format(`{"info":{"preferred_username":"{0}"}}`, "new_app_user")
	var rawUser interface{}
	err = json.Unmarshal([]byte(userJson), &rawUser)
	assert.NoError(t, err)
	realm.Users = append([]interface{}{rawUser})

	err = manager.UpdateRealm(prevRealmName, realm)
	assert.NoError(t, err)
	r, err = manager.GetRealm(realm.Name)
	assert.NoError(t, err)
	checkRealm(t, &realm, r)
	// 3. Delete Realm
	err = manager.DeleteRealm(realm.Name)
	assert.NoError(t, err)
}

func TestUpdateRealmFailsNonExistingRealm(t *testing.T) {
	manager := createTestRedisDataManager()
	realm := data.Realm{
		Name: "super-duper-realm",
	}
	nonExistingRealm := sf.Format("non_existing_{0}", uuid.New().String())
	err := manager.UpdateRealm(nonExistingRealm, realm)
	assert.Error(t, err)
	assert.True(t, errors.As(err, &appErrs.EmptyNotFoundErr))
}

func TestDeleteRealmFailsNonExistingRealm(t *testing.T) {
	manager := createTestRedisDataManager()

	nonExistingRealm := sf.Format("non_existing_{0}", uuid.New().String())
	err := manager.DeleteRealm(nonExistingRealm)
	assert.Error(t, err)
	assert.True(t, errors.As(err, &appErrs.EmptyNotFoundErr))
}

func TestGetRealmFailsNonExistingRealm(t *testing.T) {
	manager := createTestRedisDataManager()

	nonExistingRealm := sf.Format("non_existing_{0}", uuid.New().String())
	_, err := manager.GetRealm(nonExistingRealm)
	assert.Error(t, err)
	assert.True(t, errors.As(err, &appErrs.EmptyNotFoundErr))
}

func TestCreateClientSuccessfully(t *testing.T) {
	testCases := []struct {
		name       string
		clientName string
		clientType data.ClientType
	}{
		{name: "create_public_client", clientName: "sample_pub_client", clientType: data.Public},
		{name: "create_conf_client", clientName: "sample_conf_client", clientType: data.Confidential},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()
			manager := createTestRedisDataManager()
			realm := data.Realm{
				Name:                   "sample_realm_4_create_client_tests",
				TokenExpiration:        3600,
				RefreshTokenExpiration: 1800,
			}
			err := manager.CreateRealm(realm)
			assert.NoError(t, err)
			_, err = manager.GetRealm(realm.Name)
			assert.NoError(t, err)

			client := data.Client{
				Name: "app_4_update_realm",
				Type: tCase.clientType,
				ID:   uuid.New(),
			}
			if client.Type == data.Confidential {
				client.Auth = data.Authentication{
					Type:  data.ClientIdAndSecrets,
					Value: uuid.New().String(),
				}
			}

			err = manager.CreateClient(realm.Name, client)
			assert.NoError(t, err)

			c, err := manager.GetClient(realm.Name, client.Name)
			assert.NoError(t, err)
			checkClient(t, &client, c)

			err = manager.DeleteRealm(realm.Name)
			assert.NoError(t, err)
		})
	}
}

func TestCreateClientFailsDuplicateClient(t *testing.T) {
	manager := createTestRedisDataManager()
	realm := data.Realm{
		Name:                   "sample_realm_4_check_client_duplicate",
		TokenExpiration:        3600,
		RefreshTokenExpiration: 1800,
	}
	err := manager.CreateRealm(realm)
	assert.NoError(t, err)

	client := data.Client{
		Name: "app_4_check_duplicate_client_create",
		Type: data.Public,
		ID:   uuid.New(),
	}

	err = manager.CreateClient(realm.Name, client)
	assert.NoError(t, err)

	err = manager.CreateClient(realm.Name, client)
	assert.Error(t, err)
	assert.True(t, errors.As(err, &appErrs.ErrExists))

	err = manager.DeleteRealm(realm.Name)
	assert.NoError(t, err)
}

func TestUpdateClientSuccessfully(t *testing.T) {
	manager := createTestRedisDataManager()
	realm := data.Realm{
		Name:                   "sample_realm_4_test_client_update",
		TokenExpiration:        3600,
		RefreshTokenExpiration: 1800,
	}
	err := manager.CreateRealm(realm)
	assert.NoError(t, err)

	client := data.Client{
		Name: "app_4_check_client_update",
		Type: data.Public,
		ID:   uuid.New(),
	}

	err = manager.CreateClient(realm.Name, client)
	assert.NoError(t, err)

	c, err := manager.GetClient(realm.Name, client.Name)
	assert.NoError(t, err)
	checkClient(t, &client, c)

	client.Auth = data.Authentication{
		Type:  data.ClientIdAndSecrets,
		Value: uuid.New().String(),
	}

	err = manager.UpdateClient(realm.Name, client.Name, client)
	assert.NoError(t, err)
	c, err = manager.GetClient(realm.Name, client.Name)
	assert.NoError(t, err)
	checkClient(t, &client, c)

	err = manager.DeleteRealm(realm.Name)
	assert.NoError(t, err)
}

func TestUpdateClientFailsNonExistingClient(t *testing.T) {
	manager := createTestRedisDataManager()
	realm := data.Realm{
		Name:                   "sample_realm_4_test_non_existing_client_update",
		TokenExpiration:        3600,
		RefreshTokenExpiration: 1800,
	}
	err := manager.CreateRealm(realm)
	assert.NoError(t, err)

	nonExistingClient := sf.Format("non-existing-client_{0}", uuid.New().String())
	client := data.Client{
		Name: "Surprise Motherfucker",
	}

	err = manager.UpdateClient(realm.Name, nonExistingClient, client)
	assert.Error(t, err)
	assert.True(t, errors.As(err, &appErrs.EmptyNotFoundErr))

	err = manager.DeleteRealm(realm.Name)
	assert.NoError(t, err)
}

func TestDeleteClientFailsNonExistingClient(t *testing.T) {
	manager := createTestRedisDataManager()
	realm := data.Realm{
		Name:                   "sample_realm_4_test_non_existing_client_update",
		TokenExpiration:        3600,
		RefreshTokenExpiration: 1800,
	}
	err := manager.CreateRealm(realm)
	assert.NoError(t, err)

	nonExistingClient := sf.Format("non-existing-client_{0}", uuid.New().String())

	err = manager.DeleteClient(realm.Name, nonExistingClient)
	assert.Error(t, err)
	assert.True(t, errors.As(err, &appErrs.EmptyNotFoundErr))

	err = manager.DeleteRealm(realm.Name)
	assert.NoError(t, err)
}

func TestGetClientFailsNonExistingClient(t *testing.T) {
	manager := createTestRedisDataManager()
	realm := data.Realm{
		Name:                   "sample_realm_4_test_non_existing_client_get",
		TokenExpiration:        3600,
		RefreshTokenExpiration: 1800,
	}
	err := manager.CreateRealm(realm)
	assert.NoError(t, err)

	nonExistingClient := sf.Format("non-existing-client_{0}", uuid.New().String())

	_, err = manager.GetClient(realm.Name, nonExistingClient)
	assert.Error(t, err)
	assert.True(t, errors.As(err, &appErrs.EmptyNotFoundErr))

	err = manager.DeleteRealm(realm.Name)
	assert.NoError(t, err)
}

func TestCreateUserSuccessfully(t *testing.T) {
	testCases := []struct {
		name              string
		realmNameTemplate string
		userName          string
	}{
		{name: "create_min_user", realmNameTemplate: "app_realm_{0}", userName: "app_user"},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()
			manager := createTestRedisDataManager()
			// here we are going to create user separately from Realm via manager.CreateUser
			realm := data.Realm{
				Name:                   sf.Format(tCase.realmNameTemplate, uuid.New().String()),
				TokenExpiration:        3600,
				RefreshTokenExpiration: 1800,
			}
			err := manager.CreateRealm(realm)
			assert.NoError(t, err)
			r, err := manager.GetRealm(realm.Name)
			assert.NoError(t, err)
			// TODO(UMV): IMPL FULL COMPARISON, HERE WE MAKE VERY FORMAL COMPARISON
			checkRealm(t, &realm, r)

			jsonTemplate := `{"info":{"name":"{0}", "preferred_username": "{1}"}, "credentials":{"password": "123"}}`
			jsonStr := sf.Format(jsonTemplate, tCase.userName, tCase.userName)
			var rawUser interface{}
			err = json.Unmarshal([]byte(jsonStr), &rawUser)
			assert.NoError(t, err)
			user := data.CreateUser(rawUser)
			err = manager.CreateUser(realm.Name, user)
			assert.NoError(t, err)
			storedUser, err := manager.GetUser(realm.Name, tCase.userName)
			assert.NoError(t, err)
			assert.Equal(t, tCase.userName, storedUser.GetUsername())
			err = manager.DeleteRealm(realm.Name)
			assert.NoError(t, err)
		})
	}
}

/*func TestUpdateUserSuccessfully(t *testing.T) {

}*/

func TestChangeUserPasswordSuccessfully(t *testing.T) {
	manager := createTestRedisDataManager()
	// 1. Create Realm+Client+User
	realm := data.Realm{
		Name:                   sf.Format("app_4_user_pwd_change_check_{0}", uuid.New().String()),
		TokenExpiration:        3600,
		RefreshTokenExpiration: 1800,
	}

	client := data.Client{
		Name: "app_client_4_check_pwd_change",
		Type: data.Public,
		ID:   uuid.New(),
		Auth: data.Authentication{
			Type:  data.ClientIdAndSecrets,
			Value: uuid.New().String(),
		},
	}
	realm.Clients = append([]data.Client{client})

	userName := "new_app_user"
	userJson := sf.Format(`{"info":{"preferred_username":"{0}"}, "credentials":{"password": "123"}}`, userName)
	var rawUser interface{}
	err := json.Unmarshal([]byte(userJson), &rawUser)
	assert.NoError(t, err)
	realm.Users = append([]interface{}{rawUser})

	err = manager.CreateRealm(realm)
	assert.NoError(t, err)
	_, err = manager.GetRealm(realm.Name)
	assert.NoError(t, err)

	// 2. Reset Password and check ...
	err = manager.SetPassword(realm.Name, userName, "123_ololo_321")
	assert.NoError(t, err)

	err = manager.DeleteRealm(realm.Name)
	assert.NoError(t, err)
}

func createTestRedisDataManager() *RedisDataManager {
	rndNamespace := sf.Format("ferrum_test_{0}", uuid.New().String())
	dataSourceCfg := config.DataSourceConfig{
		Type:   config.REDIS,
		Source: testRedisSource,
		Options: map[config.DataSourceConnOption]string{
			config.Namespace: rndNamespace,
			config.DbNumber:  "0",
		},
		Credentials: &config.CredentialsConfig{
			Username: testUser,
			Password: testUserPassword,
		},
	}

	loggerCfg := config.LoggingConfig{}

	logger := logging.CreateLogger(&loggerCfg)
	manager, _ := CreateRedisDataManager(&dataSourceCfg, logger)
	return manager
}

func checkRealm(t *testing.T, expected *data.Realm, actual *data.Realm) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.TokenExpiration, actual.TokenExpiration)
	assert.Equal(t, expected.RefreshTokenExpiration, actual.RefreshTokenExpiration)
}

func checkClient(t *testing.T, expected *data.Client, actual *data.Client) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Auth.Type, actual.Auth.Type)
	assert.Equal(t, expected.Auth.Value, actual.Auth.Value)
}

func checkUsers(t *testing.T, expected *[]data.User, actual *[]data.User) {
	assert.Equal(t, len(*expected), len(*actual))
	for _, e := range *expected {
		// check and find actual ....
		found := false
		for _, a := range *actual {
			if e.GetId() == a.GetId() {
				checkUser(t, &e, &a)
				found = true
				break
			}
		}
		assert.True(t, found)
	}
}

func checkUser(t *testing.T, expected *data.User, actual *data.User) {
	assert.Equal(t, (*expected).GetId(), (*actual).GetId())
	assert.Equal(t, (*expected).GetUsername(), (*actual).GetUsername())
	assert.Equal(t, (*expected).GetPassword(), (*actual).GetPassword())
}
