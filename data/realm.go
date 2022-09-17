package data

type Realm struct {
	Name                   string        `json:"name"`
	Clients                []Client      `json:"clients"`
	Users                  []interface{} `json:"users"`
	TokenExpiration        int           `json:"token_expiration"`
	RefreshTokenExpiration int           `json:"refresh_expiration"`
}
