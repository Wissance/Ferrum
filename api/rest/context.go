package rest

import (
	"Ferrum/logging"
	"Ferrum/managers"
	"Ferrum/services"
)

type WebApiContext struct {
	DataProvider   *managers.DataContext
	Security       *services.SecurityService
	TokenGenerator *services.JwtGenerator
	Logger         *logging.AppLogger
}
