package dto

type Token struct {
	AccessToken     string
	Expires         int
	RefreshExpires  int
	RefreshToken    string
	TokenType       string
	NotBeforePolicy int
	Session         string
	Scope           string
}
