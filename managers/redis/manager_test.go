package redis

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
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
			// TODO(UMV): IMPL FULL COMPARISON, HERE WE MAKE VERY FORMAL COMPARISON
			assert.Equal(t, realm.Name, r.Name)
			assert.Equal(t, len(tCase.clients), len(r.Clients))
			users, err := manager.GetUsers(realm.Name)
			assert.NoError(t, err)
			assert.Equal(t, len(realm.Users), len(users))
			err = manager.DeleteRealm(realm.Name)
			assert.NoError(t, err)
		})
	}
}

func TestUpdateRealmSuccessfully(t *testing.T) {

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
			assert.Equal(t, realm.Name, r.Name)

			jsonTemplate := `{"info":{"name":"{0}", "preferred_username": "{1}"}}`
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

func TestUpdateUserSuccessfully(t *testing.T) {

}

func TestResetUserPasswordSuccessfully(t *testing.T) {

}

func TestChangeUserPasswordSuccessfully(t *testing.T) {

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
