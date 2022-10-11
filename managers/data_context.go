package managers

import (
	"github.com/google/uuid"
	"github.com/wissance/Ferrum/data"
)

type DataContext interface {
	GetRealm(realmName string) *data.Realm
	GetClient(realm *data.Realm, name string) *data.Client
	GetUser(realm *data.Realm, userName string) *data.User
	GetUserById(realm *data.Realm, userId uuid.UUID) *data.User
	GetRealmUsers(realmName string) *[]data.User
}
