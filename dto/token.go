package dto

type Token struct {
	AccessToken     string `json:"access_token"`
	Expires         int
	RefreshExpires  int
	RefreshToken    string
	TokenType       string
	NotBeforePolicy int
	Session         string
	Scope           string
}
