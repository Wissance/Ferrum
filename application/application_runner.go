package application

import (
	"github.com/wissance/Ferrum/logging"
)

/*type AppContextBase struct {
	Context context.Context
}*/

// AppRunner interface that allows to manipulate application
type AppRunner interface {
	// Start this function starts initialized application (must be called after Init)
	Start() (bool, error)
	// Stop function to stop application
	Stop() (bool, error)
	// Init function initializes application components
	Init() (bool, error)
	// GetLogger function that required after app initialized all components to log some additional information about application stop
	GetLogger() *logging.AppLogger
}

var _ AppRunner = (*Application)(nil)
