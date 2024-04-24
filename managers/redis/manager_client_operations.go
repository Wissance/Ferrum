package redis

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/wissance/Ferrum/data"
	errors2 "github.com/wissance/Ferrum/errors"
	sf "github.com/wissance/stringFormatter"
)

// GetClients - getting clients from the specified realm
/*
 * Arguments:
 *    - realmName
 * Returns: slice of client, error
 */
func (mn *RedisDataManager) GetClients(realmName string) ([]data.Client, error) {
	realmClients, err := mn.getRealmClients(realmName)
	if err != nil {
		return nil, fmt.Errorf("getRealmClients failed: %w", err)
	}
	clients := make([]data.Client, len(realmClients))
	for i, rc := range realmClients {
		// todo(UMV) get all them at once
		client, readClientErr := mn.GetClient(realmName, rc.Name)
		if readClientErr != nil {
			if errors.Is(err, errors2.ObjectNotFoundError{}) {
				mn.logger.Error(sf.Format("Realm: \"{0}\" has client: \"{1}\", that Redis does not have", realmName, rc.Name))
			}
			return nil, readClientErr
		}
		clients[i] = *client
	}
	return clients, nil
}

// GetClient function for get realm client by name
/* This function constructs Redis key by pattern combines namespace and realm name and client name (clientKeyTemplate)
 * Parameters:
 *     - realmName
 *     - clientName
 * Returns: client and error
 */
