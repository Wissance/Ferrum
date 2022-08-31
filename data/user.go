package data

import "github.com/google/uuid"

type User interface {
	GetUsername() string
	GetPassword() string
	GetId() uuid.UUID
}
