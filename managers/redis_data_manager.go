package managers

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	sf "github.com/wissance/stringFormatter"
	"strconv"
)

// This set of const of a templates to all data storing in Redis it contains prefix - a namespace {0}
const (
	userKeyTemplate               = "{0}.user_{1}"
	realmKeyTemplate              = "{0}.realm_{1}"
	realmClientsKeyTemplate       = "{0}.realm_{1}_clients"
	clientKeyTemplate             = "{0}.client_{1}"
	realmUsersKeyTemplate         = "{0}.realm_{1}_users"
	realmUsersFullDataKeyTemplate = "{0}.realm_{1}_users_full_data"
)

type objectType string

const (
	Realm        objectType = "realm"
	RealmClients            = "realm clients"
	RealmUsers              = "realm users"
	Client                  = "client"
	User                    = "user"
)

const defaultNamespace = "fe"

// RedisDataManager is a redis client
/*
 * Redis Data Manager is a service class for managing authorization server data in Redis
 * There are following store Rules:
 * 1. Realms (data.Realm) in Redis storing separately from Clients && Users, every Realm stores in Redis by key forming from template && Realm name
 *    i.e. if we have Realm with name "wissance" it could be accessed by key fe_realm_wissance (realmKeyTemplate)
 * 2. Realm Clients ([]data.ExtendedIdentifier) storing in Redis by key forming from template, Realm with name wissance has array of clients id by key
 *    fe_realm_wissance_clients (realmClientsKeyTemplate)
 * 3. Every Client (data.Client) stores separately by key forming from client id (different realms could have clients with same name but in different realm,
 *    Client Name is unique only in Realm) and template clientKeyTemplate, therefore realm with pair (ID: 6e09faca-1004-11ee-be56-0242ac120002 Name: homeApp)
 *    could be received by key - fe_client_6e09faca-1004-11ee-be56-0242ac120002
 * 4. Every User in Redis storing by it own key forming by userId + template (userKeyTemplate) -> i.e. user with id 6dee45ee-1056-11ee-be56-0242ac120002 stored
 *    by key fe_user_6dee45ee-1056-11ee-be56-0242ac120002
 * 5. Client to Realm and User to Realm relation stored by separate keys forming using template and realm name, these relations stores array of data.ExtendedIdentifier
 *    that wires together Realm Name with User.ID and User.Name.
 *    IMPORTANT NOTES:
 *    1. When save Client or User don't forget to save relations in Redis too (see templates realmClientsKeyTemplate && realmUsersKeyTemplate)
 *    2. When add/modify new or existing user don't forget to update realmUsersFullDataKeyTemplate maybe this collection will be removed in future but currently
 *       we have it.
 */
type RedisDataManager struct {
	namespace   string
	redisOption *redis.Options
	redisClient *redis.Client
	logger      *logging.AppLogger
	ctx         context.Context
}

func CreateRedisDataManager(dataSourceCfg *config.DataSourceConfig, logger *logging.AppLogger) (DataContext, error) {
	// todo(UMV): todo provide an error handling
	opts := buildRedisConfig(dataSourceCfg, logger)
	rClient := redis.NewClient(opts)
	namespace, ok := dataSourceCfg.Options[config.Namespace]
	if !ok || len(namespace) == 0 {
		namespace = defaultNamespace
	}
	mn := &RedisDataManager{logger: logger, redisOption: opts, redisClient: rClient, ctx: context.Background(),
		namespace: namespace}
	dc := DataContext(mn)
	return dc, nil
}

func (mn *RedisDataManager) GetRealm(realmName string) *data.Realm {
	realmKey := sf.Format(realmKeyTemplate, mn.namespace, realmName)
	realm := getObjectFromRedis[data.Realm](mn.redisClient, mn.ctx, mn.logger, Realm, realmKey)
	// should get realms too
	// if realms were stored without clients (we expected so), get clients related to realm and assign here
	if len(realm.Clients) == 0 {
		realm.Clients = mn.GetRealmClients(realmName)
	}
	return realm
}

func (mn *RedisDataManager) GetClient(realm *data.Realm, name string) *data.Client {
	realmClientsKey := sf.Format(realmClientsKeyTemplate, mn.namespace, realm.Name)
	// realm_%name%_clients contains array with configured clients ID (data.ExtendedIdentifier) for that realm
	realmClients := getObjectFromRedis[[]data.ExtendedIdentifier](mn.redisClient, mn.ctx, mn.logger, RealmClients, realmClientsKey)
	if realmClients == nil {
		mn.logger.Error(sf.Format("There are no clients for realm: \"{0} \" in Redis, BAD data config", realm.Name))
		return nil
	}
	realmHasClient := false
	var clientId data.ExtendedIdentifier
	for _, rc := range *realmClients {
		if rc.Name == name {
			realmHasClient = true
			clientId = rc
			break
		}
	}
	if !realmHasClient {
		mn.logger.Debug(sf.Format("Realm: \"{0}\" doesn't have client : \"{1}\" in Redis", realm.Name, name))
		return nil
	}
	clientKey := sf.Format(clientKeyTemplate, mn.namespace, clientId.ID)
	client := getObjectFromRedis[data.Client](mn.redisClient, mn.ctx, mn.logger, Client, clientKey)
	return client
}

