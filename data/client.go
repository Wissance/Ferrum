package data

import (
	"github.com/google/uuid"
)

type ClientType string

const (
	Public       ClientType = "public"
	Confidential            = "confidential"
)

type Client struct {
	Type ClientType
	ID   uuid.UUID
	Name string
	Auth Authentication
}
