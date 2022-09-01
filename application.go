package Ferrum

import (
	"Ferrum/application"
	"Ferrum/config"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	r "github.com/wissance/gwuu/api/rest"
	"github.com/wissance/stringFormatter"
	"io/ioutil"
	"path/filepath"
)

type Application struct {
	appConfigFile *string
	appConfig     *config.ServerConfig
	webApiHandler *r.WebApiHandler
}

func Create(configFile string, dataFile string) application.AppRunner {
	app := &Application{}
	app.appConfigFile = &configFile
	appRunner := application.AppRunner(app)
	return appRunner
}

func (app *Application) Start() (bool, error) {
	return false, nil
}

func (app *Application) Init() (bool, error) {
	err := app.readAppConfig()
	if err != nil {
		fmt.Println(stringFormatter.Format("An error occurred during reading app config file: {0}", err.Error()))
		return false, err
	}
	// init users, today we are reading data file
	err = app.initDataProviders()
	if err != nil {
		fmt.Println(stringFormatter.Format("An error occurred during data providers init: {0}", err.Error()))
		return false, err
	}
	// init webapi
	err = app.initRestApi()
	if err != nil {
		fmt.Println(stringFormatter.Format("An error occurred during rest api init: {0}", err.Error()))
		return false, err
	}
	return true, nil
}

func (app *Application) Stop() (bool, error) {
	return false, nil
}

func (app *Application) readAppConfig() error {
	absPath, err := filepath.Abs(*app.appConfigFile)
	if err != nil {
		fmt.Println(stringFormatter.Format("An error occurred during getting config file abs path: {0}", err.Error()))
		return err
	}

	fileData, err := ioutil.ReadFile(absPath)
	if err != nil {
		fmt.Println(stringFormatter.Format("An error occurred during config file reading: {0}", err.Error()))
		return err
	}

	if err = json.Unmarshal(fileData, app.appConfig); err != nil {
		fmt.Println(stringFormatter.Format("An error occurred during config file unmarshal: {0}", err.Error()))
		return err
	}

	return nil
}

func (app *Application) initDataProviders() error {
	return nil
}

func (app *Application) initRestApi() error {
	app.webApiHandler = r.NewWebApiHandler(true, r.AnyOrigin)
	router := app.webApiHandler.Router
	router.StrictSlash(true)
	app.initKeuCloakSimilarRestApiRoutes(router)
	return nil
}

func (app *Application) initKeuCloakSimilarRestApiRoutes(router *mux.Router) {
	// 1. Generate token endpoint - /auth/realms/{realm}/protocol/openid-connect/token
	// 2. Get userinfo endpoint - /auth/realms/SOAR/protocol/openid-connect/userinfo
}
