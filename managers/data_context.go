package managers

import (
	"Ferrum/data"
	"github.com/google/uuid"
)

type DataContext interface {
	GetRealm(realmName string) *data.Realm
	GetClient(realm *data.Realm, name string) *data.Client
	GetUser(realm *data.Realm, userName string) *data.User
	GetUserById(realm *data.Realm, userId uuid.UUID) *data.User
	GetRealmUsers(realmName string) *[]data.User
}
