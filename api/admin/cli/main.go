package main

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/google/uuid"
	appErrs "github.com/wissance/Ferrum/errors"
	"github.com/wissance/Ferrum/utils/encoding"
	"log"
	"os"

	"github.com/wissance/Ferrum/managers"

	"github.com/wissance/Ferrum/api/admin/cli/operations"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	sf "github.com/wissance/stringFormatter"
)

const defaultConfig = "./config_w_redis.json"

var (
	argConfigFile = flag.String("config", defaultConfig, "Application config for working with a persistent data store")
	argOperation  = flag.String("operation", "", "One of the available operations read|create|update|delete or user specific change/reset password")
	argResource   = flag.String("resource", "", "\"realm\", \"client\" or \"user\" or maybe other in future")
	argResourceId = flag.String("resource_id", "", "resource object identifier, id required for the update|delete or read operation")
	argParams     = flag.String("params", "", "Name of a realm for operations on client or user resources")
	argValue      = flag.String("value", "", "Json encoded resource itself")
)

func main() {
	flag.Parse()
	// TODO(UMV): extend config
	cfg, err := config.ReadAppConfig(*argConfigFile)
	if err != nil {
		log.Fatalf("readAppConfig failed: %s", err)
	}
	logger := logging.CreateLogger(&cfg.Logging)
	manager, err := managers.PrepareContext(&cfg.DataSource, logger)
	if err != nil {
		log.Fatalf("prepareContext failed: %s", err)
	}

	operation := operations.OperationType(*argOperation)
	resource := operations.ResourceType(*argResource)
	resourceId := *argResourceId
	params := *argParams
	value := []byte(*argValue)

	isInvalidOperation := operation != operations.GetOperation && operation != operations.CreateOperation &&
		operation != operations.DeleteOperation && operation != operations.UpdateOperation &&
		operation != operations.ChangePassword && operation != operations.ResetPassword
	if isInvalidOperation {
		log.Fatalf("bad Operation \"%s\"", operation)
	}
	// If there is a password change or password collection, it is not necessary to specify Resource
	if !(operation == operations.ChangePassword || operation == operations.ResetPassword) {
		isInvalidResource := resource != operations.RealmResource && resource != operations.ClientResource &&
			resource != operations.UserResource && resource != operations.ServerSettings
		if isInvalidResource {
			msg := sf.Format("Non supported resource \"{0}\"", resource)
			//nolint:govet,staticcheck
			log.Println(msg)
			os.Exit(-1)
		}
	}
	if (resource == operations.ClientResource) || (resource == operations.UserResource) {
		if params == "" {
			log.Fatalf("Not specified Params")
		}
	}

	switch operation {
	case operations.GetOperation:
		if resourceId == "" {
			log.Fatalf("Not specified ResourceId")
		}
		switch resource {
		case operations.ClientResource:
			client, err := manager.GetClient(params, resourceId)
			if err != nil {
				log.Fatalf("GetClient failed: %s", err)
			}
			fmt.Println(*client)

		case operations.UserResource:
			user, err := manager.GetUser(params, resourceId)
			if err != nil {
				log.Fatalf("GetUser failed: %s", err)
			}
			fmt.Println(user.GetUserInfo())

		case operations.RealmResource:
			realm, err := manager.GetRealm(resourceId)
			if err != nil {
				log.Fatalf("GetRealm failed: %s", err)
			}
			fmt.Println(*realm)

		case operations.UserFederationConfigResource:
			userFederation, err := manager.GetUserFederationConfig(params, resourceId)
			if err != nil {
				log.Fatalf("GetUserFederationConfig failed: %s", err)
			}
			fmt.Println(*userFederation)
		}

		return
	case operations.CreateOperation:
		if len(value) == 0 {
			log.Fatalf("Not specified Value")
		}
		switch resource {
		case operations.ClientResource:
			var clientNew data.Client
			if unmarshalErr := json.Unmarshal(value, &clientNew); unmarshalErr != nil {
				//nolint:govet,staticcheck
				log.Fatal(sf.Format("json.Unmarshal failed: {0}", unmarshalErr.Error()))
			}
			if createErr := manager.CreateClient(params, clientNew); createErr != nil {
				//nolint:govet,staticcheck
				log.Fatal(sf.Format("CreateClient failed: {0}", createErr.Error()))
			}
			log.Print(sf.Format("Client: \"{0}\" successfully created", clientNew.Name))

		case operations.UserResource:
			var userNew any
			if err := json.Unmarshal(value, &userNew); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			realm, err := manager.GetRealm(params)
			if err != nil {
				log.Fatalf("GetRealm failed: %s", err)
			}
			user := data.CreateUser(userNew, realm.Encoder)
			if err := manager.CreateUser(params, user); err != nil {
				log.Fatalf("CreateUser failed: %s", err)
			}
			fmt.Println(sf.Format("User: \"{0}\" successfully created", user.GetUsername()))

		case operations.RealmResource:
			var newRealm data.Realm
			if err := json.Unmarshal(value, &newRealm); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			if err := manager.CreateRealm(newRealm); err != nil {
				log.Fatalf("CreateRealm failed: %s", err)
			}
			fmt.Println(sf.Format("Realm: \"{0}\" successfully created", newRealm.Name))
		case operations.UserFederationConfigResource:
			var userFederationConfig data.UserFederationServiceConfig
			if err := json.Unmarshal(value, &userFederationConfig); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			if err := manager.CreateUserFederationConfig(params, userFederationConfig); err != nil {
				log.Fatalf("CreateUserFederationConfig failed: %s", err)
			}
			fmt.Println(sf.Format("User federation service config: \"{0}\" successfully created", userFederationConfig.Name))
		}

		return
	case operations.DeleteOperation:
		if resourceId == "" {
			log.Fatalf("Not specified Resource_id")
		}
		switch resource {
		case operations.ClientResource:
			if err := manager.DeleteClient(params, resourceId); err != nil {
				log.Fatalf("DeleteClient failed: %s", err)
			}
			fmt.Println(sf.Format("Client: \"{0}\" successfully deleted", resourceId))

		case operations.UserResource:
			if err := manager.DeleteUser(params, resourceId); err != nil {
				log.Fatalf("DeleteUser failed: %s", err)
			}
			fmt.Println(sf.Format("User: \"{0}\" successfully deleted", resourceId))

		case operations.RealmResource:
			if err := manager.DeleteRealm(resourceId); err != nil {
				log.Fatalf("DeleteRealm failed: %s", err)
			}
			fmt.Println(sf.Format("Realm: \"{0}\" successfully deleted", resourceId))

		case operations.UserFederationConfigResource:
			if err := manager.DeleteUserFederationConfig(params, resourceId); err != nil {
				log.Fatalf("DeleteUserFederationConfig failed: %s", err)
			}
			fmt.Println(sf.Format("User federation service config: \"{0}\" successfully deleted", resourceId))
		}

		return
	case operations.UpdateOperation:
		if resourceId == "" && resource != operations.ServerSettings {
			log.Fatalf("Not specified Resource_id")
		}
		if len(value) == 0 {
			log.Fatalf("Not specified Value")
		}
		switch resource {
		case operations.ClientResource:
			var newClient data.Client
			if err := json.Unmarshal(value, &newClient); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			if err := manager.UpdateClient(params, resourceId, newClient); err != nil {
				log.Fatalf("UpdateClient failed: %s", err)
			}
			fmt.Println(sf.Format("Client: \"{0}\" successfully updated", newClient.Name))

		case operations.UserResource:
			var newUser any
			if err := json.Unmarshal(value, &newUser); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			user := data.CreateUser(newUser, nil)
			if err := manager.UpdateUser(params, resourceId, user); err != nil {
				log.Fatalf("UpdateUser failed: %s", err)
			}
			fmt.Println(sf.Format("User: \"{0}\" successfully updated", user.GetUsername(), params))

		case operations.RealmResource:
			var newRealm data.Realm
			if parseErr := json.Unmarshal(value, &newRealm); parseErr != nil {
				parseErrMsg := sf.Format("json.Unmarshal failed: {0}", parseErr.Error())
				//nolint:govet,staticcheck
				log.Println(parseErrMsg)
				os.Exit(-1)
			}
			if updateErr := manager.UpdateRealm(resourceId, newRealm); updateErr != nil {
				updateErrMsg := sf.Format("UpdateRealm failed: {0}", updateErr)
				//nolint:govet,staticcheck
				log.Println(updateErrMsg)
				os.Exit(-1)
			}
			fmt.Println(sf.Format("Realm: \"{0}\" successfully updated", newRealm.Name))
		case operations.UserFederationConfigResource:
			var userFederationServiceConfig data.UserFederationServiceConfig
			if err := json.Unmarshal(value, &userFederationServiceConfig); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			if err := manager.UpdateUserFederationConfig(params, resourceId, userFederationServiceConfig); err != nil {
				log.Fatalf("UpdateUserFederationConfig failed: %s", err)
			}
			fmt.Println(sf.Format("User federation service config: \"{0}\" successfully updated", userFederationServiceConfig.Name, params))
		case operations.ServerSettings:
			var security config.GlobalSecurityConfig
			if parseErr := json.Unmarshal(value, &security); parseErr != nil {
				msg := sf.Format("json.Unmarshal failed: {0}", parseErr)
				//nolint:govet,staticcheck
				log.Println(msg)
				os.Exit(-1)
			}
			serverSettings, readErr := manager.GetServerSettings()
			var encoder *encoding.PasswordJsonEncoder
			if readErr != nil && errors.As(readErr, &appErrs.EmptyNotFoundErr) {
				serverSettings = &data.ServerSettings{}
				salt := encoding.GenerateRandomSalt()
				encoder = encoding.NewPasswordJsonEncoder(salt)
				serverSettings.Admin.Id, _ = uuid.NewUUID()
				serverSettings.Admin.PasswordSalt = salt
			} else {
				encoder = encoding.NewPasswordJsonEncoder(serverSettings.Admin.PasswordSalt)
			}
			serverSettings.AllowedHosts = security.AllowedHosts
			serverSettings.AdminApiUrlPrefix = security.AdminApiUrlPrefix
			serverSettings.Admin.Username = security.Admin.Username
			serverSettings.Admin.PasswordHash = encoder.GetB64PasswordHash(serverSettings.Admin.PasswordSalt)

			if setErr := manager.SetServerSettings(serverSettings); setErr != nil {
				msg := sf.Format("UpdateUserFederationConfig failed: {0}", setErr)
				//nolint:govet,staticcheck
				log.Println(msg)
				os.Exit(-1)
			}
			fmt.Println("Server settings were successfully updated")
		}
		return
	case operations.ChangePassword:
		switch resource {
		case operations.AdminResource:
			fallthrough
		case operations.UserResource:
			fallthrough
		case "":
			password := string(value)
			if resource == operations.AdminResource {
				serverSettings, readErr := manager.GetServerSettings()
				if readErr != nil {
					msg := sf.Format("There is an error while getting ServerSettings: {0}", readErr.Error())
					//nolint:govet,staticcheck
					log.Println(msg)
					os.Exit(-1)
				}
				encoder := encoding.NewPasswordJsonEncoder(serverSettings.Admin.PasswordSalt)
				serverSettings.Admin.PasswordHash = encoder.GetB64PasswordHash(serverSettings.Admin.PasswordSalt)
				setError := manager.SetServerSettings(serverSettings)
				if setError != nil {
					msg := sf.Format("SetPassword for Admin failed: {0}", setError.Error())
					//nolint:govet,staticcheck
					log.Println(msg)
					os.Exit(-1)
				}
				fmt.Printf("Admin password successfully changed")
			} else {
				if params == "" {
					log.Fatalf("Not specified Params")
				}
				if resourceId == "" {
					log.Fatalf("Not specified Resource_id")
				}
				// TODO(SIA)  Moving password verification to another location
				if len(value) < 8 {
					log.Fatalf("Password length must be greater than 7")
				}
				passwordManager := manager.(PasswordManager)
				if setErr := passwordManager.SetPassword(params, resourceId, password); setErr != nil {
					log.Fatalf("SetPassword failed: %s", err)
				}
				fmt.Printf("Password successfully changed")
			}

		default:
			log.Fatalf("Bad Resource")
		}

		return
	case operations.ResetPassword:
		switch resource {
		case operations.UserResource:
			fallthrough
		case "":
			if params == "" {
				log.Fatalf("Not specified Params")
			}
			if resourceId == "" {
				log.Fatalf("Not specified ResourceId")
			}
			password := getRandPassword()
			passwordManager := manager.(PasswordManager)
			if err := passwordManager.SetPassword(params, resourceId, password); err != nil {
				log.Fatalf("SetPassword failed: %s", err)
			}
			fmt.Printf("New password: %s", password)

		default:
			log.Fatalf("Bad Resource")
		}

		return
	default:
		log.Fatalf("Bad Operation")
	}
}

func getRandPassword() string {
	// TODO(SIA) Move password generation to another location
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatalf("rand.Read failed: %s", err)
	}
	str := base32.StdEncoding.EncodeToString(randomBytes)
	const length = 8
	password := str[:length]
	return password
}
