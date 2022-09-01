package data

type AuthenticationType int

const (
	ClientIdAndSecrets AuthenticationType = 1
)

type Authentication struct {
	Type       AuthenticationType
	Value      string
	Attributes interface{}
}