func (mn *RedisDataManager) GetUser(realm *data.Realm, userName string) *data.User {
	userRealmsKey := sf.Format(realmUsersKeyTemplate, mn.namespace, realm.Name)
	realmUsers := getObjectsListFromRedis[data.ExtendedIdentifier](mn.redisClient, mn.ctx, mn.logger, RealmUsers, userRealmsKey)
	if realmUsers == nil {
		mn.logger.Error(sf.Format("There are no user with name :\"{0}\" in realm: \"{1} \" in Redis, BAD data config", userName, realm.Name))
		return nil
	}

	var extendedUserId data.ExtendedIdentifier
	userFound := false
	for _, rc := range *realmUsers {
		if rc.Name == userName {
			userFound = true
			extendedUserId = rc
			break
		}
	}

	if !userFound {
		mn.logger.Debug(sf.Format("User with name: \"{0}\" was not found for realm: \"{1}\"", userName, realm.Name))
		return nil
	}

	userKey := sf.Format(userKeyTemplate, mn.namespace, extendedUserId.Name)
	rawUser := getObjectFromRedis[interface{}](mn.redisClient, mn.ctx, mn.logger, User, userKey)
	user := data.CreateUser(*rawUser)
	return &user
}

func (mn *RedisDataManager) GetUserById(realm *data.Realm, userId uuid.UUID) *data.User {
	// userKey := sf.Format(userKeyTemplate, mn.namespace, userId)
	var rawUser data.User
	users := mn.GetRealmUsers(realm.Name)
	for _, u := range *users {
		checkingUserId := u.GetId()
		if checkingUserId == userId {
			rawUser = u
			break
		}
	}
	// we can't get user such way
	//rawUser := getObjectFromRedis[interface{}](mn.redisClient, mn.ctx, mn.logger, User, userKey)
	user := data.CreateUser(rawUser)
	userRealmsKey := sf.Format(realmUsersKeyTemplate, mn.namespace, realm.Name)
	realmUsers := getObjectFromRedis[[]data.ExtendedIdentifier](mn.redisClient, mn.ctx, mn.logger, RealmUsers, userRealmsKey)
	if realmUsers == nil {
		mn.logger.Error(sf.Format("There are no user with ID:\"{0}\" in realm: \"{1} \" in Redis, BAD data config", userId.String(), realm.Name))
		return nil
	}

	for _, rc := range *realmUsers {
		if rc.ID == userId {
			return &user
		}
	}

	mn.logger.Debug(sf.Format("User with id: \"{0}\" wasn't found in Realm: {1}", userId, realm.Name))
	return nil
}

func (mn *RedisDataManager) GetRealmUsers(realmName string) *[]data.User {
	// TODO(UMV): possibly we should not use this method ??? what if we have 1M+ users .... ? think maybe it should be somehow optimized ...
	userRealmsKey := sf.Format(realmUsersKeyTemplate, mn.namespace, realmName)

	realmUsers := getObjectsListFromRedis[data.ExtendedIdentifier](mn.redisClient, mn.ctx, mn.logger, RealmUsers, userRealmsKey)
	if realmUsers == nil {
		mn.logger.Error(sf.Format("There are no users in realm: \"{0} \" in Redis, BAD data config", realmName))
		return nil
	}

	userRedisKeys := make([]string, len(*realmUsers))
	for i, ru := range *realmUsers {
		userRedisKeys[i] = sf.Format(userKeyTemplate, mn.namespace, ru.Name)
	}

	// userFullDataRealmsKey := sf.Format(realmUsersFullDataKeyTemplate, mn.namespace, realmName)
	// this is wrong, we can't get rawUsers such way ...
	realmUsersData := getMultipleObjectFromRedis[interface{}](mn.redisClient, mn.ctx, mn.logger, RealmUsers, userRedisKeys)
	//getObjectsListFromRedis[interface{}](mn.redisClient, mn.ctx, mn.logger, RealmUsers, userFullDataRealmsKey)

	if realmUsersData != nil {
		userData := make([]data.User, len(*realmUsersData))
		for i, u := range *realmUsersData {
			userData[i] = data.CreateUser(u)
		}
		return &userData
	}
	return nil
}

