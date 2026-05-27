package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/wissance/Ferrum/api/rest"
	"github.com/wissance/Ferrum/api/rest/filter"
	"github.com/wissance/Ferrum/api/rest/metrics"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	appErrs "github.com/wissance/Ferrum/errors"
	"github.com/wissance/Ferrum/globals"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/Ferrum/managers"
	"github.com/wissance/Ferrum/security/bruteforce"
	"github.com/wissance/Ferrum/services"
	"github.com/wissance/Ferrum/sre"
	"github.com/wissance/Ferrum/swagger"
	"github.com/wissance/Ferrum/utils/encoding"
	"github.com/wissance/Ferrum/utils/uuidtools"
	r "github.com/wissance/gwuu/api/rest"
	"github.com/wissance/stringFormatter"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const ferrumSwaggerAddressEnvVariable = "FERRUM_SWAGGER_EXT_ADDRESS"

type Application struct {
	devMode                    bool
	appConfigFile              *string
	dataConfigFile             *string
	secretKeyFile              *string
	appConfig                  *config.AppConfig
	bruteforceProtectionConfig *bruteforce.ProtectionServiceConfig // TODO(UMV): add this config to the app config
	authenticationDefs         *data.AuthenticationDefs
	secretKey                  []byte
	serverData                 *data.ServerData
	ctx                        context.Context
	cancelFunc                 context.CancelFunc
	dataProvider               *managers.DataContext
	webApiHandler              *r.GinBasedWebApiHandler
	webApiContext              *rest.WebApiContext
	logger                     *logging.AppLogger
	httpHandler                *http.Handler
	httpServer                 *http.Server
	shutdownTimeout            time.Duration
	metricsCollector           *sre.MetricsCollector
}

// CreateAppWithConfigs creates but not Init new Application as AppRunner
/* This function creates new Application and pass configFile to newly created object
 * Parameters:
 *     - configFile - path to config
 *     - devMode - developer mode (for showing some non production info = swagger, ...)
 * Returns: new Application as AppRunner
 */
func CreateAppWithConfigs(configFile string, devMode bool, ctx context.Context) AppRunner {
	app := &Application{}
	contextWithCancel, cancelFunc := context.WithCancel(ctx)
	app.devMode = devMode
	app.appConfigFile = &configFile
	app.ctx = contextWithCancel
	app.cancelFunc = cancelFunc
	app.authenticationDefs = &data.AuthenticationDefs{}
	app.metricsCollector = sre.CreateMetricsCollector()
	appRunner := AppRunner(app)
	return appRunner
}

// CreateAppWithData creates but not Init new Application as AppRunner
/* This function creates new Application and pass already decoded json of appConfig and serverData plus secretKey
 * Parameters:
 *     - appConfig  - decoded config
 *     - serverData - decoded server config
 *     - secretKey  - secret key that is using for signing JWT
 * Returns: new Application as AppRunner
 */
func CreateAppWithData(appConfig *config.AppConfig, serverData *data.ServerData, ctx context.Context,
	secretKey []byte, devMode bool) AppRunner {
	contextWithCancel, cancelFunc := context.WithCancel(ctx)
	app := &Application{appConfig: appConfig, secretKey: secretKey, serverData: serverData,
		ctx: contextWithCancel, cancelFunc: cancelFunc, devMode: devMode,
		metricsCollector: sre.CreateMetricsCollector()}
	app.authenticationDefs = &data.AuthenticationDefs{}
	appRunner := AppRunner(app)
	return appRunner
}

// Start function that starts application
/* This function must be called after Init it starts application web server either on HTTP or HTTPS (depends on config Schema value)
 * Parameters: no
 * Return start result (true if Start was successful) and error (nil if start was successful)
 */
func (app *Application) Start() (bool, error) {
	err := app.startWebService()
	if err != nil {
		app.logger.Error(stringFormatter.Format("An error occurred during API Service Start"))
	}
	return err == nil, err
}

// Init initializes application
/* This function implements application subsystem init:
 *    1. Read config if Application was not Created via CreateAppWithData
 *    2. After config was decoded this function implements following initialization
 *       2.1 Read secret file for signing JWT
 *       2.2 Initializes logger
 *       2.3 Initializes Data Provider
 *       2.4 Initializes REST API
 * Parameters: no
 * Return result of init (true if init was successful) and error (nil if init was successful)
 */
