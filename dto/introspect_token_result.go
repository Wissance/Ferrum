package dto

// StringOrArray represents a value that can either be a string or an array of strings
type StringOrArray []string

type IntrospectTokenResult struct {
	Exp      int           `json:"exp,omitempty"`
	Nbf      int           `json:"nbf,omitempty"`
	Iat      int           `json:"iat,omitempty"`
	Aud      StringOrArray `json:"aud,omitempty"`
	Active   bool          `json:"active,omitempty"`
	AuthTime int           `json:"auth_time,omitempty"`
	Jti      string        `json:"jti,omitempty"`
	Type     string        `json:"typ,omitempty"`
}
