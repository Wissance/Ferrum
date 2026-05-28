package dto

type Roles struct {
	Roles []string `json:"roles"`
}

type AccountRoles struct {
	AccountRoles Roles `json:"account"`
}

type IntrospectTokenResult struct {
	Exp            int64        `json:"exp,omitempty"`
	Iss            string       `json:"iss,omitempty"`
	Iat            int64        `json:"iat,omitempty"`
	Active         bool         `json:"active,omitempty"`
	Username       string       `json:"username,omitempty"`
	ClientId       string       `json:"client_id,omitempty"`
	Scope          string       `json:"scope,omitempty"`
	RealmAccess    Roles        `json:"realm_access,omitempty"`
	ResourceAccess AccountRoles `json:"resource_access,omitempty"`
}
