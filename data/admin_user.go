package data

import "github.com/google/uuid"

// AdminUser is a User that is not related to any Realm, this is an administration user
// In general this is a user that is using for server administration: i.e. initial data.Realm creation
type AdminUser struct {
	Id           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordSalt string    `json:"password_salt"`
	PasswordHash string    `json:"password_hash"`
}
