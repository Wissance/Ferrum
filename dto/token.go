package dto

type Token struct {
	AccessToken     string `json:"access_token"`
	Expires         int    `json:"expires_in"`
	RefreshExpires  int    `json:"refresh_expires_in"`
	RefreshToken    string `json:"refresh_token"`
	TokenType       string `json:"token_type"`
	NotBeforePolicy int    `json:"not-before-policy"`
	Session         string `json:"session_state"`
	Scope           string `json:"scope"`
}
