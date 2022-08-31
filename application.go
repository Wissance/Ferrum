package Ferrum

import (
	"Ferrum/application"
	"Ferrum/config"
)

type Application struct {
	appConfig *config.ServerConfig
}

func Create(configFile string, dataFile string) application.AppRunner {
	app := &Application{}
	appRunner := application.AppRunner(app)
	return appRunner
}

func (app *Application) Start() (bool, error) {
	return false, nil
}

func (app *Application) Init() (bool, error) {
	// init users
	// init webapi
	return false, nil
}

func (app *Application) Stop() (bool, error) {
	return false, nil
}
