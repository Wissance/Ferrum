package operations

type ResourceType string

const (
	RealmResource                ResourceType = "realm"
	ClientResource                            = "client"
	UserResource                              = "user"
	UserFederationConfigResource              = "user_federation"
)

type OperationType string

const (
	GetOperation    OperationType = "get"
	CreateOperation               = "create"
	DeleteOperation               = "delete"
	UpdateOperation               = "update"
	ChangePassword                = "change_password"
	ResetPassword                 = "reset_password"
)
