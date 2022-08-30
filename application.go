package Ferrum

import "Ferrum/application"

type Application struct {
}

func Create() application.AppRunner {
	app := &Application{}
	appRunner := application.AppRunner(app)
	return appRunner
}

func (app *Application) Start() (bool, error) {
	return false, nil
}

func (app *Application) Init() (bool, error) {
	return false, nil
}

func (app *Application) Stop() (bool, error) {
	return false, nil
}
