package data

import "github.com/google/uuid"

type KeyCloakUser struct {
	rawData interface{}
}

func Create(rawData interface{}) User {
	kcUser := &KeyCloakUser{rawData: rawData}
	user := User(kcUser)
	return user
}

func (user *KeyCloakUser) GetUsername() string {
	return ""
}
func (user *KeyCloakUser) GetPassword() string {
	return ""
}

func (user *KeyCloakUser) GetId() uuid.UUID {
	return uuid.Nil
}
