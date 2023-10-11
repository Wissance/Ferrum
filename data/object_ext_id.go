package data

import "github.com/google/uuid"

// ExtendedIdentifier is a service struct that is using for association identifier and name of object like Client and User
type ExtendedIdentifier struct {
	ID   uuid.UUID
	Name string
}
