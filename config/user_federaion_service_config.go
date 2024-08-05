package config

import "crypto/tls"

type UserFederationServiceConfig struct {
	// Name is internal Unique identifier
	Name   string      `json:"name"`
	Url    string      `json:"url"`
	TlsCfg *tls.Config `json:"tls_cfg"`
}