func (mn *RedisDataManager) GetClient(realmName string, clientName string) (*data.Client, error) {
	clientKey := sf.Format(clientKeyTemplate, mn.namespace, realmName, clientName)
	client, err := getSingleRedisObject[data.Client](mn.redisClient, mn.ctx, mn.logger, Client, clientKey)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// CreateClient - new client creation
/* Returns an error if the client exists in redis
 * Arguments:
 *    - realmName
 *    - clientNew
 * Returns: error
 */
func (mn *RedisDataManager) CreateClient(realmName string, clientNew data.Client) error {
	// TODO(SIA) Add transaction
	// TODO(SIA) use function isExists
	_, err := mn.getRealmObject(realmName)
	if err != nil {
		return err
	}
	// TODO(SIA) use function isExists
	_, err = mn.GetClient(realmName, clientNew.Name)
	if err == nil {
		return errors2.ErrExists
	}
	if !errors.Is(err, errors2.ObjectNotFoundError{}) {
		return err
	}

	clientBytes, err := json.Marshal(clientNew)
	if err != nil {
		mn.logger.Error(sf.Format("An error occurred during Marshal Client"))
		return fmt.Errorf("json.Marshal failed: %w", err)
	}
	err = mn.upsertClientObject(realmName, clientNew.Name, string(clientBytes))
	if err != nil {
		return fmt.Errorf("upsertClientObject failed: %w", err)
	}

	if addClientErr := mn.addClientToRealm(realmName, clientNew); addClientErr != nil {
		return fmt.Errorf("addClientToRealm failed: %w", err)
	}

	return nil
}

// DeleteClient - deleting an existing client
/* It also deletes the client from realmClients
 * Arguments:
 *    - realmName
 *    - clientName
 * Returns: error
 */
func (mn *RedisDataManager) DeleteClient(realmName string, clientName string) error {
	if err := mn.deleteClientObject(realmName, clientName); err != nil {
		return fmt.Errorf("deleteClientObject failed: %w", err)
	}
	if err := mn.deleteClientFromRealm(realmName, clientName); err != nil {
		if errors.Is(err, errors2.ObjectNotFoundError{}) || errors.Is(err, errors2.ErrZeroLength) {
			return nil
		}
		return err
	}
	return nil
}

// UpdateClient - upgrading an existing client
/*
 * Arguments:
 *    - realmName
 *    - clientName
 *    - clientNew
 * Returns: error
 */
func (mn *RedisDataManager) UpdateClient(realmName string, clientName string, clientNew data.Client) error {
	// TODO(SIA) Add transaction
	oldClient, err := mn.GetClient(realmName, clientName)
	if err != nil {
		return fmt.Errorf("GetClient failed: %w", err)
	}
	if clientNew.ID != oldClient.ID || clientNew.Name != oldClient.Name {
		if err := mn.DeleteClient(realmName, oldClient.Name); err != nil {
			return fmt.Errorf("DeleteClient failed: %w", err)
		}
		if err := mn.addClientToRealm(realmName, clientNew); err != nil {
			return fmt.Errorf("addClientToRealm failed: %w", err)
		}
	}

	clientBytes, err := json.Marshal(clientNew)
	if err != nil {
		mn.logger.Error(sf.Format("An error occurred during Marshal Client"))
		return fmt.Errorf("json.Marshal failed: %w", err)
	}
	err = mn.upsertClientObject(realmName, clientNew.Name, string(clientBytes))
	if err != nil {
		return fmt.Errorf("upsertClientObject failed: %w", err)
	}

	return nil
}

// getRealmClients - get realmClients entity.
/* realmClientsKeyTemplate is used inside.
 * Arguments:
 *    - realmName
 * Returns: slice of ExtendedIdentifier, error
 */
func (mn *RedisDataManager) getRealmClients(realmName string) ([]data.ExtendedIdentifier, error) {
	realmClientsKey := sf.Format(realmClientsKeyTemplate, mn.namespace, realmName)
	realmClients, err := getObjectsListFromRedis[data.ExtendedIdentifier](mn.redisClient, mn.ctx, mn.logger, RealmClients, realmClientsKey)
	if err != nil {
		return nil, fmt.Errorf("getObjectsListFromRedis failed: %w", err)
	}
	return realmClients, nil
}

// getRealmClient - get ExtendedIdentifier entity.
/* First, getRealmClients happens. Then there are comparisons by name
 * Arguments:
 *    - realmName
 *    - clientName
 * Returns: *ExtendedIdentifier, error
 */
func (mn *RedisDataManager) getRealmClient(realmName string, clientName string) (*data.ExtendedIdentifier, error) {
	realmClients, err := mn.getRealmClients(realmName)
	if err != nil {
		return nil, fmt.Errorf("getRealmClients failed: %w", err)
	}

	realmHasClient := false
	var client data.ExtendedIdentifier
	for _, rc := range realmClients {
		if rc.Name == clientName {
			realmHasClient = true
			client = rc
			break
		}
	}
	if !realmHasClient {
		mn.logger.Debug(sf.Format("Realm: \"{0}\" doesn't have Client: \"{1}\" in Redis", realmName, clientName))
		return nil, errors2.NewObjectNotFoundError(string(Realm), realmName, "")
	}
	return &client, nil
}

// upsertClientObject - create or update a client
/* If such a key exists, the value will be overwritten without error
 * Arguments:
 *    - realmName
 *    - clientName
 *    - clientJson - string
 * Returns: error
 */
func (mn *RedisDataManager) upsertClientObject(realmName string, clientName string, clientJson string) error {
	clientKey := sf.Format(clientKeyTemplate, mn.namespace, realmName, clientName)
	if err := mn.upsertRedisString(Client, clientKey, clientJson); err != nil {
		return fmt.Errorf("upsertRedisString failed: %w", err)
	}
	return nil
}

// addClientToRealm - adding a client to the realmClient entity
/* Uses createRealmClients internally
 * Arguments:
 *    - realmName
 *    - client
 * Returns: error
 */
func (mn *RedisDataManager) addClientToRealm(realmName string, client data.Client) error {
	realmClient := data.ExtendedIdentifier{
		ID:   client.ID,
		Name: client.Name,
	}
	sliceRealmClient := []data.ExtendedIdentifier{realmClient}
	if err := mn.createRealmClients(realmName, sliceRealmClient, false); err != nil {
		return fmt.Errorf("createRealmClients failed: %w", err)
	}
	return nil
}

// createRealmClients - To add a new item to the list or create a new realmClients
/* Adds clients to the realm. If the argument isAllPreDelete = true, all other clients will be deleted before they are added
 * Arguments:
 *    - realmName
 *    - realmClients - slice of ExtendedIdentifier
 *    - isAllPreDelete - flag, If true, the already existing realmClients will be deleted. If false, new ones will be added to it.
 * Returns: error
 */
func (mn *RedisDataManager) createRealmClients(realmName string, realmClients []data.ExtendedIdentifier, isAllPreDelete bool) error {
	// TODO(SIA) maybe split into two functions
	bytesRealmClients, err := json.Marshal(realmClients)
	if err != nil {
		mn.logger.Error(sf.Format("An error occurred during Marshal realmClients"))
		return fmt.Errorf("json.Marshal failed: %w", err)
	}
	if isAllPreDelete {
		if err := mn.deleteRealmClientsObject(realmName); err != nil {
			if err != nil && !errors.Is(err, errors2.ErrNotExists) {
				return fmt.Errorf("deleteRealmClientsObject failed: %w", err)
			}
		}
	}
	realmClientsKey := sf.Format(realmClientsKeyTemplate, mn.namespace, realmName)
	if err := mn.appendStringToRedisList(RealmClients, realmClientsKey, string(bytesRealmClients)); err != nil {
		return fmt.Errorf("appendStringToRedisList failed: %w", err)
	}
	return nil
}

// deleteClientObject - deleting a client
/* Inside uses clientKeyTemplate
 * Arguments:
 *    - realmName
 *    - clientName
 * Returns: error
 */
func (mn *RedisDataManager) deleteClientObject(realmName string, clientName string) error {
	clientKey := sf.Format(clientKeyTemplate, mn.namespace, realmName, clientName)
	if err := mn.deleteRedisObject(Client, clientKey); err != nil {
		return fmt.Errorf("deleteRedisObject failed: %w", err)
	}
	return nil
}

// deleteClientFromRealm - deleting a client from realmClients entity
/* Deletes client from realmClients, does not delete client. Will return an error if there is no client in realm.
 * After deletion, all items in the list are merged into one.
 * A lot of things happen to delete a client: get clients, find the client, delete it from the array,
 * delete all clients from the realm, add a new array of clients to the realm.
 * Arguments:
 *    - realmName
 *    - clientName
 * Returns: error
 */
func (mn *RedisDataManager) deleteClientFromRealm(realmName string, clientName string) error {
	realmClients, err := mn.getRealmClients(realmName)
	if err != nil {
		mn.logger.Warn(sf.Format("deleteClientFromRealm failed: {0}", err.Error()))
		return err
	}

	isHasClient := false
	for i := range realmClients {
		if realmClients[i].Name == clientName {
			isHasClient = true
			if i != (len(realmClients) - 1) {
				realmClients[i] = realmClients[len(realmClients)-1]
			}
			realmClients = realmClients[:len(realmClients)-1]
			break
		}
	}
	if !isHasClient {
		return errors2.NewObjectNotFoundError(Client, clientName, sf.Format("realm: {0}", realmName))
	}
	if err := mn.createRealmClients(realmName, realmClients, true); err != nil {
		return fmt.Errorf("createRealmClients failed: %w", err)
	}
	return nil
}
