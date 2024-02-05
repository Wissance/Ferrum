package dto

// OpenIdConfiguration struct that represents OpenId Configuration like in KeyCloak
// todo(UMV): this class would have only those properties (others commented)
type OpenIdConfiguration struct {
	Issuer                             string   `json:"issuer"`
	AuthorizationEndpoint              string   `json:"authorization_endpoint"`
	TokenEndpoint                      string   `json:"token_endpoint"`
	IntrospectionEndpoint              string   `json:"introspection_endpoint"`
	UserInfoEndpoint                   string   `json:"userinfo_endpoint"`
	EndSessionEndpoint                 string   `json:"end_session_endpoint"`
	DeviceAuthorizationEndpoint        string   `json:"device_authorization_endpoint"`
	RegistrationEndpoint               string   `json:"registration_endpoint"`
	PushedAuthorizationRequestEndpoint string   `json:"pushed_authorization_request_endpoint"`
	BachChannelAuthorizationEndpoint   string   `json:"bach_channel_authorization_endpoint"`
	GrantTypesSupported                []string `json:"grant_types_supported"`
	ResponseTypesSupported             []string `json:"response_types_supported"`
	// JwksUri                            string   `json:"jwks_uri"` // TODO (UMV): Uncomment if required
	// FrontChannelLogoutSessionSupported bool         // TODO (UMV): Uncomment if required
	// FrontChannelLogoutSupported bool                // TODO (UMV): Uncomment if required
	// CheckSessionIframe string                       // TODO (UMV): Uncomment if required
	// SubjectTypeSupported []string                   // TODO (UMV): Uncomment if required
	//IdTokenSigningAlgValuesSupported                   []string `json:"id_token_signing_alg_values_supported"`
	//IdTokenEncryptionEncValuesSupported                []string `json:"id_token_encryption_enc_values_supported"`
	//UserInfoSigningAlgValuesSupported                  []string `json:"userinfo_signing_alg_values_supported"`
	//RequestObjectSigningAlgValuesSupported             []string `json:"request_object_signing_alg_values_supported"`
	//RequestEncryptionEncValuesSupported                []string `json:"request_encryption_enc_values_supported"`
	ResponseModesSupported []string `json:"response_modes_supported"`
	//TokenEndpointAuthSigningAlgValuesSupported         []string `json:"token_endpoint_auth_signing_alg_values_supported"`
	//IntrospectionEndpointAuthMethodsSupported          []string `json:"introspection_endpoint_auth_methods_supported"`
	//IntrospectionEndpointAuthSigningAlgValuesSupported []string `json:"introspection_endpoint_auth_signing_alg_values_supported"`
	//AuthorizationSigningAlgValuesSupported             []string `json:"authorization_signing_alg_values_supported"`
	//AuthorizationEncryptionAlgValuesSupported          []string `json:"authorization_encryption_alg_values_supported"`
	//AuthorizationEncryptionEncValuesSupported          []string `json:"authorization_encryption_enc_values_supported"`
	ClaimsSupported                      []string `json:"claims_supported"`
	ClaimTypesSupported                  []string `json:"claim_types_supported"`
	ClaimsParameterSupported             bool     `json:"claims_parameter_supported"`
	RequestParameterSupported            bool     `json:"request_parameter_supported"`
	CodeChallengeMethodsSupported        []string `json:"code_challenge_methods_supported"`
	TlsClientCertificateBoundAccessToken bool     `json:"tls_client_certificate_bound_access_token"`
	//RevocationEndpointAuthMethodsSupported             []string `json:"revocation_endpoint_auth_methods_supported"`
	//RevocationEndpointAuthSigningAlgValuesSupported    []string `json:"revocation_endpoint_auth_signing_alg_values_supported"`
	//BackChannelLogoutSupported                         bool     // TODO (UMV): Uncomment if required
	//BackChannelLogoutSessionSupported                  bool     // TODO (UMV): Uncomment if required
}
