package rest

import (
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/Ferrum/managers"
	"github.com/wissance/Ferrum/services"
)

// WebApiContext is a central Application logic processor manages from Web via HTTP/HTTPS
type WebApiContext struct {
	Address        string
	Schema         string
	DataProvider   *managers.DataContext
	AuthDefs       *data.AuthenticationDefs
	Security       *services.SecurityService
	TokenGenerator *services.JwtGenerator
	Logger         *logging.AppLogger
}
