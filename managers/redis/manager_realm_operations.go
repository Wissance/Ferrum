package redis

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/wissance/Ferrum/data"
	errors2 "github.com/wissance/Ferrum/errors"
	sf "github.com/wissance/stringFormatter"
)

// GetRealm function for getting realm by name, returns the realm with clients but no users.
/* This function constructs Redis key by pattern combines namespace and realm name (realmKeyTemplate). Unlike from FILE provider.
 * Realm stored in Redis does not have Clients and Users inside Realm itself, these objects must be queried separately.
 * Parameters:
 *     - realmName name of a realm
 * Returns: realm and error
 */
func (mn *RedisDataManager) GetRealm(realmName string) (*data.Realm, error) {
	realm, err := mn.getRealmObject(realmName)
	if err != nil {
		return nil, fmt.Errorf("getRealmObject failed: %w", err)
	}

	// should get realms too
	// if realms were stored without clients (we expected so), get clients related to realm and assign here
	clients, err := mn.GetClients(realmName)
	if err != nil {
		return nil, fmt.Errorf("GetClients failed: %w", err)
	}
	realm.Clients = clients

	return realm, nil
}

// CreateRealm - creates a realm, if the realm has users and clients, they will also be created.
/*
 * Arguments:
 *    - newRealm
 * Returns: error
 */
