package main

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"log"

	"github.com/wissance/Ferrum/api/admin/cli/config_cli"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/managers/redis_data_manager"

	"github.com/wissance/Ferrum/api/admin/cli/domain_cli"

	"github.com/wissance/Ferrum/logging"
	sf "github.com/wissance/stringFormatter"
)

// type ManagerForCli interface {
// 	GetRealm(realmName string) (*data.Realm, error)
// 	GetClient(realmName string, clientName string) (*data.Client, error)
// 	GetUser(realmName string, userName string) (data.User, error)

// 	CreateRealm(newRealm data.Realm) error
// 	CreateClient(realmName string, clientNew data.Client) error
// 	CreateUser(realmName string, userNew data.User) error

// 	DeleteRealm(realmName string) error
// 	DeleteClient(realmName string, clientName string) error
// 	DeleteUser(realmName string, userName string) error

// 	UpdateRealm(realmName string, realmNew data.Realm) error
// 	UpdateClient(realmName string, clientName string, clientNew data.Client) error
// 	UpdateUser(realmName string, userName string, userNew data.User) error
// }

func main() {
	cfg, err := config_cli.NewConfig()
	if err != nil {
		log.Fatalf("NewConfig failed: %s", err)
	}

	//var manager ManagerForCli
	//{
	logger := logging.CreateLogger(&cfg.LoggingConfig)
	redisManager, err := redis_data_manager.CreateRedisDataManager(&cfg.DataSourceConfig, logger)
	if err != nil {
		log.Fatalf("CreateRedisDataManager failed: %s", err)
	}
	manager := redisManager
	//}

	if cfg.Operation != domain_cli.GetOperation && cfg.Operation != domain_cli.CreateOperation &&
		cfg.Operation != domain_cli.DeleteOperation && cfg.Operation != domain_cli.UpdateOperation &&
		cfg.Operation != domain_cli.ChangePassword && cfg.Operation != domain_cli.ResetPassword {
		log.Fatalf("bad Operation \"%s\"", cfg.Operation)
	}
	if !(cfg.Operation == domain_cli.ChangePassword || cfg.Operation == domain_cli.ResetPassword) {
		if cfg.Resource != domain_cli.RealmResource && cfg.Resource != domain_cli.ClientResource && cfg.Resource != domain_cli.UserResource {
			log.Fatalf("bad Resource \"%s\"", cfg.Resource)
		}
	}
	if (cfg.Resource_id == domain_cli.ClientResource) || (cfg.Resource_id == domain_cli.UserResource) {
		if cfg.Params == "" {
			log.Fatalf("Not specified Params")
		}
	}

	switch cfg.Operation {
	case domain_cli.GetOperation:
		if cfg.Resource_id == "" {
			log.Fatalf("Not specified Resource_id")
		}
		switch cfg.Resource {
		case domain_cli.ClientResource:
			client, err := manager.GetClient(cfg.Params, cfg.Resource_id)
			if err != nil {
				log.Fatalf("GetClient failed: %s", err)
			}
			fmt.Println(*client)

		case domain_cli.UserResource:
			user, err := manager.GetUser(cfg.Params, cfg.Resource_id)
			if err != nil {
				log.Fatalf("GetUser failed: %s", err)
			}
			fmt.Println(user.GetUserInfo())

		case domain_cli.RealmResource:
			realm, err := manager.GetRealm(cfg.Resource_id)
			if err != nil {
				log.Fatalf("GetRealm failed: %s", err)
				// log.Fatal(sf.Format("Realm: \"{0}\" doesn't exist", cfg.Resource_id))
			}
			fmt.Println(*realm)
		}

		return
	case domain_cli.CreateOperation:
		if len(cfg.Value) == 0 {
			log.Fatalf("Not specified Value")
		}
		switch cfg.Resource {
		case domain_cli.ClientResource:
			var clientNew data.Client
			if err := json.Unmarshal(cfg.Value, &clientNew); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			if err := manager.CreateClient(cfg.Params, clientNew); err != nil {
				log.Fatalf("CreateClient failed: %s", err)
			}
			fmt.Println(sf.Format("Client: \"{0}\" successfully created", clientNew.Name))

		case domain_cli.UserResource:
			var userNew any
			if err := json.Unmarshal(cfg.Value, &userNew); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			user := data.CreateUser(userNew)
			if err := manager.CreateUser(cfg.Params, user); err != nil {
				log.Fatalf("CreateUser failed: %s", err)
			}
			fmt.Println(sf.Format("User: \"{0}\" successfully created", user.GetUsername()))

		case domain_cli.RealmResource:
			var newRealm data.Realm
			if err := json.Unmarshal(cfg.Value, &newRealm); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			if err := manager.CreateRealm(newRealm); err != nil {
				log.Fatalf("CreateRealm failed: %s", err)
			}
			fmt.Println(sf.Format("Realm: \"{0}\" successfully created", newRealm.Name))
		}

		return
	case domain_cli.DeleteOperation:
		if cfg.Resource_id == "" {
			log.Fatalf("Not specified Resource_id")
		}
		switch cfg.Resource {
		case domain_cli.ClientResource:
			if err := manager.DeleteClient(cfg.Params, cfg.Resource_id); err != nil {
				log.Fatalf("DeleteClient failed: %s", err)
			}
			fmt.Println(sf.Format("Client: \"{0}\" successfully deleted", cfg.Resource_id))

		case domain_cli.UserResource:
			if err := manager.DeleteUser(cfg.Params, cfg.Resource_id); err != nil {
				log.Fatalf("DeleteUser failed: %s", err)
			}
			fmt.Println(sf.Format("User: \"{0}\" successfully deleted", cfg.Resource_id))

		case domain_cli.RealmResource:
			if err := manager.DeleteRealm(cfg.Resource_id); err != nil {
				log.Fatalf("DeleteRealm failed: %s", err)
			}
			fmt.Println(sf.Format("Realm: \"{0}\" successfully deleted", cfg.Resource_id))
		}

		return
	case domain_cli.UpdateOperation:
		if cfg.Resource_id == "" {
			log.Fatalf("Not specified Resource_id")
		}
		if len(cfg.Value) == 0 {
			log.Fatalf("Not specified Value")
		}
		switch cfg.Resource {
		case domain_cli.ClientResource:
			var newClient data.Client
			if err := json.Unmarshal(cfg.Value, &newClient); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			if err := manager.UpdateClient(cfg.Params, cfg.Resource_id, newClient); err != nil {
				log.Fatalf("UpdateClient failed: %s", err)
			}
			fmt.Println(sf.Format("Client: \"{0}\" successfully updated", newClient.Name))

		case domain_cli.UserResource:
			var newUser any
			if err := json.Unmarshal(cfg.Value, &newUser); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			user := data.CreateUser(newUser)
			if err := manager.UpdateUser(cfg.Params, cfg.Resource_id, user); err != nil {
				log.Fatalf("UpdateUser failed: %s", err)
			}
			fmt.Println(sf.Format("User: \"{0}\" successfully updated", user.GetUsername(), cfg.Params))

		case domain_cli.RealmResource:
			var newRealm data.Realm
			if err := json.Unmarshal(cfg.Value, &newRealm); err != nil {
				log.Fatalf("json.Unmarshal failed: %s", err)
			}
			if err := manager.UpdateRealm(cfg.Resource_id, newRealm); err != nil {
				log.Fatalf("UpdateRealm failed: %s", err)
			}
			fmt.Println(sf.Format("Realm: \"{0}\" successfully updated", newRealm.Name))
		}

		return
	case domain_cli.ChangePassword:
		switch cfg.Resource {
		case domain_cli.UserResource:
			fallthrough
		case "":
			if cfg.Params == "" {
				log.Fatalf("Not specified Params")
			}
			if cfg.Resource_id == "" {
				log.Fatalf("Not specified Resource_id")
			}
			// TODO(SIA)  Вынести валидацию пароля в другое место
			if len(cfg.Value) < 8 {
				log.Fatalf("Password length must be greater than 8")
			}
			password := string(cfg.Value)
			if err := manager.SetPassword(cfg.Params, cfg.Resource_id, password); err != nil {
				log.Fatalf("SetPassword failed: %s", err)
			}
			fmt.Printf("Password successfully changed")

		default:
			log.Fatalf("Bad Resource")
		}

		return
	case domain_cli.ResetPassword:
		switch cfg.Resource {
		case domain_cli.UserResource:
			fallthrough
		case "":
			if cfg.Params == "" {
				log.Fatalf("Not specified Params")
			}
			if cfg.Resource_id == "" {
				log.Fatalf("Not specified Resource_id")
			}
			password := getRandPassword()
			if err := manager.SetPassword(cfg.Params, cfg.Resource_id, password); err != nil {
				log.Fatalf("SetPassword failed: %s", err)
			}
			fmt.Printf("New password: %s", password)

		default:
			log.Fatalf("Bad Resource")
		}

		return
	default:
		log.Fatalf("Bad Operation") // TODO
	}
}

func getRandPassword() string {
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
