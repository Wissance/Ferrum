package redis_data_manager

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/managers/errors_managers"
	sf "github.com/wissance/stringFormatter"
)

func (mn *RedisDataManager) GetUsersFromRealm(realmName string) ([]data.User, error) {
	// TODO(UMV): possibly we should not use this method ??? what if we have 1M+ users .... ? think maybe it should be somehow optimized ...
	realmUsers, err := mn.getRealmUsers(realmName)
	if err != nil {
		return nil, fmt.Errorf("getRealmUsers failed: %w", err)
	}

	// todo(UMV): probably we should organize batching here if we have many users i.e. 100K+
	userRedisKeys := make([]string, len(realmUsers))
	for i, ru := range realmUsers {
		userRedisKeys[i] = sf.Format(userKeyTemplate, mn.namespace, realmName, ru.Name)
	}

	// userFullDataRealmsKey := sf.Format(realmUsersFullDataKeyTemplate, mn.namespace, realmName)
	// this is wrong, we can't get rawUsers such way ...
	realmUsersData, err := getMultipleObjectFromRedis[interface{}](mn.redisClient, mn.ctx, mn.logger, RealmUsers, userRedisKeys)
	if err != nil {
		return nil, fmt.Errorf("getMultipleObjectFromRedis failed: %w", err)
	}
	// getObjectsListFromRedis[interface{}](mn.redisClient, mn.ctx, mn.logger, RealmUsers, userFullDataRealmsKey)
	if len(realmUsersData) == 0 {
		mn.logger.Error(sf.Format("Redis does not have all users that belong to Realm: \"{0}\"", realmName))
		return nil, fmt.Errorf("getMultipleObjectFromRedis failed: %w", errors_managers.ErrZeroLength)
	}
	if len(realmUsers) != len(realmUsersData) {
		mn.logger.Error(sf.Format("Realm: \"{0}\" has users, that Redis does not have part of it", realmName))
		return nil, errors_managers.ErrNotAll
	}

	userData := make([]data.User, len(realmUsersData))
	for i, u := range realmUsersData {
		userData[i] = data.CreateUser(u)
	}
	return userData, nil
}

func (mn *RedisDataManager) GetUser(realmName string, userName string) (data.User, error) {
	userKey := sf.Format(userKeyTemplate, mn.namespace, realmName, userName)
	rawUser, err := getObjectFromRedis[interface{}](mn.redisClient, mn.ctx, mn.logger, User, userKey)
	if err != nil {
		if errors.Is(err, errors_managers.ErrNotFound) {
			mn.logger.Error(sf.Format("Realm: \"{0}\" does not have Client: \"{1}\"", realmName, userName))
		}
		return nil, fmt.Errorf("getObjectFromRedis failed: %w", err)
	}
	user := data.CreateUser(*rawUser)
	return user, nil
}

func (mn *RedisDataManager) GetUserById(realmName string, userId uuid.UUID) (data.User, error) {
	realmUser, err := mn.getRealmUserById(realmName, userId)
	if err != nil {
		return nil, fmt.Errorf("getRealmUserById failed: %w", err)
	}
	user, err := mn.GetUser(realmName, realmUser.Name)
	if err != nil {
		if errors.Is(err, errors_managers.ErrNotFound) { // TODO(SIA) check
			mn.logger.Error(sf.Format("Realm: \"{0}\" has client: \"{1}\", that Redis does not have", realmName, userId))
		}
		return nil, fmt.Errorf("GetUser failed: %w", err)
	}
	return user, nil
}

// Returns an error if the user exists in redis
func (mn *RedisDataManager) CreateUser(realmName string, userNew any) error {
	// TODO(SIA) транзакции
	// TODO(SIA) возможно нужно проверять, что есть какие-то поля у user
	user := data.CreateUser(userNew)
	userName := user.GetUsername()
	_, err := mn.GetUser(realmName, userName) // TODO(SIA) use function isExists
	if err == nil {
		return errors_managers.ErrExists
	}
	if !errors.Is(err, errors_managers.ErrNotFound) {
		return fmt.Errorf("GetUser failed: %w", err)
	}
	_, err = mn.GetRealm(realmName) // TODO(SIA) use function isExists
	if err != nil {
		return fmt.Errorf("GetRealm failed: %w", err)
	}

	userBytes, err := json.Marshal(userNew)
	if err != nil {
		mn.logger.Error(sf.Format("An error occurred during Client marshal")) // TODO(SIA) ADD NAME
		return fmt.Errorf("json.Marshal failed: %w", err)
	}
	err = mn.createUserRedis(realmName, userName, string(userBytes))
	if err != nil {
		return fmt.Errorf("createClientRedis failed: %w", err)
	}

	if err := mn.addUserToRealm(realmName, user); err != nil {
		return fmt.Errorf("addUserToRealm failed: %w", err)
	}

	return nil
}

