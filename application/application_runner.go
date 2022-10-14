package application

import (
	"context"
	"github.com/wissance/Ferrum/logging"
)

type AppContextBase struct {
	Context context.Context
}

type AppRunner interface {
	Start() (bool, error)
	Stop() (bool, error)
	Init() (bool, error)
	GetLogger() *logging.AppLogger
}
