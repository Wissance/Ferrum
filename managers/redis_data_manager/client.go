package redis_data_manager

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/managers/errors_managers"
	sf "github.com/wissance/stringFormatter"
)

func (mn *RedisDataManager) GetClientsFromRealm(realmName string) ([]data.Client, error) {
	realmClients, err := mn.getRealmClients(realmName)
	if err != nil {
		return nil, fmt.Errorf("getRealmClients failed: %w", err)
	}
	clients := make([]data.Client, len(realmClients))
	for i, rc := range realmClients {
		// todo(UMV) get all them at once
		client, err := mn.GetClient(rc.Name)
		if err != nil {
			if errors.Is(err, errors_managers.ErrNotFound) { // TODO(SIA) check
				mn.logger.Error(sf.Format("Realm: \"{0}\" has client: \"{1}\", that Redis does not have", realmName, rc.Name))
			}
			return nil, fmt.Errorf("GetClient failed: %w", err)
		}
		clients[i] = *client
	}
	return clients, nil
}

func (mn *RedisDataManager) GetClient(clientName string) (*data.Client, error) {
	clientKey := sf.Format(clientKeyTemplate, mn.namespace, clientName)
	client, err := getObjectFromRedis[data.Client](mn.redisClient, mn.ctx, mn.logger, Client, clientKey)
	if err != nil {
		return nil, fmt.Errorf("getObjectFromRedis failed: %w", err)
	}
	if client == nil {
		mn.logger.Error(sf.Format("Redis does not have Client: \"{0}\"", clientName))
		return nil, errors_managers.ErrNotFound
	}
	return client, nil
}

func (mn *RedisDataManager) GetClientFromRealm(realmName string, clientName string) (*data.Client, error) {
	realmClient, err := mn.getRealmClient(realmName, clientName)
	if err != nil {
		return nil, fmt.Errorf("getRealmClient failed: %w", err)
	}
	client, err := mn.GetClient(realmClient.Name)
	if err != nil {
		if errors.Is(err, errors_managers.ErrNotFound) { // TODO(SIA) check
			mn.logger.Error(sf.Format("Realm: \"{0}\" has client: \"{1}\", that Redis does not have", realmName, clientName))
		}
		return nil, fmt.Errorf("GetClient failed: %w", err)
	}
	return client, nil
}

// Returns an error if the client exists in redis
func (mn *RedisDataManager) CreateClient(clientValue []byte) (*data.Client, error) {
	// TODO(SIA) транзакции
	// TODO(SIA) возможно нужно проверять, что есть какие-то поля у clients
	var clientNew data.Client
	err := json.Unmarshal(clientValue, &clientNew)
	if err != nil {
		mn.logger.Error(sf.Format("An error occurred during Client unmarshall"))
		return nil, fmt.Errorf("json.Unmarshal failed: %w", err)
	}
	_, err = mn.GetClient(clientNew.Name)
	if err == nil {
		return nil, errors_managers.ErrExists
	}
	if !errors.Is(err, errors_managers.ErrNotFound) {
		return nil, fmt.Errorf("GetClient failed: %w", err)
	}

	err = mn.createClientRedis(clientNew.Name, string(clientValue))
	if err != nil {
		return nil, fmt.Errorf("createClientRedis failed: %w", err)
	}
	return &clientNew, nil
}

