package redis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/wissance/Ferrum/errors"

	"github.com/redis/go-redis/v9"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/logging"
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
 *    i.e. if we have Realm with name "wissance" it could be accessed by key fe.realm_wissance (realmKeyTemplate)
 * 2. Realm Clients ([]data.ExtendedIdentifier) storing in Redis by key forming from template, Realm with name wissance has array of clients id by key
 *    fe.realm_wissance_clients (realmClientsKeyTemplate)
 * 3. Every Client (data.Client) stores separately by key forming from client name (different realms could have clients with same name but in different realm,
 *    Client Name is unique only in Realm) and template clientKeyTemplate, therefore realm with pair (ID: 6e09faca-1004-11ee-be56-0242ac120002 Name: homeApp)
 *    could be received by key - fe.wissance_client_homeApp
 * 4. Every User in Redis storing by it own key forming by userName + template (userKeyTemplate) -> i.e. user with (ID: 6e09faca-1004-11ee-be56-0242ac120002 Name: homeApp) stored
 *    by key fe.wissance_user_homeApp
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

// IsAvailable methods that checks whether DataContext could be used or not
/* Availability means that redisClient is not NULL and Ready for receive requests
 * Parameters: no
 * Returns true if DataContext is available
 */
func (mn *RedisDataManager) IsAvailable() bool {
	return mn.redisClient != nil && mn.redisClient.Ping(mn.ctx) == nil
}

// CreateRedisDataManager is factory function for instance of RedisDataManager creation
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

// upsertRedisString - inserting or updating a value by key
/* If there is no key, a key-value will be created. If the key is present, the value will be updated
 * Arguments:
 *    - objName - for logger
 *    - objKey - key object in redis
 *	  - objValue - new string value
 * Returns: error
 */
func (mn *RedisDataManager) upsertRedisString(objName objectType, objKey string, objValue string) error {
	statusCmd := mn.redisClient.Set(mn.ctx, objKey, objValue, 0)
	if statusCmd.Err() != nil {
		mn.logger.Warn(sf.Format("An error occurred during Set {0}: \"{1}\" from Redis server", objName, objKey))
		return statusCmd.Err()
	}
	return nil
}

// deleteRedisObject - delete key
/* Returns an error, if 0 items are deleted
 * Arguments:
 *    - objName - for logger
 *    - objKey - key object in redis
 * Returns: error
 */
func (mn *RedisDataManager) deleteRedisObject(objName objectType, objKey string) error {
	redisIntCmd := mn.redisClient.Del(mn.ctx, objKey)
	if redisIntCmd.Err() != nil {
		mn.logger.Warn(sf.Format("An error occurred during Del {0}: \"{1}\" from Redis server", objName, objKey))
		return redisIntCmd.Err()
	}
	res := redisIntCmd.Val()
	if res == 0 {
		mn.logger.Warn(sf.Format("An error occurred during Del, 0 items deleted {0}: \"{1}\" from Redis server", objName, objKey))
		return errors.NewObjectExistsError(string(objName), objKey, "")
	}
	return nil
}

// appendStringToRedisList - inserts a string into the list
/* Adds a row to the list. Internally, it uses RPush
 * Arguments:
 *    - objName - for logger
 *    - objKey - key object in redis
 *    - objValue - new string value
 * Returns: error
 */
func (mn *RedisDataManager) appendStringToRedisList(objName objectType, objKey string, objValue string) error {
	redisIntCmd := mn.redisClient.RPush(mn.ctx, objKey, objValue)
	if redisIntCmd.Err() != nil {
		mn.logger.Warn(sf.Format("An error occurred during RPush {0}: \"{1}\" from Redis server", objName, objKey))
		return redisIntCmd.Err()
	}
	return nil
}

// getSingleRedisObject is a method that DOESN'T work with List type object, only a String object type.
func getSingleRedisObject[T any](redisClient *redis.Client, ctx context.Context, logger *logging.AppLogger,
	objName objectType, objKey string,
) (*T, error) {
	redisCmd := redisClient.Get(ctx, objKey)
	if redisCmd.Err() != nil {
		logger.Warn(sf.Format("An error occurred during fetching {0}: \"{1}\" from Redis server", objName, objKey))
		if redisCmd.Err() == redis.Nil {
			return nil, errors.NewObjectNotFoundError(string(objName), objKey, "")
		}
		return nil, redisCmd.Err()
	}

	var obj T
	jsonBin := []byte(redisCmd.Val())
	err := json.Unmarshal(jsonBin, &obj)
	if err != nil {
		logger.Error(sf.Format("An error occurred during unmarshall {0} : \"{1}\"", objName, objKey))
		return nil, fmt.Errorf("json.Unmarshal failed: %w", err)
	}
	return &obj, nil
}

// getMultipleRedisObjects is a method that DOESN'T work with List type object, only a String object type
// Does not return an error if the object is not found
func getMultipleRedisObjects[T any](redisClient *redis.Client, ctx context.Context, logger *logging.AppLogger,
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
			logger.Error(sf.Format("An error occurred during unmarshall {0} : \"{1}\"", objName, objKey))
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
		logger.Warn(sf.Format("Received zero list items {0}: \"{1}\" from Redis server", objName, objKey))
		return nil, errors.ErrZeroLength
	}
	var result []T
	var portion []T
	for _, rawVal := range items {
		jsonBin := []byte(rawVal)
		err := json.Unmarshal(jsonBin, &portion) // already contains all SLICE in one object
		if err != nil {
			logger.Error(sf.Format("An error occurred during unmarshall {0} : \"{1}\", err: ", objName, objKey, err.Error()))
			return nil, err
		}
		result = append(result, portion...)
	}
	return result, nil
}

// TODO(SIA) add function keyExists
// TODO(SIA) Add a function to delete multiple keys at once
