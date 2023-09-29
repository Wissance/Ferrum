package data

type AuthenticationType int

// ClientIdAndSecrets AuthenticationType represents Confidential Clients
const (
	ClientIdAndSecrets AuthenticationType = 1
)

// Authentication struct for Clients authentication data, for ClientIdAndSecrets Value stores ClientSecret
type Authentication struct {
	Type       AuthenticationType
	Value      string
	Attributes interface{}
}
