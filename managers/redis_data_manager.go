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

const (
	userKeyTemplate               = "fe_user_{0}"
	realmKeyTemplate              = "fe_realm_{0}"
	realmClientsKeyTemplate       = "fe_realm_{0}_clients"
	clientKeyTemplate             = "fe_client_{0}"
	realmUsersKeyTemplate         = "fe_realm_{0}_users"
	realmUsersFullDataKeyTemplate = "fe_realm_{0}_users_full_data"
)

type objectType string

const (
	Realm        objectType = "realm"
	RealmClients            = "realm clients"
	RealmUsers              = "realm users"
	Client                  = "client"
	User                    = "user"
)

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
	redisOption *redis.Options
	redisClient *redis.Client
	logger      *logging.AppLogger
	ctx         context.Context
}

func CreateRedisDataManager(dataSourceCfd *config.DataSourceConfig, logger *logging.AppLogger) DataContext {
	opts := buildRedisConfig(dataSourceCfd, logger)
	rClient := redis.NewClient(opts)
	mn := &RedisDataManager{logger: logger, redisOption: opts, redisClient: rClient, ctx: context.Background()}
	dc := DataContext(mn)
	return dc
}

func (mn *RedisDataManager) GetRealm(realmName string) *data.Realm {
	realmKey := sf.Format(realmKeyTemplate, realmName)
	realm := getObjectFromRedis[data.Realm](mn.redisClient, mn.ctx, mn.logger, Realm, realmKey)
	return realm
}

func (mn *RedisDataManager) GetClient(realm *data.Realm, name string) *data.Client {
	realmClientsKey := sf.Format(realmClientsKeyTemplate, realm.Name)
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
	clientKey := sf.Format(clientKeyTemplate, clientId.ID)
	client := getObjectFromRedis[data.Client](mn.redisClient, mn.ctx, mn.logger, Client, clientKey)
	return client
}

func (mn *RedisDataManager) GetUser(realm *data.Realm, userName string) *data.User {
	userRealmsKey := sf.Format(realmUsersKeyTemplate, realm.Name)
	realmUsers := getObjectFromRedis[[]data.ExtendedIdentifier](mn.redisClient, mn.ctx, mn.logger, RealmUsers, userRealmsKey)
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

	userKey := sf.Format(userKeyTemplate, extendedUserId.Name)
	rawUser := getObjectFromRedis[interface{}](mn.redisClient, mn.ctx, mn.logger, User, userKey)
	user := data.CreateUser(rawUser)
	return &user
}

func (mn *RedisDataManager) GetUserById(realm *data.Realm, userId uuid.UUID) *data.User {
	userKey := sf.Format(userKeyTemplate, userId)
	rawUser := getObjectFromRedis[interface{}](mn.redisClient, mn.ctx, mn.logger, User, userKey)
	user := data.CreateUser(rawUser)
	userRealmsKey := sf.Format(realmUsersKeyTemplate, realm.Name)
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
	userRealmsKey := sf.Format(realmUsersKeyTemplate, realmName)
	realmUsers := getObjectFromRedis[[]data.ExtendedIdentifier](mn.redisClient, mn.ctx, mn.logger, RealmUsers, userRealmsKey)
	if realmUsers == nil {
		mn.logger.Error(sf.Format("There are no users in realm: \"{0} \" in Redis, BAD data config", realmName))
		return nil
	}

	userFullDataRealmsKey := sf.Format(realmUsersFullDataKeyTemplate, realmName)
	realmUsersData := getObjectFromRedis[[]interface{}](mn.redisClient, mn.ctx, mn.logger, RealmUsers, userFullDataRealmsKey)

	if realmUsersData != nil {
		userData := make([]data.User, len(*realmUsersData))
		for i, u := range *realmUsersData {
			userData[i] = data.CreateUser(u)
		}
		return &userData
	}
	return nil
}

func getObjectFromRedis[T any](redisClient *redis.Client, ctx context.Context, logger *logging.AppLogger,
	objName objectType, objKey string) *T {
	redisCmd := redisClient.Get(ctx, objKey)
	if redisCmd.Err() != nil {
		logger.Warn(sf.Format("An error occurred during fetching {0}: \"{1}\" from Redis server", objName, objKey))
		return nil
	}

	var obj T
	realmJson := []byte(redisCmd.Val())
	err := json.Unmarshal(realmJson, &obj)
	if err != nil {
		logger.Error(sf.Format("An error occurred during {0} : \"{1}\" unmarshall", objName, objKey))
	}
	return &obj
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
