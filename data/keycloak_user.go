package data

import "github.com/google/uuid"

type KeyCloakUser struct {
	raw interface{}
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
