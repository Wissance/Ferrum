package data

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ohler55/ojg/jp"
)

type KeyCloakUser struct {
	rawData     interface{}
	jsonRawData string
}

func CreateUser(rawData interface{}) User {
	jsonData, _ := json.Marshal(&rawData)
	kcUser := &KeyCloakUser{rawData: rawData, jsonRawData: string(jsonData)}
	user := User(kcUser)
	return user
}

func (user *KeyCloakUser) GetUsername() string {
	return getPathStringValue[string](user.rawData, "info.preferred_username")
}
func (user *KeyCloakUser) GetPassword() string {
	return getPathStringValue[string](user.rawData, "credentials.password")
}

func (user *KeyCloakUser) GetId() uuid.UUID {
	idStrValue := getPathStringValue[string](user.rawData, "info.sub")
	id, err := uuid.Parse(idStrValue)
	if err != nil {
		//todo(UMV): think what to do here
	}
	return id
}

func (user *KeyCloakUser) GetUserInfo() interface{} {
	var jsonResult interface{}
	result := getPathStringValue[interface{}](user.rawData, "info")
	str, _ := json.Marshal(&result)
	_ = json.Unmarshal(str, &jsonResult)
	return jsonResult
}

func /*(user *KeyCloakUser)*/ getPathStringValue[T any](rawData interface{}, path string) T {
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