func (app *Application) Init() (bool, error) {
	// part that initializes app from configs
	isAppCreatedWithConfigs := app.appConfigFile != nil
	if isAppCreatedWithConfigs {
		// init logger
		cfg, err := config.ReadAppConfig(*app.appConfigFile)
		if err != nil {
			return false, fmt.Errorf("an error occurred during reading app config file: %w", err)
		}
		app.appConfig = cfg
		app.logger = logging.CreateLogger(&app.appConfig.Logging)
		app.logger.Init()

		// after config read init secretKey file and data file (if provider.type == FILE)
		app.secretKeyFile = &app.appConfig.ServerCfg.SecretFile
		if app.appConfig.DataSource.Type == config.FILE {
			app.dataConfigFile = &app.appConfig.DataSource.Source
		}
		// reading secrets key
		key := app.readKey()
		if key == nil {
			app.logger.Error("An error occurred during reading secret key, key is nil")
			return false, errors.New("secret key is nil")
		}
		app.secretKey = key
		app.shutdownTimeout = cfg.ServerCfg.ShutdownTimeout * time.Second
	} else {
		// init logger
		app.logger = logging.CreateLogger(&app.appConfig.Logging)
		app.shutdownTimeout = 30 * time.Second
		app.logger.Init()
	}
	// common part: both configs and direct struct pass
	// init users, today we are reading data file
	err := app.initDataProviders()
	if err != nil {
		app.logger.Error(stringFormatter.Format("An error occurred during data providers init: {0}", err.Error()))
		return false, err
	}
	// todo(umv): add init data,

	// init auth defs
	app.initAuthServerDefs()

	// init webapi
	err = app.initRestApi()
	if err != nil {
		app.logger.Error(stringFormatter.Format("An error occurred during rest api init: {0}", err.Error()))
		return false, err
	}

	app.httpServer = &http.Server{Handler: *app.httpHandler}
	return true, nil
}

// Stop function that stops application
/* Now doesn't do anything, just a stub
 * Parameters : no
 * Returns result of app stop and error
 */
