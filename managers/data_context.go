package managers

import (
	"Ferrum/data"
	"github.com/google/uuid"
)

type DataContext interface {
	Init() bool
	GetRealm(realmId *string) *data.Realm
	GetClient(clientId uuid.UUID) *data.Client
	GetUser(userId uuid.UUID) *data.User
	GetClientUsers(clientId uuid.UUID) *[]data.User
}
