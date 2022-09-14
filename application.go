package main

import (
	"Ferrum/api/rest"
	"Ferrum/application"
	"Ferrum/config"
	"Ferrum/managers"
	"Ferrum/services"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	r "github.com/wissance/gwuu/api/rest"
	"github.com/wissance/stringFormatter"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type Application struct {
	appConfigFile  *string
	dataConfigFile *string
	appConfig      *config.AppConfig
	webApiHandler  *r.WebApiHandler
	webApiContext  *rest.WebApiContext
}

func Create(configFile string, dataFile string) application.AppRunner {
	app := &Application{}
	app.appConfigFile = &configFile
	app.dataConfigFile = &dataFile
	appRunner := application.AppRunner(app)
	return appRunner
}

func (app *Application) Start() (bool, error) {
	var err error
	go func() {
		err = app.startWebService()
		if err != nil {
			fmt.Println(stringFormatter.Format("An error occurred during API Service Start"))
		}
	}()
	return err == nil, err
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
	return true, nil
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

	app.appConfig = &config.AppConfig{}
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
	dataProvider := managers.Create(*app.dataConfigFile)
	securityService := services.CreateSecurityService(&dataProvider)
	// todo: provide GOOD key as a file ....
	app.webApiContext = &rest.WebApiContext{DataProvider: &dataProvider, Security: &securityService,
		TokenGenerator: &services.JwtGenerator{SignKey: []byte("secureSecretText")}}
	router := app.webApiHandler.Router
	router.StrictSlash(true)
	app.initKeyCloakSimilarRestApiRoutes(router)
	// Setting up listener for logging
	return nil
}

func (app *Application) initKeyCloakSimilarRestApiRoutes(router *mux.Router) {
	// 1. Generate token endpoint - /auth/realms/{realm}/protocol/openid-connect/token
	app.webApiHandler.HandleFunc(router, "/auth/realms/{realm}/protocol/openid-connect/token/", app.webApiContext.IssueNewToken, http.MethodPost)
	// 2. Get userinfo endpoint - /auth/realms/SOAR/protocol/openid-connect/userinfo
	app.webApiHandler.HandleFunc(router, "/auth/realms/{realm}/protocol/openid-connect/userinfo/", app.webApiContext.GetUserInfo, http.MethodGet)
}

func (app *Application) startWebService() error {
	var err error
	addressTemplate := "{0}:{1}"
	address := stringFormatter.Format(addressTemplate, app.appConfig.ServerCfg.Address, app.appConfig.ServerCfg.Port)
	switch app.appConfig.ServerCfg.Schema { //nolint:exhaustive
	case config.HTTP:
		fmt.Println(stringFormatter.Format("Starting \"HTTP\" WEB API Service on address: \"{0}\"", address))
		err = http.ListenAndServe(address, app.webApiHandler.Router)
		if err != nil {
			fmt.Println(stringFormatter.Format("An error occurred during attempt to start \"HTTP\" WEB API Service: {0}", err.Error()))
		}
	case config.HTTPS:
		fmt.Println(stringFormatter.Format("5. Starting \"HTTPS\" REST API Service on address: \"{0}\"", address))
		cert := app.appConfig.ServerCfg.Security.CertificateFile
		key := app.appConfig.ServerCfg.Security.KeyFile
		err = http.ListenAndServeTLS(address, cert, key, app.webApiHandler.Router)
		if err != nil {
			fmt.Println(stringFormatter.Format("An error occurred during attempt tp start \"HTTPS\" REST API Service: {0}", err.Error()))
		}
	}
	return err
}
