package managers

import (
	"Ferrum/data"
	"github.com/google/uuid"
)

type DataContext interface {
	GetRealm(realmId string) *data.Realm
	GetClient(realm *data.Realm, name string) *data.Client
	GetUser(realm *data.Realm, userName string) *data.User
	GetClientUsers(clientId uuid.UUID) *[]data.User
}