func (mn *RedisDataManager) DeleteUser(realmName string, userName string) error {
	if err := mn.deleteUserRedis(realmName, userName); err != nil {
		return fmt.Errorf("deleteUserRedis failed: %w", err)
	}
	if err := mn.deleteUserFromRealm(realmName, userName); err != nil {
		if errors.Is(err, errors_managers.ErrNotFound) || errors.Is(err, errors_managers.ErrZeroLength) {
			return nil
		}
		return fmt.Errorf("deleteUserFromRealm failed: %w", err)
	}
	return nil
}

func (mn *RedisDataManager) UpdateUser(realmName string, userName string, userNew any) error {
	// TODO(SIA) транзакции
	oldUser, err := mn.GetUser(realmName, userName)
	if err != nil {
		return fmt.Errorf("GetUser failed: %w", err)
	}
	oldUserName := oldUser.GetUsername()
	oldUserId := oldUser.GetId()

	user := data.CreateUser(userNew)
	newUserName := user.GetUsername()
	newUserId := user.GetId()

	if newUserId != oldUserId || newUserName != oldUserName {
		if err := mn.DeleteUser(realmName, oldUserName); err != nil {
			return fmt.Errorf("DeleteUser failed: %w", err)
		}
		if err := mn.addUserToRealm(realmName, user); err != nil {
			return fmt.Errorf("addUserToRealm failed: %w", err)
		}
	}

	userBytes, err := json.Marshal(userNew)
	if err != nil {
		mn.logger.Error(sf.Format("An error occurred during Client marshal")) // TODO(SIA) ADD NAME
		return fmt.Errorf("json.Marshal failed: %w", err)
	}
	err = mn.createUserRedis(realmName, newUserName, string(userBytes))
	if err != nil {
		return fmt.Errorf("createUserRedis failed: %w", err)
	}

	return nil
}

func (mn *RedisDataManager) getRealmUsers(realmName string) ([]data.ExtendedIdentifier, error) {
	userRealmsKey := sf.Format(realmUsersKeyTemplate, mn.namespace, realmName)
	realmUsers, err := getObjectsListFromRedis[data.ExtendedIdentifier](mn.redisClient, mn.ctx, mn.logger, RealmUsers, userRealmsKey)
	if err != nil {
		if errors.Is(err, errors_managers.ErrNotFound) {
			mn.logger.Error(sf.Format("There are no users in realm: \"{0}\" in Redis", realmName))
		}
		return nil, fmt.Errorf("getObjectsListFromRedis failed: %w", err)
	}
	return realmUsers, nil
}

func (mn *RedisDataManager) getRealmUser(realmName string, userName string) (*data.ExtendedIdentifier, error) {
	realmUsers, err := mn.getRealmUsers(realmName)
	if err != nil {
		return nil, fmt.Errorf("getRealmUsers failed: %w", err)
	}

	var user data.ExtendedIdentifier
	userFound := false
	for _, rc := range realmUsers {
		if rc.Name == userName {
			userFound = true
			user = rc
			break
		}
	}
	if !userFound {
		mn.logger.Debug(sf.Format("User with name: \"{0}\" was not found for realm: \"{1}\"", userName, realmName))
		return nil, errors_managers.ErrNotFound
	}
	return &user, nil
}

