package rest

import (
	"Ferrum/managers"
	"Ferrum/services"
)

type WebApiContext struct {
	DataProvider   *managers.DataContext
	Security       *services.SecurityService
	TokenGenerator *services.JwtGenerator
}
