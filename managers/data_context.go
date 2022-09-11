package managers

import (
	"Ferrum/data"
)

type DataContext interface {
	GetRealm(realmName string) *data.Realm
	GetClient(realm *data.Realm, name string) *data.Client
	GetUser(realm *data.Realm, userName string) *data.User
	GetRealmUsers(realmName string) *[]data.User
}
