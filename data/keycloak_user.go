package data

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ohler55/ojg/jp"
)

// KeyCloakUser this structure is for user data that looks similar to KeyCloak, Users in Keycloak have info field with preferred_username and sub
// and others fields, Ferrum users have credentials built-in in user (temporary it stores in non encrypted mode)
type KeyCloakUser struct {
	rawData     interface{}
	jsonRawData string
}

// CreateUser function creates User interface instance (KeyCloakUser) from raw json
/* Function create User instance from any json (interface{})
 * Parameters:
 *    - rawData - any json
 * Return: instance of User as KeyCloakUser
 */
func CreateUser(rawData interface{}) User {
	jsonData, _ := json.Marshal(&rawData)
	kcUser := &KeyCloakUser{rawData: rawData, jsonRawData: string(jsonData)}
	user := User(kcUser)
	return user
}

// GetUsername returns username as it stores in KeyCloak
/* this function use internal map to navigate over info.preferred_username keys, the last one key is a login key
 * We are expecting that username is unique in the Realm
 * Parameters: no
 * Returns: username
 */
func (user *KeyCloakUser) GetUsername() string {
	return getPathStringValue[string](user.rawData, "info.preferred_username")
}

// GetPassword returns password
/* this function use internal map to navigate over credentials.password keys to retrieve a password
 * Parameters: no
 * Returns: password
 */
func (user *KeyCloakUser) GetPassword() string {
	return getPathStringValue[string](user.rawData, "credentials.password")
}

// GetId returns unique user identifier
/* this function use internal map to navigate over info.sun keys to retrieve a user id
 * Parameters: no
 * Returns: user id
 */
func (user *KeyCloakUser) GetId() uuid.UUID {
	idStrValue := getPathStringValue[string](user.rawData, "info.sub")
	id, err := uuid.Parse(idStrValue)
	if err != nil {
		// todo(UMV): think what to do here, return error!
	}
	return id
}

// GetUserInfo returns Json with all non-confidential user data as KeyCloak do
/* this function use internal map to navigate over key info ant retrieve all public userinfo
 * Parameters: no
 * Returns: user info
 */
func (user *KeyCloakUser) GetUserInfo() interface{} {
	result := getPathStringValue[interface{}](user.rawData, "info")
	return result
}

// getPathStringValue is a generic function to get actually map by key, key represents as a jsonpath navigation property
/* this function uses json path to navigate over nested maps and return any required type
 * Parameters:
 *    - rawData - json object
 *    - path - json path to retrieve part of json with specified type (T)
 * Returns: part of json
 */
func getPathStringValue[T any](rawData interface{}, path string) T {
	var result T
	mask, err := jp.ParseString(path)
	if err != nil {
		// todo(UMV): log and think what to do ...
	}
	res := mask.Get(rawData)
	if res != nil && len(res) == 1 {
		result = res[0].(T)
	}
	return result
}
