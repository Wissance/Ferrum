package rest

import (
	"github.com/wissance/Ferrum/managers"
	"github.com/wissance/Ferrum/services"
)

type WebApiContext struct {
	DataProvider   *managers.DataContext
	Security       *services.SecurityService
	TokenGenerator *services.JwtGenerator
}
