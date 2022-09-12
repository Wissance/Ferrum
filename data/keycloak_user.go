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

func Create(rawData interface{}) User {
	jsonData, _ := json.Marshal(&rawData)
	kcUser := &KeyCloakUser{rawData: rawData, jsonRawData: string(jsonData)}
	user := User(kcUser)
	return user
}

func (user *KeyCloakUser) GetUsername() string {
	return user.getPathStringValue("info.preferred_username")
}
func (user *KeyCloakUser) GetPassword() string {
	return user.getPathStringValue("credentials.password")
}

func (user *KeyCloakUser) GetId() uuid.UUID {
	idStrValue := user.getPathStringValue("info.sub")
	id, err := uuid.Parse(idStrValue)
	if err != nil {
		//todo(UMV): think what to do here
	}
	return id
}

func (user *KeyCloakUser) GetUserInfo() interface{} {
	str := user.getPathStringValue("info")
	var result interface{}
	_ = json.Unmarshal([]byte(str), &result)
	return result
}

func (user *KeyCloakUser) getPathStringValue(path string) string {
	mask, err := jp.ParseString(path)
	if err != nil {
		// todo(UMV): log and think what to do ...
	}
	res := mask.Get(user.rawData)
	if res != nil && len(res) == 1 {
		return res[0].(string)
	}
	return ""
}
