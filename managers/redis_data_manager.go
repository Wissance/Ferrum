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
	realmCollection         = "realms"
	clientsCollection       = "clients"
	usersCollection         = "users"
	userKeyTemplate         = "fe_user_{0}"
	realmKeyTemplate        = "fe_realm_{0}"
	realmClientsKeyTemplate = "fe_realm_{0}_clients"
	clientKeyTemplate       = "fe_client_{0}"
	realmUsersKeyTemplate   = "fe_realm_{0}_users"
)

// RedisDataManager is a redis client
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
	// todo(umv) think about preload ???
	dc := DataContext(mn)
	return dc
}

func (mn *RedisDataManager) GetRealm(realmName string) *data.Realm {
	realmKey := sf.Format(realmKeyTemplate, realmName)
	realm := getObjectFromRedis[data.Realm](mn.redisClient, mn.ctx, mn.logger, "realm", realmKey)
	return realm
	/*redisCmd := mn.redisClient.Get(mn.ctx, realmKey)
	if redisCmd.Err() != nil {
		mn.logger.Warn(sf.Format("An error occurred during fetching realm: \"{0}\" from Redis server", realmName))
		return nil
	}

	var realm data.Realm
	realmJson := []byte(redisCmd.Val())
	err := json.Unmarshal(realmJson, &realm)
	if err != nil {
		mn.logger.Error(sf.Format("An error occurred during realm : \"{0}\" unmarshall", realmName))
	}
	return &realm*/
}

func (mn *RedisDataManager) GetClient(realm *data.Realm, name string) *data.Client {
	clientKey := sf.Format(clientKeyTemplate, name)
	client := getObjectFromRedis[data.Client](mn.redisClient, mn.ctx, mn.logger, "client", clientKey)
	// todo (UMV): check that client is from Realm, if not warn
	return client
}

func (mn *RedisDataManager) GetUser(realm *data.Realm, userName string) *data.User {
	return nil
}

func (mn *RedisDataManager) GetUserById(realm *data.Realm, userId uuid.UUID) *data.User {
	return nil
}

func (mn *RedisDataManager) GetRealmUsers(realmName string) *[]data.User {
	return nil
}

func getObjectFromRedis[T any](redisClient *redis.Client, ctx context.Context, logger *logging.AppLogger,
	objName string, objKey string) *T {
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
