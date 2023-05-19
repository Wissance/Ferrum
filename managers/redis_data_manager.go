package managers

import (
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/wissance/Ferrum/data"
)

type RedisDataManager struct {
	redisClient *redis.Client
}

func CreateRedisDataManager() DataContext {
	rClient := redis.NewClient(&redis.Options{})
	mn := &RedisDataManager{redisClient: rClient}
	// todo(umv) think about preload ???
	dc := DataContext(mn)
	return dc
}

func (mn *RedisDataManager) GetRealm(realmName string) *data.Realm {
	return nil
}

func (mn *RedisDataManager) GetClient(realm *data.Realm, name string) *data.Client {
	return nil
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