func (mn *RedisDataManager) getRealmUserById(realmName string, userId uuid.UUID) (*data.ExtendedIdentifier, error) {
	realmUsers, err := mn.getRealmUsers(realmName)
	if err != nil {
		return nil, fmt.Errorf("getRealmUsers failed: %w", err)
	}
	var user data.ExtendedIdentifier
	userFound := false
	for _, rc := range realmUsers {
		if rc.ID == userId {
			userFound = true
			user = rc
			break
		}
	}
	if !userFound {
		mn.logger.Debug(sf.Format("User with id: \"{0}\" was not found for realm: \"{1}\"", userId, realmName))
		return nil, errors_managers.ErrNotFound
	}
	return &user, nil
}

// If such a key exists, the value will be overwritten without error
func (mn *RedisDataManager) createUserRedis(realmName string, userName string, userJson string) error {
	userKey := sf.Format(userKeyTemplate, mn.namespace, realmName, userName)
	if err := setString(mn.redisClient, mn.ctx, mn.logger, User, userKey, userJson); err != nil {
		// TODO(SIA) add log
		return fmt.Errorf("setString failed: %w", err)
	}
	return nil
}

// Returns an error if the user is in the realm
func (mn *RedisDataManager) addUserToRealm(realmName string, user data.User) error {
	userId := user.GetId()
	userName := user.GetUsername()
	realmUser := data.ExtendedIdentifier{
		ID:   userId,
		Name: userName,
	}
	sliceRealmUser := []data.ExtendedIdentifier{realmUser}
	if err := mn.createRealmUsers(realmName, sliceRealmUser, false); err != nil {
		return fmt.Errorf("createRealmUsers failed: %w", err)
	}
	return nil
}

// Adds users to the realm. If the argument isAllPreDelete = true, all other users will be deleted before they are added
func (mn *RedisDataManager) createRealmUsers(realmName string, realmUsers []data.ExtendedIdentifier, isAllPreDelete bool) error {
	bytesRealmUsers, err := json.Marshal(realmUsers)
	if err != nil {
		mn.logger.Error(sf.Format("An error occurred during realmUsers unmarshall"))
		return fmt.Errorf("json.Marshal failed: %w", err)
	}

	if isAllPreDelete {
		if err := mn.deleteRealmUsersRedis(realmName); err != nil {
			if err != nil && !errors.Is(err, errors_managers.ErrNotExists) {
				return fmt.Errorf("deleteRealmUsersRedis failed: %w", err)
			}
		}
	}
	realmUsersKey := sf.Format(realmUsersKeyTemplate, mn.namespace, realmName)
	redisIntCmd := mn.redisClient.RPush(mn.ctx, realmUsersKey, string(bytesRealmUsers)) // TODO(SIA) переписать по аналогии с другими
	if redisIntCmd.Err() != nil {
		// TODO(SIA) add log
		return redisIntCmd.Err()
	}
	return nil
}

func (mn *RedisDataManager) deleteUserRedis(realmName string, userName string) error {
	userKey := sf.Format(userKeyTemplate, mn.namespace, realmName, userName)
	if err := delKey(mn.redisClient, mn.ctx, mn.logger, User, userKey); err != nil {
		return fmt.Errorf("delKey failed: %w", err)
	}
	return nil
}

// Deletes user from realmUsers, does not delete user. Will return an error if there is no user in realm
func (mn *RedisDataManager) deleteUserFromRealm(realmName string, userName string) error {
	// TODO(SIA) Много действий происходит, для удаления user: происходит получение users, нахождение user, удаление его из массива,
	// удаление всех users из редис, добавление нового массива users в редис
	realmUsers, err := mn.getRealmUsers(realmName)
	if err != nil {
		return fmt.Errorf("getRealmUsers failed: %w", err)
	}

	isHasUser := false
	for i := range realmUsers {
		if realmUsers[i].Name == userName {
			isHasUser = true
			if i != (len(realmUsers) - 1) {
				realmUsers[i] = realmUsers[len(realmUsers)-1]
			}
			realmUsers = realmUsers[:len(realmUsers)-1]
			break
		}
	}
	if !isHasUser {
		// TODO(SIA) add log ("realm \"%s\" doesn't have user \"%s\" in Redis: %w", realmName, userName, errors_repo.ErrNotFound)
		return errors_managers.ErrNotFound
	}
	if err := mn.createRealmUsers(realmName, realmUsers, true); err != nil {
		return fmt.Errorf("createRealmClients failed: %w", err)
	}
	return nil
}
