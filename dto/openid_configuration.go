package dto

type OpenIdConfiguration struct {
	Issuer                 string   `json:"issuer"`
	AuthorizationEndpoint  string   `json:"authorization_endpoint"`
	TokenEndpoint          string   `json:"token_endpoint"`
	IntrospectionEndpoint  string   `json:"introspection_endpoint"`
	UserInfoEndpoint       string   `json:"userinfo_endpoint"`
	EndSessionEndpoint     string   `json:"end_session_endpoint"`
	RegistrationEndpoint   string   `json:"registration_endpoint"`
	GrantTypesSupported    []string `json:"grant_types_supported"`
	ResponseTypesSupported []string `json:"response_types_supported"`
}
