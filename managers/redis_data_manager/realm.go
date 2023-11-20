package redis_data_manager

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/managers/errors_managers"
	sf "github.com/wissance/stringFormatter"
)

// GetRealm function for getting realm by name, returns the realm without clients and users. Unlike from FILE provider
// Realm stored in Redis does not have Clients and Users inside Realm itself, these objects must be queried separately.
func (mn *RedisDataManager) GetRealm(realmName string) (*data.Realm, error) {
	realmKey := sf.Format(realmKeyTemplate, mn.namespace, realmName)
	realm, err := getObjectFromRedis[data.Realm](mn.redisClient, mn.ctx, mn.logger, Realm, realmKey)
	if err != nil {
		return nil, fmt.Errorf("getObjectFromRedis failed: %w", err)
	}
	if realm == nil {
		mn.logger.Error(sf.Format("Redis does not have Realm: \"{0}\"", realmName))
		return nil, errors_managers.ErrNotFound
	}
	return realm, nil
}

func (mn *RedisDataManager) GetRealmWithClients(realmName string) (*data.Realm, error) {
	realm, err := mn.GetRealm(realmName)
	if err != nil {
		return nil, fmt.Errorf("GetRealm failed: %w", err)
	}

	// should get realms too
	// if realms were stored without clients (we expected so), get clients related to realm and assign here
	clients, err := mn.GetClientsFromRealm(realmName)
	if err != nil {
		return nil, fmt.Errorf("GetClientsFromRealm failed: %w", err)
	}
	realm.Clients = clients

	return realm, nil
}

