package rest

import (
	"github.com/wissance/Ferrum/managers"
	"github.com/wissance/Ferrum/services"
	"Ferrum/logging"
)

type WebApiContext struct {
	DataProvider   *managers.DataContext
	Security       *services.SecurityService
	TokenGenerator *services.JwtGenerator
	Logger         *logging.AppLogger
}