func (app *Application) Stop() (bool, error) {
	ctxWithTimeout, cancel := context.WithTimeout(app.ctx, app.shutdownTimeout)
	defer cancel()
	app.metricsCollector.UnRegisterAllMetrics()
	err := app.httpServer.Shutdown(ctxWithTimeout)
	app.cancelFunc()
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetLogger function that returns logger from initialized application
/* Returns logger from application (app.Logger), must be called after Init
 * Parameters : no
 * Returns: logger
 */
func (app *Application) GetLogger() *logging.AppLogger {
	return app.logger
}

func (app *Application) initDataProviders() error {
	var err error
	if app.serverData != nil {
		dataProvider, prepareErr := managers.PrepareContextUsingData(&app.appConfig.DataSource, app.serverData, app.logger)
		observableProvider := sre.CreateObservableDataContext(app.metricsCollector, &dataProvider)
		app.dataProvider = &observableProvider
		return prepareErr
	}

	if app.dataConfigFile != nil {
		dataProvider, prepareErr := managers.PrepareContextUsingFile(&app.appConfig.DataSource, app.dataConfigFile, app.logger)
		observableProvider := sre.CreateObservableDataContext(app.metricsCollector, &dataProvider)
		app.dataProvider = &observableProvider
		err = prepareErr
	} else {
		dataProvider, prepareErr := managers.PrepareContext(&app.appConfig.DataSource, app.logger)
		observableProvider := sre.CreateObservableDataContext(app.metricsCollector, &dataProvider)
		app.dataProvider = &observableProvider
		err = prepareErr
		if err != nil {
			return err
		}
		return app.initData(dataProvider)
	}
	return err
}

// iniData inits ServerSetting before AppStart
// we are passing here dataProvider because data provider from app does not have initialized SRE yet
func (app *Application) initData(dataProvider managers.DataContext) error {
	// this function init some required data after managers.DataContext creation
	settings, err := dataProvider.GetServerSettings()
	if err == nil && !uuidtools.IsUUIDEmpty(&settings.Admin.Id) {
		// nothing to be done here, settings already exists
		return nil
	}
	if errors.As(err, &appErrs.EmptyNotFoundErr) {
		// ServerSetting were not found, create
		// this salt is going to be used when you are going to change admin pass
		salt := encoding.GenerateRandomSalt()
		passwordEncoder := encoding.NewPasswordJsonEncoder(salt)
		settings = &data.ServerSettings{
			AllowedHosts: app.appConfig.Security.AllowedHosts,
			Admin: data.AdminUser{
				Id:           uuid.New(),
				Username:     app.appConfig.Security.Admin.Username,
				PasswordSalt: salt,
				PasswordHash: passwordEncoder.GetB64PasswordHash(app.appConfig.Security.Admin.Password),
			},
			AdminApiUrlPrefix: app.appConfig.Security.AdminApiUrlPrefix,
		}

		err = dataProvider.SetServerSettings(settings)
	}
	return err
}

func (app *Application) initRestApi() error {
	app.webApiHandler = r.NewGinBasedWebApiHandler(true, r.AnyOrigin)
	securityService := services.CreateSecurityService(app.dataProvider, app.logger, app.ctx)
	bruteforceProtectionConfig := bruteforce.ProtectionServiceConfig{
		WatchTimeSec: 600,
	}
	app.bruteforceProtectionConfig = &bruteforceProtectionConfig
	bruteForceProtectionService := bruteforce.CreateProtectionService(app.ctx, app.bruteforceProtectionConfig, app.logger)
	serverAddress := stringFormatter.Format("{0}:{1}", app.appConfig.ServerCfg.Address, app.appConfig.ServerCfg.Port)
	app.webApiContext = &rest.WebApiContext{
		Address: serverAddress, Schema: string(app.appConfig.ServerCfg.Schema),
		AuthDefs:     app.authenticationDefs,
		DataProvider: app.dataProvider, Security: &securityService,
		TokenGenerator:       &services.JwtGenerator{SignKey: app.secretKey, Logger: app.logger},
		Logger:               app.logger,
		BruteforceProtection: &bruteForceProtectionService,
	}
	router := app.webApiHandler.Router
	router.RedirectTrailingSlash = true
	router.Use(app.metricsCollector.HttpMetricsCollectMiddleware())
	router.Use(filter.AttackersFilterMiddleware(app.webApiContext.BruteforceProtection, app.logger))
	rootRoutesGroup := router.Group("/")
	app.initKeyCloakSimilarRestApiRoutes(rootRoutesGroup)
	app.initSRERestApiRoutes(rootRoutesGroup)
	if app.devMode {
		app.initSwaggerRoutes(rootRoutesGroup)
	}
	routerHandler := router.Handler()
	app.httpHandler = &routerHandler
	// Setting up listener for logging
	appenderIndex := app.logger.GetAppenderIndex(config.RollingFile, app.appConfig.Logging.Appenders)
	if appenderIndex == -1 {
		app.logger.Info("The RollingFile appender was not found.")
		var resultRouter http.Handler = router
		app.httpHandler = &resultRouter
		return nil
	}
	app.createHttpLoggingHandler(appenderIndex)
	return nil
}

func (app *Application) initSwaggerRoutes(router *gin.RouterGroup) {
	// Swagger docs router and config
	swagger.SwaggerInfo.Version = "v0.9.3"
	swagger.SwaggerInfo.Title = "Ferrum Authorization Server"
	swagger.SwaggerInfo.Description = "Ferrum a better Authorization server compatible by API with a KeyCloak"
	address := app.getSwaggerAddress()
	if address == "" {
		address = app.appConfig.ServerCfg.Address
	}
	swagger.SwaggerInfo.Host = stringFormatter.Format("{0}:{1}", address, app.appConfig.ServerCfg.Port)

	app.webApiHandler.GET(router, "/swagger", func(ctx *gin.Context) {
		httpSwagger.Handler()
	})
}

func (app *Application) initSRERestApiRoutes(router *gin.RouterGroup) {
	app.webApiHandler.GET(router, "/metrics", metrics.GetPrometheusHandler())
}

func (app *Application) initAuthServerDefs() {
	app.authenticationDefs.SupportedGrantTypes = []string{
		globals.AuthorizationTokenGrantType,
		globals.RefreshTokenGrantType,
		globals.PasswordGrantType,
	}

	app.authenticationDefs.SupportedResponseTypes = []string{
		globals.TokenResponseType,
		globals.CodeResponseType,
		globals.CodeTokenResponseType,
	}

	app.authenticationDefs.SupportedResponses = []string{
		globals.JwtResponse,
	}

	app.authenticationDefs.SupportedScopes = []string{
		globals.ProfileEmailScope,
		globals.OpenIdScope,
		globals.ProfileScope,
		globals.EmailScope,
	}

	app.authenticationDefs.SupportedClaimTypes = []string{
		"normal",
	}

	app.authenticationDefs.SupportedClaims = []string{
		globals.SubClaimType,
		globals.EmailClaimType,
		globals.PreferredUsernameClaim,
	}
}

func (app *Application) initKeyCloakSimilarRestApiRoutes(router *gin.RouterGroup) {
	// 1. Introspect endpoint - /auth/realms/:realm/protocol/openid-connect/introspect
	app.webApiHandler.POST(router, "/auth/realms/:realm/protocol/openid-connect/token/introspect", app.webApiContext.Introspect)
	app.webApiHandler.POST(router, "/realms/:realm/protocol/openid-connect/token/introspect", app.webApiContext.Introspect)
	// 2. Generate token endpoint - /auth/realms/:realm/protocol/openid-connect/token
	app.webApiHandler.POST(router, "/auth/realms/:realm/protocol/openid-connect/token", app.webApiContext.IssueNewToken)
	app.webApiHandler.POST(router, "/realms/:realm/protocol/openid-connect/token", app.webApiContext.IssueNewToken)
	// 3. Get userinfo endpoint - /auth/realms/:realm/protocol/openid-connect/userinfo
	app.webApiHandler.GET(router, "/auth/realms/:realm/protocol/openid-connect/userinfo", app.webApiContext.GetUserInfo)
	app.webApiHandler.GET(router, "/realms/:realm/protocol/openid-connect/userinfo", app.webApiContext.GetUserInfo)
	// 4. OpenId Configuration endpoint
	app.webApiHandler.GET(router, "/auth/realms/:realm/.well-known/openid-configuration", app.webApiContext.GetOpenIdConfiguration)
	app.webApiHandler.GET(router, "/realms/:realm/.well-known/openid-configuration", app.webApiContext.GetOpenIdConfiguration)
}

func (app *Application) startWebService() error {
	var err error
	addressTemplate := "{0}:{1}"
	address := stringFormatter.Format(addressTemplate, app.appConfig.ServerCfg.Address, app.appConfig.ServerCfg.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	app.httpServer.Addr = address
	switch app.appConfig.ServerCfg.Schema { //nolint:exhaustive
	case config.HTTP:
		app.logger.Info(stringFormatter.Format("Starting \"HTTP\" WEB API Service on address: \"{0}\"", address))
		go func() {
			err = app.httpServer.Serve(listener)
			if err != nil {
				app.logger.Error(
					stringFormatter.Format("An error occurred during attempt to start \"HTTP\" WEB API Service: {0}", err.Error()))
			}
		}()
	case config.HTTPS:
		app.logger.Info(stringFormatter.Format("Starting \"HTTPS\" REST API Service on address: \"{0}\"", address))
		cert := app.appConfig.ServerCfg.Security.CertificateFile
		key := app.appConfig.ServerCfg.Security.KeyFile
		go func() {
			err = app.httpServer.ServeTLS(listener, cert, key)
			if err != nil {
				app.logger.Error(
					stringFormatter.Format("An error occurred during attempt tp start \"HTTPS\" REST API Service: {0}", err.Error()))
			}
		}()
	}
	return err
}

func (app *Application) readKey() []byte {
	absPath, err := filepath.Abs(*app.appConfigFile)
	if err != nil {
		app.logger.Error(stringFormatter.Format("An error occurred during getting key file abs path: {0}", err.Error()))
		return nil
	}

	fileData, err := os.ReadFile(absPath)
	if err != nil {
		app.logger.Error(stringFormatter.Format("An error occurred during key file reading: {0}", err.Error()))
		return nil
	}

	return fileData
}

func (app *Application) createHttpLoggingHandler(index int) {
	destination := app.appConfig.Logging.Appenders[index].Destination

	lumberjackWriter := lumberjack.Logger{
		Filename:   string(destination.File),
		MaxSize:    destination.MaxSize,
		MaxAge:     destination.MaxAge,
		MaxBackups: destination.MaxBackups,
		LocalTime:  destination.LocalTime,
		Compress:   false,
	}

	if app.appConfig.Logging.LogHTTP {
		if app.appConfig.Logging.ConsoleOutHTTP {
			writer := io.MultiWriter(&lumberjackWriter, os.Stdout)
			gin.DefaultWriter = writer
			gin.DefaultErrorWriter = writer
		} else {
			gin.DefaultWriter = &lumberjackWriter
			gin.DefaultErrorWriter = &lumberjackWriter
		}
	}
}

func (app *Application) getSwaggerAddress() string {
	// 1. Get ENV Variable - FERRUM_SWAGGER_EXT_ADDRESS (see .env file)
	envAddr := os.Getenv(ferrumSwaggerAddressEnvVariable)
	if len(envAddr) > 0 {
		return envAddr
	}

	// 2. Get Address from Network Interfaces
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addresses {
		// check the address type and if it is not a loop back the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
