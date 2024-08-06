package data

import "crypto/tls"

type UserFederationServiceType string

type UserFederationServiceConfig struct {
	// Name is internal Unique identifier, MUST be unique across all providers
	Name string `json:"name"`
	// Url is a base url to fetch data
	Url string `json:"url"`
	// TlsCfg is an HTTPS configuration options, use InsecureSkipVerify=true to allow to use self-signed certificate
	TlsCfg *tls.Config `json:"tls_cfg"`
	// EntryPoint is case of LDAP is a catalog where we should fetch data, i.e.
	EntryPoint string `json:"entry_point"`
	Password   string `json:"password"`
}
