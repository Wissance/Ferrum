package ferrum

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wissance/Ferrum/api/rest"
	"github.com/wissance/Ferrum/application"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/managers"
	"github.com/wissance/Ferrum/services"
	r "github.com/wissance/gwuu/api/rest"
	"github.com/wissance/stringFormatter"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type Application struct {
	appConfigFile  *string
	dataConfigFile *string
	secretKeyFile  *string
	appConfig      *config.AppConfig
	secretKey      *[]byte
	serverData     *data.ServerData
	dataProvider   *managers.DataContext
	webApiHandler  *r.WebApiHandler
	webApiContext  *rest.WebApiContext
	Logger         *logging.AppLogger
}

func CreateAppWithConfigs(configFile string, dataFile string, secretKeyFile string) application.AppRunner {
	app := &Application{}
	app.appConfigFile = &configFile
	app.dataConfigFile = &dataFile
	app.secretKeyFile = &secretKeyFile
	appRunner := application.AppRunner(app)
	return appRunner
}

func CreateAppWithData(appConfig *config.AppConfig, serverData *data.ServerData, secretKey []byte) application.AppRunner {
	app := &Application{appConfig: appConfig, secretKey: &secretKey, serverData: serverData}
	appRunner := application.AppRunner(app)
	return appRunner
}

func (app *Application) Start() (bool, error) {
	var err error
	go func() {
		err = app.startWebService()
		if err != nil {
			app.Logger.Error(stringFormatter.Format("An error occurred during API Service Start"))
		}
	}()
	return err == nil, err
}

func (app *Application) Init() (bool, error) {
	// part that initializes app from configs
	if app.appConfigFile != nil {
		err := app.readAppConfig()
		if err != nil {
			fmt.Println(stringFormatter.Format("An error occurred during reading app config file: {0}", err.Error()))
			return false, err
		}

		// reading secrets key
		key := app.readKey()
		if key == nil {
			fmt.Println(stringFormatter.Format("An error occurred during reading secret key, key is nil"))
			return false, errors.New("secret key is nil")
		}
		app.secretKey = key
	}
	// common part: both configs and direct struct pass
	// init logger
	app.Logger = logging.CreateLogger(&app.appConfig.Logging)
	app.Logger.Init()
	// init users, today we are reading data file
	err := app.initDataProviders()
	if err != nil {
		app.Logger.Error(stringFormatter.Format("An error occurred during data providers init: {0}", err.Error()))
		return false, err
	}

	// init webapi
	err = app.initRestApi()
	if err != nil {
		app.Logger.Error(stringFormatter.Format("An error occurred during rest api init: {0}", err.Error()))
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
		app.Logger.Error(stringFormatter.Format("An error occurred during getting config file abs path: {0}", err.Error()))
		return err
	}

	fileData, err := ioutil.ReadFile(absPath)
	if err != nil {
		app.Logger.Error(stringFormatter.Format("An error occurred during config file reading: {0}", err.Error()))
		return err
	}

	app.appConfig = &config.AppConfig{}
	if err = json.Unmarshal(fileData, app.appConfig); err != nil {
		app.Logger.Error(stringFormatter.Format("An error occurred during config file unmarshal: {0}", err.Error()))
		return err
	}

	return nil
}

func (app *Application) initDataProviders() error {
	if app.dataConfigFile != nil {
		dataProvider := managers.CreateAndContextInitWithDataFile(*app.dataConfigFile)
		app.dataProvider = &dataProvider
	} else {
		dataProvider := managers.CreateAndContextInitUsingData(app.serverData)
		app.dataProvider = &dataProvider
	}
	return nil
}

func (app *Application) initRestApi() error {
	app.webApiHandler = r.NewWebApiHandler(true, r.AnyOrigin)
	securityService := services.CreateSecurityService(app.dataProvider)
	app.webApiContext = &rest.WebApiContext{DataProvider: app.dataProvider, Security: &securityService,
		TokenGenerator: &services.JwtGenerator{SignKey: *app.secretKey}}
	router := app.webApiHandler.Router
	router.StrictSlash(true)
	app.initKeyCloakSimilarRestApiRoutes(router)
	// Setting up listener for logging
	return nil
}

func (app *Application) initKeyCloakSimilarRestApiRoutes(router *mux.Router) {
	app.webApiHandler.HandleFunc(router, "/auth/realms/{realm}/protocol/openid-connect/token/introspect/", app.webApiContext.Introspect, http.MethodPost)
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
		app.Logger.Info(stringFormatter.Format("Starting \"HTTP\" WEB API Service on address: \"{0}\"", address))
		err = http.ListenAndServe(address, app.webApiHandler.Router)
		if err != nil {
			app.Logger.Error(stringFormatter.Format("An error occurred during attempt to start \"HTTP\" WEB API Service: {0}", err.Error()))
		}
	case config.HTTPS:
		app.Logger.Info(stringFormatter.Format("Starting \"HTTPS\" REST API Service on address: \"{0}\"", address))
		cert := app.appConfig.ServerCfg.Security.CertificateFile
		key := app.appConfig.ServerCfg.Security.KeyFile
		err = http.ListenAndServeTLS(address, cert, key, app.webApiHandler.Router)
		if err != nil {
			app.Logger.Error(stringFormatter.Format("An error occurred during attempt tp start \"HTTPS\" REST API Service: {0}", err.Error()))
		}
	}
	return err
}

func (app *Application) readKey() *[]byte {
	absPath, err := filepath.Abs(*app.appConfigFile)
	if err != nil {
		app.Logger.Error(stringFormatter.Format("An error occurred during getting key file abs path: {0}", err.Error()))
		return nil
	}

	fileData, err := ioutil.ReadFile(absPath)
	if err != nil {
		app.Logger.Error(stringFormatter.Format("An error occurred during key file reading: {0}", err.Error()))
		return nil
	}

	return &fileData
}
