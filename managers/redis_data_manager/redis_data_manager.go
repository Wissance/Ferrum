package redis_data_manager

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/Ferrum/managers/errors_managers"
	sf "github.com/wissance/stringFormatter"
)

// This set of const of a templates to all data storing in Redis it contains prefix - a namespace {0}
const (
	userKeyTemplate         = "{0}.{1}_user_{2}"
	realmKeyTemplate        = "{0}.realm_{1}"
	realmClientsKeyTemplate = "{0}.realm_{1}_clients"
	clientKeyTemplate       = "{0}.{1}_client_{2}"
	realmUsersKeyTemplate   = "{0}.realm_{1}_users"
	// realmUsersFullDataKeyTemplate = "{0}.realm_{1}_users_full_data"
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

// CreateRedisDataManager is factory function for instance of RedisDataManager creation and return as interface DataContext
/* Simply creates instance of RedisDataManager and initializes redis client, this function requires config.Namespace to be set up in configs, otherwise
 * defaultNamespace is using
 * Parameters:
 *     - dataSourceCfg contains Redis specific settings in Options map (see allowed keys of map in config.DataSourceConnOption)
 *     - logger - initialized logger instance
 */
func CreateRedisDataManager(dataSourceCfg *config.DataSourceConfig, logger *logging.AppLogger) (*RedisDataManager, error) {
	// todo(UMV): todo provide an error handling
	opts := buildRedisConfig(dataSourceCfg, logger)
	rClient := redis.NewClient(opts)
	namespace, ok := dataSourceCfg.Options[config.Namespace]
	if !ok || len(namespace) == 0 {
		namespace = defaultNamespace
	}
	mn := &RedisDataManager{
		logger: logger, redisOption: opts, redisClient: rClient, ctx: context.Background(),
		namespace: namespace,
	}
	return mn, nil
}

// buildRedisConfig builds redis.Options from map of values by known in config package set of keys
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

// getObjectFromRedis is a method that DOESN'T work with List type object, only a String object type.
func getObjectFromRedis[T any](redisClient *redis.Client, ctx context.Context, logger *logging.AppLogger,
	objName objectType, objKey string,
) (*T, error) {
	redisCmd := redisClient.Get(ctx, objKey)
	if redisCmd.Err() != nil {
		logger.Warn(sf.Format("An error occurred during fetching {0}: \"{1}\" from Redis server", objName, objKey))
		if redisCmd.Err() == redis.Nil {
			return nil, errors_managers.ErrNotFound
		}
		return nil, redisCmd.Err()
	}

	var obj T
	jsonBin := []byte(redisCmd.Val())
	err := json.Unmarshal(jsonBin, &obj)
	if err != nil {
		logger.Error(sf.Format("An error occurred during {0} : \"{1}\" unmarshall", objName, objKey))
		return nil, fmt.Errorf("json.Unmarshal failed: %w", err)
	}
	return &obj, nil
}

// getObjectFromRedis is a method that DOESN'T work with List type object, only a String object type
// Does not return an error if the object is not found
func getMultipleObjectFromRedis[T any](redisClient *redis.Client, ctx context.Context, logger *logging.AppLogger,
	objName objectType, objKey []string,
) ([]T, error) {
	redisCmd := redisClient.MGet(ctx, objKey...)
	if redisCmd.Err() != nil {
		// todo(UMV): print when this will be done https://github.com/Wissance/stringFormatter/issues/14
		logger.Warn(sf.Format("An error occurred during fetching {0}: from Redis server", objName))
		return nil, redisCmd.Err()
	}

	raw := redisCmd.Val()
	result := make([]T, len(raw))
	var unMarshalledRaw interface{}
	for i, v := range raw {
		err := json.Unmarshal([]byte(v.(string)), &unMarshalledRaw)
		if err != nil {
			logger.Error(sf.Format("An error occurred during {0} : \"{1}\" unmarshall", objName, objKey))
			return nil, fmt.Errorf("json.Unmarshal failed: %w", err)
		}
		result[i] = unMarshalledRaw.(T)
	}
	return result, nil
}

// this functions gets object that stored as a LIST Object type
func getObjectsListFromRedis[T any](redisClient *redis.Client, ctx context.Context, logger *logging.AppLogger,
	objName objectType, objKey string,
) ([]T, error) {
	redisCmd := redisClient.LRange(ctx, objKey, 0, -1)
	if redisCmd.Err() != nil {
		logger.Warn(sf.Format("An error occurred during fetching {0}: \"{1}\" from Redis server", objName, objKey))
		return nil, redisCmd.Err()
	}

	// var obj T
	items := redisCmd.Val()
	if len(items) == 0 {
		return nil, errors_managers.ErrNotFound
	}
	var result []T
	var portion []T
	for _, rawVal := range items {
		jsonBin := []byte(rawVal)
		err := json.Unmarshal(jsonBin, &portion) // already contains all SLICE in one object
		if err != nil {
			logger.Error(sf.Format("An error occurred during {0} : \"{1}\" unmarshall", objName, objKey))
			return nil, fmt.Errorf("json.Unmarshal failed: %w", err)
		}
		result = append(result, portion...)
	}
	return result, nil
}

// If such a key exists, the value will be overwritten without error
func setString(redisClient *redis.Client, ctx context.Context, logger *logging.AppLogger,
	objName objectType, objKey string, objValue string,
) error {
	statusCmd := redisClient.Set(ctx, objKey, objValue, 0)
	if statusCmd.Err() != nil {
		logger.Warn(sf.Format("An error occurred during set {0}: \"{1}\": \"{2}\" from Redis server", objName, objKey, objValue))
		return statusCmd.Err()
	}
	return nil
}

func delKey(redisClient *redis.Client, ctx context.Context, logger *logging.AppLogger, objName objectType, objKey string) error {
	redisIntCmd := redisClient.Del(ctx, objKey)
	res := redisIntCmd.Val()
	if res == 0 {
		// TODO(SIA) add log
		return errors_managers.ErrNotExists
	}
	return nil
}

// TODO(SIA) add function
// func rPush(redisClient *redis.Client, ctx context.Context, logger *logging.AppLogger,
// 	objName objectType, objKey string, objValue string,
// ) error {
// 	statusCmd := redisClient.Set(ctx, objKey, objValue, 0)
// 	if statusCmd.Err() != nil {
// 		logger.Warn(sf.Format("An error occurred during set {0}: \"{1}\": \"{2}\" from Redis server", objName, objKey, objValue))
// 		return statusCmd.Err()
// 	}
// 	return nil
// }

// TODO(SIA) add function
// func isKeyExists(redisClient *redis.Client, ctx context.Context, logger *logging.AppLogger, key string) bool {
// 	redisIntCmd := redisClient.Exists(ctx, key)
// 	if redisIntCmd.Err() != nil {
// 		logger.Warn(sf.Format("An error occurred during fetching {0}: \"{1}\" from Redis server", objName, objKey))
// 		return nil, redisCmd.Err()
// 	}
// 	return nil
// }