// Returns an error if the client is in the realm
func (mn *RedisDataManager) AddClientToRealm(realmName string, clientName string) error {
	_, err := mn.getRealmClient(realmName, clientName)
	if err == nil {
		return errors_managers.ErrExists
	}
	if !errors.Is(err, errors_managers.ErrNotFound) {
		return fmt.Errorf("getRealmClient failed: %w", err)
	}

	client, err := mn.GetClient(clientName)
	if err != nil {
		return fmt.Errorf("GetClient failed: %w", err)
	}
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

func (mn *RedisDataManager) DeleteClient(clientName string) error {
	// TODO(SIA) add cascading deletion to all realms
	clientKey := sf.Format(clientKeyTemplate, mn.namespace, clientName)
	redisIntCmd := mn.redisClient.Del(mn.ctx, clientKey)
	if redisIntCmd.Err() != nil {
		// TODO(SIA) add log
		return redisIntCmd.Err() // TODO(SIA) проверить, будет ли ошибка, если нет такого клиента
	}
	return nil
}

// Deletes client from realmClients, does not delete client. Will return an error if there is no client in realm
func (mn *RedisDataManager) DeleteClientFromRealm(realmName string, clientName string) error {
	// TODO(SIA) Много действий происходит, для удаления клиента: происходит получение клиентов, нахождение клиента, удаление его из массива,
	// удаление всех клиентов из редис, добавление нового массива клиентов в редис
	realmClients, err := mn.getRealmClients(realmName)
	if err != nil {
		return fmt.Errorf("getRealmClients failed: %w", err)
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
		// TODO(SIA) add log ("realm \"%s\" doesn't have client \"%s\" in Redis", realmName, clientName)
		return errors_managers.ErrNotFound
	}
	if err := mn.createRealmClients(realmName, realmClients, true); err != nil {
		return fmt.Errorf("createRealmClients failed: %w", err)
	}
	return nil
}

func (mn *RedisDataManager) UpdateClient(clientName string, clientValue []byte) (*data.Client, error) {
	// TODO(SIA) транзакции
	oldClient, err := mn.GetClient(clientName)
	if err != nil {
		return nil, fmt.Errorf("GetClient failed: %w", err)
	}
	var newClient data.Client
	if err := json.Unmarshal(clientValue, &newClient); err != nil {
		mn.logger.Error(sf.Format("An error occurred during Client unmarshall"))
		return nil, fmt.Errorf("json.Unmarshal failed: %w", err)
	}
	if newClient.ID != oldClient.ID || newClient.Name != oldClient.Name {
		// TODO(SIA) каскадно обновлять информацию во всех realm где был этот клиент. И удалить сам клиент, т.к.
		// следующее создание через setString не перезапишет старый клиент с прошлым именем
	}
	if err := mn.createClientRedis(newClient.Name, string(clientValue)); err != nil {
		return nil, fmt.Errorf("createClientRedis failed: %w", err)
	}
	return &newClient, nil
}

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
		mn.logger.Debug(sf.Format("Realm: \"{0}\" doesn't have client: \"{1}\" in Redis", realmName, clientName))
		return nil, errors_managers.ErrNotFound
	}

	return &client, nil
}

func (mn *RedisDataManager) getRealmClients(realmName string) ([]data.ExtendedIdentifier, error) {
	realmClientsKey := sf.Format(realmClientsKeyTemplate, mn.namespace, realmName)
	realmClients, err := getObjectsListFromRedis[data.ExtendedIdentifier](mn.redisClient, mn.ctx, mn.logger, RealmClients, realmClientsKey)
	if err != nil {
		return nil, fmt.Errorf("getObjectsListFromRedis failed: %w", err)
	}
	if len(realmClients) == 0 {
		mn.logger.Error(sf.Format("There are no clients for realm: \"{0}\" in Redis", realmName))
		return nil, errors_managers.ErrZeroLength
	}
	return realmClients, nil
}

// If such a key exists, the value will be overwritten without error
func (mn *RedisDataManager) createClientRedis(clientName string, clientJson string) error {
	clientKey := sf.Format(clientKeyTemplate, mn.namespace, clientName)
	if err := setString(mn.redisClient, mn.ctx, mn.logger, Client, clientKey, clientJson); err != nil {
		// TODO(SIA) add log
		return fmt.Errorf("setString failed: %w", err)
	}
	return nil
}

// Adds clients to the realm. If the argument isAllPreDelete = true, all other clients will be deleted before they are added
func (mn *RedisDataManager) createRealmClients(realmName string, realmClients []data.ExtendedIdentifier, isAllPreDelete bool) error {
	bytesRealmClients, err := json.Marshal(realmClients)
	if err != nil {
		mn.logger.Error(sf.Format("An error occurred during realmClients unmarshall"))
		return fmt.Errorf("json.Marshal failed: %w", err)
	}
	realmClientsKey := sf.Format(realmClientsKeyTemplate, mn.namespace, realmName)
	if isAllPreDelete {
		redisIntCmd := mn.redisClient.Del(mn.ctx, realmClientsKey)
		if redisIntCmd.Err() != nil {
			// TODO(SIA) add log
			return redisIntCmd.Err()
		}
	}
	redisIntCmd := mn.redisClient.RPush(mn.ctx, realmClientsKey, string(bytesRealmClients))
	if redisIntCmd.Err() != nil {
		// TODO(SIA) add log
		return redisIntCmd.Err()
	}
	return nil
}