// Creates a realm, if the realm has users and clients, they will also be created.
// If realm already existed, an error will be returned. If at least one client or user existed, an error will be returned as well
func (mn *RedisDataManager) CreateRealm(realmValue []byte) (*data.Realm, error) {
	// TODO(SIA) транзакции
	var newRealm data.Realm
	err := json.Unmarshal(realmValue, &newRealm)
	if err != nil {
		mn.logger.Error(sf.Format("An error occurred during Realm unmarshall"))
		return nil, fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	_, err = mn.GetRealm(newRealm.Name)
	if err == nil {
		return nil, errors_managers.ErrExists
	}
	if !errors.Is(err, errors_managers.ErrNotFound) {
		return nil, fmt.Errorf("GetRealm failed: %w", err)
	}

	if len(newRealm.Clients) != 0 {
		realmClients := make([]data.ExtendedIdentifier, len(newRealm.Clients))
		for i, client := range newRealm.Clients {
			bytesClient, err := json.Marshal(client)
			if err != nil {
				mn.logger.Error(sf.Format("An error occurred during Client marshal")) // TODO(SIA) ADD NAME
				return nil, fmt.Errorf("json.Marshal failed: %w", err)
			}
			newClient, err := mn.CreateClient(bytesClient)
			if err != nil {
				return nil, fmt.Errorf("CreateClient failed: %w", err)
			}
			realmClients[i] = data.ExtendedIdentifier{
				ID:   newClient.ID,
				Name: newClient.Name,
			}
		}
		if err := mn.createRealmClients(newRealm.Name, realmClients, true); err != nil {
			return nil, fmt.Errorf("createRealmClients failed: %w", err)
		}
	}

	if len(newRealm.Users) != 0 {
		realmUsers := make([]data.ExtendedIdentifier, len(newRealm.Users))
		for i, user := range newRealm.Users {
			bytesUser, err := json.Marshal(user)
			if err != nil {
				mn.logger.Error(sf.Format("An error occurred during User marshal")) // TODO(SIA) ADD NAME
				return nil, fmt.Errorf("json.Marshal failed: %w", err)
			}
			newUser, err := mn.CreateUser(bytesUser)
			if err != nil {
				return nil, fmt.Errorf("CreateUser failed: %w", err)
			}
			userName := newUser.GetUsername()
			userId := newUser.GetId()
			realmUsers[i] = data.ExtendedIdentifier{
				ID:   userId,
				Name: userName,
			}
		}
		if err := mn.createRealmUsers(newRealm.Name, realmUsers, true); err != nil {
			return nil, fmt.Errorf("createRealmUsers failed: %w", err)
		}
	}

	shortRealm := data.Realm{
		Name:                   newRealm.Name,
		Clients:                []data.Client{},
		Users:                  []any{},
		TokenExpiration:        newRealm.TokenExpiration,
		RefreshTokenExpiration: newRealm.RefreshTokenExpiration,
	}
	jsonShortRealm, err := json.Marshal(shortRealm)
	if err != nil {
		mn.logger.Error(sf.Format("An error occurred during Realm marshal")) // TODO(SIA) ADD NAME
		return nil, fmt.Errorf("json.Marsha failed: %w", err)
	}
	if err := mn.createRealmRedis(newRealm.Name, string(jsonShortRealm)); err != nil {
		// TODO(SIA) add log
		return nil, fmt.Errorf("createRealmRedis failed: %w", err)
	}
	return &newRealm, nil
}

// Removes realmClietns and realmUsers. Does not delete clients and users
func (mn *RedisDataManager) DeleteRealm(realmName string) error {
	// TODO(SIA) транзакции
	// TODO добавить ошибку если такого realm нет. Проверить возникает ли ошибка, при удалении не существующего ключа
	realmClientsKey := sf.Format(realmClientsKeyTemplate, mn.namespace, realmName)
	redisIntCmd := mn.redisClient.Del(mn.ctx, realmClientsKey)
	if redisIntCmd.Err() != nil {
		// TODO(SIA) add log
		return redisIntCmd.Err()
	}

	realmUsersKey := sf.Format(realmUsersKeyTemplate, mn.namespace, realmName)
	redisIntCmd = mn.redisClient.Del(mn.ctx, realmUsersKey)
	if redisIntCmd.Err() != nil {
		// TODO(SIA) add log
		return redisIntCmd.Err()
	}

	realmKey := sf.Format(realmKeyTemplate, mn.namespace, realmName)
	redisIntCmd = mn.redisClient.Del(mn.ctx, realmKey)
	if redisIntCmd.Err() != nil {
		// TODO(SIA) add log
		return redisIntCmd.Err()
	}

	return nil
}

func (mn *RedisDataManager) UpdateRealm(realmName string, realmValue []byte) (*data.Realm, error) {
	oldRealm, err := mn.GetRealm(realmName)
	if err != nil {
		return nil, fmt.Errorf("GetRealm failed: %w", err)
	}
	var newRealm data.Realm
	if err := json.Unmarshal(realmValue, &newRealm); err != nil {
		mn.logger.Error(sf.Format("An error occurred during Realm unmarshall"))
		return nil, fmt.Errorf("json.Unmarshal failed: %w", err)
	}
	if oldRealm.Name != newRealm.Name {
		// TODO(SIA) каскадно обновлять информацию у всех клиентов и user у realm. И удалить сам realm
	}

	shortRealm := data.Realm{
		Name:                   newRealm.Name,
		Clients:                []data.Client{},
		Users:                  []any{},
		TokenExpiration:        newRealm.TokenExpiration,
		RefreshTokenExpiration: newRealm.RefreshTokenExpiration,
	}
	jsonShortRealm, err := json.Marshal(shortRealm)
	if err != nil {
		mn.logger.Error(sf.Format("An error occurred during Realm marshal")) // TODO(SIA) ADD NAME
		return nil, fmt.Errorf("json.Marshal failed: %w", err)
	}
	if err := mn.createRealmRedis(shortRealm.Name, string(jsonShortRealm)); err != nil {
		// TODO(SIA) add log
		return nil, fmt.Errorf("createRealmRedis failed: %w", err)
	}
	return &shortRealm, nil
}

// If such a key exists, the value will be overwritten without error
func (mn *RedisDataManager) createRealmRedis(realmName string, realmJson string) error {
	realmKey := sf.Format(realmKeyTemplate, mn.namespace, realmName)
	if err := setString(mn.redisClient, mn.ctx, mn.logger, Realm, realmKey, realmJson); err != nil {
		// TODO(SIA) add log
		return fmt.Errorf("setString failed: %w", err)
	}
	return nil
}
