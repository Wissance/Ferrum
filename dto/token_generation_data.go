package dto

type TokenGenerationData struct {
	ClientId     string `json:"client_id" schema:"client_id" form:"client_id"`
	ClientSecret string `json:"client_secret" schema:"client_secret" form:"client_secret"`
	GrantType    string `json:"grant_type" schema:"grant_type" form:"grant_type"`
	Scope        string `json:"scope" schema:"scope" form:"scope"`
	Username     string `json:"username" schema:"username" form:"username"`
	Password     string `json:"password" schema:"password" form:"password"`
	RefreshToken string `json:"refresh_token" schema:"refresh_token" form:"refresh_token"`
}
