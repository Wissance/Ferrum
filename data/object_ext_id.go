package data

import "github.com/google/uuid"

type ExtendedIdentifier struct {
	ID   uuid.UUID
	Name string
}