func (mn *RedisDataManager) GetRealmClients(realmName string) []data.Client {
	realmClientsKey := sf.Format(realmClientsKeyTemplate, mn.namespace, realmName)
	realmClients := getObjectsListFromRedis[data.ExtendedIdentifier](mn.redisClient, mn.ctx, mn.logger, RealmClients, realmClientsKey)
	if realmClients == nil {
		mn.logger.Error(sf.Format("There are no clients for realm: \"{0} \" in Redis, BAD data config", realmName))
		return nil
	}
	clients := make([]data.Client, len(*realmClients))
	for i, rc := range *realmClients {
		// todo(UMV) get all them at once
		clientKey := sf.Format(clientKeyTemplate, mn.namespace, rc.Name)
		client := getObjectFromRedis[data.Client](mn.redisClient, mn.ctx, mn.logger, Client, clientKey)
		clients[i] = *client
	}

	return clients
}

// getObjectFromRedis is a method that DOESN'T work with List type object, only a String object type
func getObjectFromRedis[T any](redisClient *redis.Client, ctx context.Context, logger *logging.AppLogger,
	objName objectType, objKey string) *T {
	redisCmd := redisClient.Get(ctx, objKey)
	if redisCmd.Err() != nil {
		logger.Warn(sf.Format("An error occurred during fetching {0}: \"{1}\" from Redis server", objName, objKey))
		return nil
	}

	var obj T
	jsonBin := []byte(redisCmd.Val())
	err := json.Unmarshal(jsonBin, &obj)
	if err != nil {
		logger.Error(sf.Format("An error occurred during {0} : \"{1}\" unmarshall", objName, objKey))
	}
	return &obj
}

// getObjectFromRedis is a method that DOESN'T work with List type object, only a String object type
func getMultipleObjectFromRedis[T any](redisClient *redis.Client, ctx context.Context, logger *logging.AppLogger,
	objName objectType, objKey []string) *[]T {
	redisCmd := redisClient.MGet(ctx, objKey...)
	if redisCmd.Err() != nil {
		// todo(UMV): print when this will be done https://github.com/Wissance/stringFormatter/issues/14
		logger.Warn(sf.Format("An error occurred during fetching {0}: from Redis server", objName))
		return nil
	}

	raw := redisCmd.Val()
	result := make([]T, len(raw))
	for i, v := range raw {
		result[i] = v.(T)
	}
	return &result
}

// this functions gets object that stored as a LIST Object type
func getObjectsListFromRedis[T any](redisClient *redis.Client, ctx context.Context, logger *logging.AppLogger,
	objName objectType, objKey string) *[]T {

	redisCmd := redisClient.LRange(ctx, objKey, 0, -1)
	if redisCmd.Err() != nil {
		logger.Warn(sf.Format("An error occurred during fetching {0}: \"{1}\" from Redis server", objName, objKey))
		return nil
	}

	//var obj T
	items := redisCmd.Val()
	var result []T
	var portion []T
	for _, rawVal := range items {
		jsonBin := []byte(rawVal)
		err := json.Unmarshal(jsonBin, &portion) // already contains all SLICE in one object
		if err != nil {
			logger.Error(sf.Format("An error occurred during {0} : \"{1}\" unmarshall", objName, objKey))
		}
		result = append(result, portion...)
	}

	return &result
}

func buildRedisConfig(dataSourceCfd *config.DataSourceConfig, logger *logging.AppLogger) *redis.Options {
	dbNum, err := strconv.Atoi(dataSourceCfd.Options[config.DbNumber])
	if err != nil {
		logger.Error(sf.Format("can't be because we already called Validate(), but in any case: parsing error: {0}", err.Error()))
		return nil
	}
	opts := redis.Options{
		Addr: dataSourceCfd.Source,
		DB:   dbNum,
	}
	// passing credentials if we have it
	if dataSourceCfd.Credentials != nil {
		opts.Username = dataSourceCfd.Credentials.Username
		opts.Password = dataSourceCfd.Credentials.Password
	}
	// passing TLS if we have it
	val, ok := dataSourceCfd.Options[config.UseTls]
	if ok {
		useTls, parseErr := strconv.ParseBool(val)
		if parseErr == nil && useTls {
			opts.TLSConfig = &tls.Config{}
			val, ok = dataSourceCfd.Options[config.InsecureTls]
			if ok {
				inSecTls, parseInSecValErr := strconv.ParseBool(val)
				if parseInSecValErr == nil {
					opts.TLSConfig.InsecureSkipVerify = inSecTls
				}
			}
		}
	}

	return &opts
}
