package data

type ObjectType string

const (
	REALM                          ObjectType = "realm"
	CLIENT                         ObjectType = "client"
	USER                           ObjectType = "user"
	USER_FEDERATION_SERVICE_CONFIG ObjectType = "user_federation_service_config"
	SERVER_SETTINGS                ObjectType = "server_settings"
)
