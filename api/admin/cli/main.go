package main

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/wissance/Ferrum/api/admin/cli/operations"
	"github.com/wissance/Ferrum/config"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/logging"
	sf "github.com/wissance/stringFormatter"
)

const defaultConfig = "./config_w_redis.json"

var (
	argConfigFile  = flag.String("config", defaultConfig, "")
	argOperation   = flag.String("operation", "", "")
	argResource    = flag.String("resource", "", "")
	argResource_id = flag.String("resource_id", "", "")
	argParams      = flag.String("params", "", "This is the name of the realm for operations on client or user resources")
	argValue       = flag.String("value", "", "Json object")
)

func main() {
	flag.Parse()
	cfg, err := config.ReadAppConfig(*argConfigFile)
	if err != nil {
		log.Fatalf("readAppConfig failed: %s", err)
	}
	logger := logging.CreateLogger(&cfg.Logging)
	manager, err := prepareContext(&cfg.DataSource, logger)
	if err != nil {
		log.Fatalf("prepareContext failed: %s", err)
	}

	operation := operations.OperationType(*argOperation)
	resource := operations.ResourceType(*argResource)
	resource_id := *argResource_id
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
		isInvalidResource := resource != operations.RealmResource && resource != operations.ClientResource && resource != operations.UserResource
		if isInvalidResource {
			log.Fatalf("bad Resource \"%s\"", resource)
		}
	}
	if (resource == operations.ClientResource) || (resource == operations.UserResource) {
		if params == "" {
			log.Fatalf("Not specified Params")
		}
	}

	switch operation {
	case operations.GetOperation:
		if resource_id == "" {
			log.Fatalf("Not specified Resource_id")
		}
		switch resource {
		case operations.ClientResource:
			client, err := manager.GetClient(params, resource_id)
			if err != nil {
				log.Fatalf("GetClient failed: %s", err)
			}
			fmt.Println(*client)

		case operations.UserResource:
			user, err := manager.GetUser(params, resource_id)
			if err != nil {
				log.Fatalf("GetUser failed: %s", err)
			}
			fmt.Println(user.GetUserInfo())

		case operations.RealmResource:
			realm, err := manager.GetRealm(resource_id)
			if err != nil {
				log.Fatalf("GetRealm failed: %s", err)
			}
			fmt.Println(*realm)
		}

		return
	case operations.CreateOperation:
		if len(value) == 0 {
			log.Fatalf("Not specified Value")
		}
		switch resource {
		case operations.ClientResource:
			var clientNew data.Client
			if err := json.Unmarshal(value, &clientNew); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			if err := manager.CreateClient(params, clientNew); err != nil {
				log.Fatalf("CreateClient failed: %s", err)
			}
			fmt.Println(sf.Format("Client: \"{0}\" successfully created", clientNew.Name))

		case operations.UserResource:
			var userNew any
			if err := json.Unmarshal(value, &userNew); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			user := data.CreateUser(userNew)
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
		}

		return
	case operations.DeleteOperation:
		if resource_id == "" {
			log.Fatalf("Not specified Resource_id")
		}
		switch resource {
		case operations.ClientResource:
			if err := manager.DeleteClient(params, resource_id); err != nil {
				log.Fatalf("DeleteClient failed: %s", err)
			}
			fmt.Println(sf.Format("Client: \"{0}\" successfully deleted", resource_id))

		case operations.UserResource:
			if err := manager.DeleteUser(params, resource_id); err != nil {
				log.Fatalf("DeleteUser failed: %s", err)
			}
			fmt.Println(sf.Format("User: \"{0}\" successfully deleted", resource_id))

		case operations.RealmResource:
			if err := manager.DeleteRealm(resource_id); err != nil {
				log.Fatalf("DeleteRealm failed: %s", err)
			}
			fmt.Println(sf.Format("Realm: \"{0}\" successfully deleted", resource_id))
		}

		return
	case operations.UpdateOperation:
		if resource_id == "" {
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
			if err := manager.UpdateClient(params, resource_id, newClient); err != nil {
				log.Fatalf("UpdateClient failed: %s", err)
			}
			fmt.Println(sf.Format("Client: \"{0}\" successfully updated", newClient.Name))

		case operations.UserResource:
			var newUser any
			if err := json.Unmarshal(value, &newUser); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			user := data.CreateUser(newUser)
			if err := manager.UpdateUser(params, resource_id, user); err != nil {
				log.Fatalf("UpdateUser failed: %s", err)
			}
			fmt.Println(sf.Format("User: \"{0}\" successfully updated", user.GetUsername(), params))

		case operations.RealmResource:
			var newRealm data.Realm
			if err := json.Unmarshal(value, &newRealm); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			if err := manager.UpdateRealm(resource_id, newRealm); err != nil {
				log.Fatalf("UpdateRealm failed: %s", err)
			}
			fmt.Println(sf.Format("Realm: \"{0}\" successfully updated", newRealm.Name))
		}

		return
	case operations.ChangePassword:
		switch resource {
		case operations.UserResource:
			fallthrough
		case "":
			if params == "" {
				log.Fatalf("Not specified Params")
			}
			if resource_id == "" {
				log.Fatalf("Not specified Resource_id")
			}
			// TODO(SIA)  Moving password verification to another location
			if len(value) < 8 {
				log.Fatalf("Password length must be greater than 8")
			}
			password := string(value)
			if err := manager.SetPassword(params, resource_id, password); err != nil {
				log.Fatalf("SetPassword failed: %s", err)
			}
			fmt.Printf("Password successfully changed")

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
			if resource_id == "" {
				log.Fatalf("Not specified Resource_id")
			}
			password := getRandPassword()
			if err := manager.SetPassword(params, resource_id, password); err != nil {
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
