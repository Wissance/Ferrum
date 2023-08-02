package dto

type TokenGenerationData struct {
	ClientId     string `json:"client_id" schema:"client_id"`
	ClientSecret string `json:"client_secret" schema:"client_secret"`
	GrantType    string `json:"grant_type" schema:"grant_type"`
	Scope        string `json:"scope" schema:"scope"`
	Username     string `json:"username" schema:"username"`
	Password     string `json:"password" schema:"password"`
	RefreshToken string `json:"refresh_token" schema:"refresh_token"`
}
