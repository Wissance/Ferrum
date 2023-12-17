package data

import "github.com/google/uuid"

// User is a common user interface with all Required methods to get information about user, in future we probably won't have GetPassword method
// because Password is not an only method for authentication
type User interface {
	GetUsername() string
	GetPassword() string
	SetPassword(password string) error
	GetId() uuid.UUID
	GetUserInfo() interface{}
	GetRawData() interface{}
	GetJsonString() string
}

var _ User = (*KeyCloakUser)(nil)
