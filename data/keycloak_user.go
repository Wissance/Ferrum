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
	mask, err := jp.ParseString("info.preferred_username")
	if err != nil {
		// todo(UMV): log and think what to do ...
	}
	res := mask.Get(user.rawData)
	if res != nil && len(res) == 1 {
		return res[0].(string)
	}
	return ""
}
func (user *KeyCloakUser) GetPassword() string {
	return ""
}

func (user *KeyCloakUser) GetId() uuid.UUID {
	return uuid.Nil
}
