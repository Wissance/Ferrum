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
		if errors.Is(err, errors_managers.ErrNotFound) {
			mn.logger.Debug(sf.Format("Redis does not have Realm: \"{0}\"", realmName))
		}
		return nil, fmt.Errorf("getObjectFromRedis failed: %w", err)
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
func (mn *RedisDataManager) CreateRealm(newRealm data.Realm) error {
	// TODO(SIA) транзакции
	_, err := mn.GetRealm(newRealm.Name) // TODO(SIA) use function isExists
	if err == nil {
		return errors_managers.ErrExists
	}
	if !errors.Is(err, errors_managers.ErrNotFound) {
		return fmt.Errorf("GetRealm failed: %w", err)
	}

	if len(newRealm.Clients) != 0 {
		realmClients := make([]data.ExtendedIdentifier, len(newRealm.Clients))
		for i, client := range newRealm.Clients {
			bytesClient, err := json.Marshal(client)
			if err != nil {
				mn.logger.Error(sf.Format("An error occurred during Client marshal")) // TODO(SIA) ADD NAME
				return fmt.Errorf("json.Marshal failed: %w", err)
			}
			if err := mn.createClientRedis(newRealm.Name, client.Name, string(bytesClient)); err != nil {
				return fmt.Errorf("createClientRedis failed: %w", err)
			}
			realmClients[i] = data.ExtendedIdentifier{
				ID:   client.ID,
				Name: client.Name,
			}
		}
		if err := mn.createRealmClients(newRealm.Name, realmClients, true); err != nil {
			return fmt.Errorf("createRealmClients failed: %w", err)
		}
	}

	if len(newRealm.Users) != 0 {
		realmUsers := make([]data.ExtendedIdentifier, len(newRealm.Users))
		for i, user := range newRealm.Users {
			newUser := data.CreateUser(user)
			newUserName := newUser.GetUsername()
			bytesUser, err := json.Marshal(user)
			if err != nil {
				mn.logger.Error(sf.Format("An error occurred during User marshal")) // TODO(SIA) ADD NAME
				return fmt.Errorf("json.Marshal failed: %w", err)
			}
			if err := mn.createUserRedis(newRealm.Name, newUserName, string(bytesUser)); err != nil {
				return fmt.Errorf("createUserRedis failed: %w", err)
			}
			newUserId := newUser.GetId()
			realmUsers[i] = data.ExtendedIdentifier{
				ID:   newUserId,
				Name: newUserName,
			}
		}
		if err := mn.createRealmUsers(newRealm.Name, realmUsers, true); err != nil {
			return fmt.Errorf("createRealmUsers failed: %w", err)
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
		return fmt.Errorf("json.Marsha failed: %w", err)
	}
	if err := mn.createRealmRedis(newRealm.Name, string(jsonShortRealm)); err != nil {
		// TODO(SIA) add log
		return fmt.Errorf("createRealmRedis failed: %w", err)
	}
	return nil
}

func (mn *RedisDataManager) DeleteRealm(realmName string) error {
	// TODO(SIA) транзакции
	if err := mn.deleteRealmRedis(realmName); err != nil {
		return fmt.Errorf("deleteRealmRedis failed: %w", err)
	}

	clients, err := mn.getRealmClients(realmName)
	if err != nil {
		if !errors.Is(err, errors_managers.ErrZeroLength) {
			return fmt.Errorf("getRealmClients failed: %w", err)
		}
	} else {
		// TODO(SIA) переписать на удаление сразу всех ключей
		for _, client := range clients {
			if err := mn.deleteClientRedis(realmName, client.Name); err != nil {
				return fmt.Errorf("deleteClientRedis failed: %w", err)
			}
		}
		if err := mn.deleteRealmClientsRedis(realmName); err != nil {
			return fmt.Errorf("deleteRealmClientsRedis failed: %w", err)
		}
	}

	users, err := mn.getRealmUsers(realmName)
	if err != nil {
		if !errors.Is(err, errors_managers.ErrZeroLength) {
			return fmt.Errorf("getRealmClients failed: %w", err)
		}
	} else {
		// TODO(SIA) переписать на удаление сразу всех ключей
		for _, user := range users {
			if err := mn.deleteUserRedis(realmName, user.Name); err != nil {
				return fmt.Errorf("deleteUserRedis failed: %w", err)
			}
		}
		if err := mn.deleteRealmUsersRedis(realmName); err != nil {
			return fmt.Errorf("deleteRealmUsersRedis failed: %w", err)
		}
	}

	return nil
}

// It is expected that realmValue will not contain clients and users
func (mn *RedisDataManager) UpdateRealm(realmName string, realmNew data.Realm) error {
	// TODO(SIA) транзакции
	oldRealm, err := mn.GetRealm(realmName)
	if err != nil {
		return fmt.Errorf("GetRealm failed: %w", err)
	}
	if oldRealm.Name != realmNew.Name {
		_, err := mn.GetRealm(realmNew.Name) // TODO(SIA) use function isExists
		if err == nil {
			mn.logger.Error(sf.Format("Realm with a new name \"{0}\" already exists in Redis", realmNew.Name))
			return errors_managers.ErrExists
		}
		if !errors.Is(err, errors_managers.ErrNotFound) {
			return fmt.Errorf("GetRealm failed: %w", err)
		}

		clients, err := mn.GetClientsFromRealm(oldRealm.Name)
		if err != nil && !errors.Is(err, errors_managers.ErrZeroLength) {
			return fmt.Errorf("GetClientsFromRealm failed: %w", err)
		}
		users, err := mn.GetUsersFromRealm(oldRealm.Name)
		if err != nil && !errors.Is(err, errors_managers.ErrZeroLength) {
			return fmt.Errorf("GetUsersFromRealm failed: %w", err)
		}
		usersData := make([]any, len(users))
		for i, u := range users {
			usersData[i] = u.GetRawData()
		}
		newRealmWithOldClientsAndUsers := data.Realm{
			Name:                   realmNew.Name,
			Clients:                clients,
			Users:                  usersData,
			TokenExpiration:        realmNew.TokenExpiration,
			RefreshTokenExpiration: realmNew.RefreshTokenExpiration,
		}
		if err := mn.DeleteRealm(oldRealm.Name); err != nil {
			return fmt.Errorf("DeleteRealm failed: %w", err)
		}
		if err = mn.CreateRealm(newRealmWithOldClientsAndUsers); err != nil {
			return fmt.Errorf("CreateRealm failed: %w", err)
		}
		return nil
	}

	shortRealm := data.Realm{
		Name:                   realmNew.Name,
		Clients:                []data.Client{},
		Users:                  []any{},
		TokenExpiration:        realmNew.TokenExpiration,
		RefreshTokenExpiration: realmNew.RefreshTokenExpiration,
	}
	jsonShortRealm, err := json.Marshal(shortRealm)
	if err != nil {
		mn.logger.Error(sf.Format("An error occurred during Realm marshal")) // TODO(SIA) ADD NAME
		return fmt.Errorf("json.Marshal failed: %w", err)
	}
	if err := mn.createRealmRedis(shortRealm.Name, string(jsonShortRealm)); err != nil {
		// TODO(SIA) add log
		return fmt.Errorf("createRealmRedis failed: %w", err)
	}
	return nil
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

func (mn *RedisDataManager) deleteRealmRedis(realmName string) error {
	realmKey := sf.Format(realmKeyTemplate, mn.namespace, realmName)
	if err := delKey(mn.redisClient, mn.ctx, mn.logger, Realm, realmKey); err != nil {
		return fmt.Errorf("delKey failed: %w", err)
	}
	return nil
}

func (mn *RedisDataManager) deleteRealmClientsRedis(realmName string) error {
	realmClientsKey := sf.Format(realmClientsKeyTemplate, mn.namespace, realmName)
	if err := delKey(mn.redisClient, mn.ctx, mn.logger, RealmClients, realmClientsKey); err != nil {
		return fmt.Errorf("delKey failed: %w", err)
	}
	return nil
}

func (mn *RedisDataManager) deleteRealmUsersRedis(realmName string) error {
	realmUsersKey := sf.Format(realmUsersKeyTemplate, mn.namespace, realmName)
	if err := delKey(mn.redisClient, mn.ctx, mn.logger, RealmUsers, realmUsersKey); err != nil {
		return fmt.Errorf("delKey failed: %w", err)
	}
	return nil
}
