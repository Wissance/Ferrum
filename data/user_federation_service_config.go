package data

import "crypto/tls"

type UserFederationServiceType string

const (
	LDAP    UserFederationServiceType = "ldap"
	FreeIPA UserFederationServiceType = "freeipa"
)

type UserFederationServiceConfig struct {
	Type UserFederationServiceType `json:"type"`
	// Url is a base url to fetch data
	Url string `json:"url"`
	// Name is internal Unique identifier, MUST be unique across all providers
	Name string `json:"name"`
	// SysUser is a system User, if SysUser == "" then mode IsAnonymous
	SysUser string `json:"sys_user"`
	// SysPassword is a system password
	SysPassword string `json:"sys_password"`
	// TlsCfg is an HTTPS configuration options, use InsecureSkipVerify=true to allow to use self-signed certificate
	TlsCfg *tls.Config `json:"tls_cfg"`
	// EntryPoint is case of LDAP is a catalog where we should fetch data, i.e.
	EntryPoint string `json:"entry_point"`
}

func (u UserFederationServiceConfig) IsAnonymousAccess() bool {
	return len(u.SysUser) == 0
}