func (mn *RedisDataManager) CreateRealm(newRealm data.Realm) error {
	// TODO(SIA) Add transaction
	// TODO(SIA) use function isExists
	_, err := mn.getRealmObject(newRealm.Name)
	if err == nil {
		return errors2.ErrExists
	}
	if !errors.Is(err, errors2.ErrNotFound) {
		return fmt.Errorf("getRealmObject failed: %w", err)
	}

	if len(newRealm.Clients) != 0 {
		realmClients := make([]data.ExtendedIdentifier, len(newRealm.Clients))
		for i, client := range newRealm.Clients {
			bytesClient, err := json.Marshal(client)
			if err != nil {
				mn.logger.Error(sf.Format("An error occurred during Marshal Client"))
				return fmt.Errorf("json.Marshal failed: %w", err)
			}
			if err := mn.upsertClientObject(newRealm.Name, client.Name, string(bytesClient)); err != nil {
				return fmt.Errorf("upsertClientObject failed: %w", err)
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
			if err := mn.upsertUserObject(newRealm.Name, newUserName, newUser.GetJsonString()); err != nil {
				return fmt.Errorf("upsertUserObject failed: %w", err)
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
		mn.logger.Error(sf.Format("An error occurred during Marshal Realm"))
		return fmt.Errorf("json.Marsha failed: %w", err)
	}
	if err := mn.upsertRealmObject(newRealm.Name, string(jsonShortRealm)); err != nil {
		return fmt.Errorf("upsertRealmObject failed: %w", err)
	}
	return nil
}

// DeleteRealm - deleting the realm with all its clients and users
/*
 * Arguments:
 *    - realmName
 * Returns: error
 */
func (mn *RedisDataManager) DeleteRealm(realmName string) error {
	// TODO(SIA) Add transaction
	if err := mn.deleteRealmObject(realmName); err != nil {
		return fmt.Errorf("deleteRealmObject failed: %w", err)
	}

	clients, err := mn.getRealmClients(realmName)
	if err != nil {
		if !errors.Is(err, errors2.ErrZeroLength) {
			return fmt.Errorf("getRealmClients failed: %w", err)
		}
	} else {
		// TODO(SIA) overwrite to delete all keys at once
		for _, client := range clients {
			if err := mn.deleteClientObject(realmName, client.Name); err != nil {
				return fmt.Errorf("deleteClientObject failed: %w", err)
			}
		}
		if err := mn.deleteRealmClientsObject(realmName); err != nil {
			return fmt.Errorf("deleteRealmClientsObject failed: %w", err)
		}
	}

	users, err := mn.getRealmUsers(realmName)
	if err != nil {
		if !errors.Is(err, errors2.ErrZeroLength) {
			return fmt.Errorf("getRealmClients failed: %w", err)
		}
	} else {
		// TODO(SIA) overwrite to delete all keys at once
		for _, user := range users {
			if err := mn.deleteUserObject(realmName, user.Name); err != nil {
				return fmt.Errorf("deleteUserObject failed: %w", err)
			}
		}
		if err := mn.deleteRealmUsersObject(realmName); err != nil {
			return fmt.Errorf("deleteRealmUsersObject failed: %w", err)
		}
	}

	return nil
}

// UpdateRealm - realm update. It is expected that realmValue will not contain clients and users.
/* If the name or id of the realm has changed.  Then this information will be cascaded to all dependent objects.
 * Arguments:
 *    - realmName
 *    - realmNew
 * Returns: error
 */
func (mn *RedisDataManager) UpdateRealm(realmName string, realmNew data.Realm) error {
	// TODO(SIA) Add transaction
	oldRealm, err := mn.getRealmObject(realmName)
	if err != nil {
		return fmt.Errorf("getRealmObject failed: %w", err)
	}
	if oldRealm.Name != realmNew.Name {
		// TODO(SIA) use function isExists
		_, err := mn.getRealmObject(realmNew.Name)
		if err == nil {
			mn.logger.Error(sf.Format("Realm with a new name \"{0}\" already exists in Redis", realmNew.Name))
			return errors2.ErrExists
		}
		if !errors.Is(err, errors2.ErrNotFound) {
			return fmt.Errorf("getRealmObject failed: %w", err)
		}

		clients, err := mn.GetClients(oldRealm.Name)
		if err != nil && !errors.Is(err, errors2.ErrZeroLength) {
			return fmt.Errorf("GetClients failed: %w", err)
		}
		users, err := mn.GetUsers(oldRealm.Name)
		if err != nil && !errors.Is(err, errors2.ErrZeroLength) {
			return fmt.Errorf("GetUsers failed: %w", err)
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
		mn.logger.Error(sf.Format("An error occurred during Marshal Realm"))
		return fmt.Errorf("json.Marshal failed: %w", err)
	}
	if err := mn.upsertRealmObject(shortRealm.Name, string(jsonShortRealm)); err != nil {
		return fmt.Errorf("upsertRealmObject failed: %w", err)
	}
	return nil
}

// getRealmObject - getting realm without clients and users
/*
 * Arguments:
 *    - realmName
 * Returns: *Realm, error
 */
func (mn *RedisDataManager) getRealmObject(realmName string) (*data.Realm, error) {
	realmKey := sf.Format(realmKeyTemplate, mn.namespace, realmName)
	realm, err := getSingleRedisObject[data.Realm](mn.redisClient, mn.ctx, mn.logger, Realm, realmKey)
	if err != nil {
		if errors.Is(err, errors2.ErrNotFound) {
			mn.logger.Debug(sf.Format("Redis does not have Realm: \"{0}\"", realmName))
		}
		return nil, fmt.Errorf("getSingleRedisObject failed: %w", err)
	}
	return realm, nil
}

// upsertRealmObject - create or update a realm without clients and users
/* If such a key exists, the value will be overwritten without error
 * Arguments:
 *    - realmName
 *    - realmJson - string
 * Returns: *Realm, error
 */
func (mn *RedisDataManager) upsertRealmObject(realmName string, realmJson string) error {
	realmKey := sf.Format(realmKeyTemplate, mn.namespace, realmName)
	if err := mn.upsertRedisString(Realm, realmKey, realmJson); err != nil {
		return fmt.Errorf("upsertRedisString failed: %w", err)
	}
	return nil
}

// deleteRealmObject - deleting a realm without clients and users
/* Inside uses realmKeyTemplate
 * Arguments:
 *    - realmName
 * Returns: error
 */
func (mn *RedisDataManager) deleteRealmObject(realmName string) error {
	realmKey := sf.Format(realmKeyTemplate, mn.namespace, realmName)
	if err := mn.deleteRedisObject(Realm, realmKey); err != nil {
		return fmt.Errorf("deleteRedisObject failed: %w", err)
	}
	return nil
}

// deleteRealmClientsObject - deleting realmClients only
/* Inside uses realmClientsKeyTemplate
 * Arguments:
 *    - realmName
 * Returns: error
 */
func (mn *RedisDataManager) deleteRealmClientsObject(realmName string) error {
	realmClientsKey := sf.Format(realmClientsKeyTemplate, mn.namespace, realmName)
	if err := mn.deleteRedisObject(RealmClients, realmClientsKey); err != nil {
		return fmt.Errorf("deleteRedisObject failed: %w", err)
	}
	return nil
}

// deleteRealmUsersObject - deleting realmUsers only
/* Inside uses realmUsersKeyTemplate
 * Arguments:
 *    - realmName
 * Returns: error
 */
func (mn *RedisDataManager) deleteRealmUsersObject(realmName string) error {
	realmUsersKey := sf.Format(realmUsersKeyTemplate, mn.namespace, realmName)
	if err := mn.deleteRedisObject(RealmUsers, realmUsersKey); err != nil {
		return fmt.Errorf("deleteRedisObject failed: %w", err)
	}
	return nil
}
