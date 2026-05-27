package rest

import (
	"context"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/Ferrum/managers"
	"github.com/wissance/Ferrum/security/bruteforce"
	"github.com/wissance/Ferrum/services"
)

// WebApiContext is a central Application logic processor manages from Web via HTTP/HTTPS
type WebApiContext struct {
	Address              string
	Schema               string
	ctx                  context.Context
	DataProvider         *managers.DataContext
	AuthDefs             *data.AuthenticationDefs
	Security             *services.SecurityService
	BruteforceProtection *bruteforce.ProtectionService
	TokenGenerator       *services.JwtGenerator
	Logger               *logging.AppLogger
}
