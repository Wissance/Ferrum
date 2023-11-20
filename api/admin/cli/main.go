package main

import (
	"fmt"
	"log"

	"github.com/wissance/Ferrum/api/admin/cli/config_cli"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/managers/redis_data_manager"

	"github.com/wissance/Ferrum/api/admin/cli/domain_cli"

	"github.com/wissance/Ferrum/logging"
	sf "github.com/wissance/stringFormatter"
)

type ManagerForCli interface {
	GetRealm(realmName string) (*data.Realm, error)
	GetClient(clientName string) (*data.Client, error)
	GetClientFromRealm(realmName string, clientName string) (*data.Client, error)
	GetUser(userName string) (data.User, error)
	GetUserFromRealm(realmName string, clientName string) (data.User, error)

	CreateRealm(realmValue []byte) (*data.Realm, error)
	CreateClient(clientValue []byte) (*data.Client, error)
	AddClientToRealm(realmName string, clientName string) error
	CreateUser(userValue []byte) (data.User, error)
	AddUserToRealm(realmName string, userName string) error

	DeleteRealm(realmName string) error
	DeleteClient(clientName string) error
	DeleteClientFromRealm(realmName string, clientName string) error
	DeleteUser(userName string) error
	DeleteUserFromRealm(realmName string, userName string) error

	UpdateClient(clientName string, clientValue []byte) (*data.Client, error)
	UpdateUser(userName string, userValue []byte) (data.User, error)
	UpdateRealm(realmName string, realmValue []byte) (*data.Realm, error)
}

func main() {
	cfg, err := config_cli.NewConfig()
	if err != nil {
		log.Fatalf("NewConfig failed: %s", err)
	}

	var manager ManagerForCli
	{
		logger := logging.CreateLogger(&cfg.LoggingConfig)
		redisManager, err := redis_data_manager.CreateRedisDataManager(&cfg.DataSourceConfig, logger)
		if err != nil {
			log.Fatalf("CreateRedisDataManager failed: %s", err)
		}
		manager = redisManager
	}

	switch cfg.Operation {
	case domain_cli.GetOperation:
		if cfg.Resource_id == "" {
			log.Fatalf("Not specified Resource_id")
		}
		switch cfg.Resource {
		case domain_cli.ClientResource:
			if cfg.Params == "" {
				client, err := manager.GetClient(cfg.Resource_id)
				if err != nil {
					log.Fatalf("GetClient failed: %s", err)
				}
				fmt.Println(*client)
			} else {
				clientIdAndName, err := manager.GetClientFromRealm(cfg.Params, cfg.Resource_id)
				if err != nil {
					log.Fatalf("GetClientFromRealm failed: %s", err)
				}
				fmt.Println(*clientIdAndName)
			}

		case domain_cli.UserResource:
			if cfg.Params == "" {
				user, err := manager.GetUser(cfg.Resource_id)
				if err != nil {
					log.Fatalf("GetUser failed: %s", err)
				}
				fmt.Println(user.GetUserInfo())
			} else {
				user, err := manager.GetUserFromRealm(cfg.Params, cfg.Resource_id)
				if err != nil {
					log.Fatalf("GetUserFromRealm failed: %s", err)
				}
				fmt.Println(user.GetUserInfo())
			}

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
		switch cfg.Resource {
		case domain_cli.ClientResource:
			if cfg.Params == "" {
				if len(cfg.Value) == 0 {
					log.Fatalf("Not specified Value")
				}
				client, err := manager.CreateClient(cfg.Value)
				if err != nil {
					log.Fatalf("CreateClient failed: %s", err)
				}
				fmt.Println(sf.Format("Client: \"{0}\" successfully created", client.Name))

			} else {
				if cfg.Resource_id == "" {
					log.Fatalf("Not specified Resource_id")
				}
				if err := manager.AddClientToRealm(cfg.Params, cfg.Resource_id); err != nil {
					log.Fatalf("AddClientToRealm failed: %s", err)
				}
				fmt.Println(sf.Format("Client: \"{0}\" successfully added to Realm: \"{1}\"", cfg.Resource_id, cfg.Params))
			}

		case domain_cli.UserResource:
			if cfg.Params == "" {
				if len(cfg.Value) == 0 {
					log.Fatalf("Not specified Value")
				}
				user, err := manager.CreateUser(cfg.Value)
				if err != nil {
					log.Fatalf("CreateUser failed: %s", err)
				}
				fmt.Println(sf.Format("User: \"{0}\" successfully created", user.GetUsername()))

			} else {
				if cfg.Resource_id == "" {
					log.Fatalf("Not specified Resource_id")
				}
				if err := manager.AddUserToRealm(cfg.Params, cfg.Resource_id); err != nil {
					log.Fatalf("AddUserToRealm failed: %s", err)
				}
				fmt.Println(sf.Format("User: \"{0}\" successfully added to Realm: \"{1}\"", cfg.Resource_id, cfg.Params))
			}

		case domain_cli.RealmResource:
			if len(cfg.Value) == 0 {
				log.Fatalf("Not specified Value")
			}
			// создает клиентов и пользователей, создает новые realmClients и realmUsers, создает realm
			realm, err := manager.CreateRealm(cfg.Value)
			if err != nil {
				log.Fatalf("CreateRealm failed: %s", err)
			}
			fmt.Println(sf.Format("Realm: \"{0}\" successfully created", realm.Name))
			return
		}

		return
	case domain_cli.DeleteOperation:
		if cfg.Resource_id == "" {
			log.Fatalf("Not specified Resource_id")
		}
		switch cfg.Resource {
		case domain_cli.ClientResource:
			if cfg.Params == "" {
				if err := manager.DeleteClient(cfg.Resource_id); err != nil {
					log.Fatalf("DeleteClient failed: %s", err)
				}
				fmt.Println(sf.Format("Client: \"{0}\" successfully deleted", cfg.Resource_id))
			} else {
				// Удаляет клиента из realmClients. Удаление самого клиента не происходит
				if err := manager.DeleteClientFromRealm(cfg.Params, cfg.Resource_id); err != nil {
					log.Fatalf("DeleteClientFromRealm failed: %s", err)
				}
				fmt.Println(sf.Format("Client: \"{0}\" successfully deleted in Realm: \"{1}\"", cfg.Resource_id, cfg.Params))
			}

		case domain_cli.UserResource:
			if cfg.Params == "" {
				if err := manager.DeleteUser(cfg.Resource_id); err != nil {
					log.Fatalf("DeleteUser failed: %s", err)
				}
				fmt.Println(sf.Format("User: \"{0}\" successfully deleted", cfg.Resource_id))
			} else {
				// Удаляет user из realmUsers. Удаление самого клиента не происходит
				if err := manager.DeleteUserFromRealm(cfg.Params, cfg.Resource_id); err != nil {
					log.Fatalf("DeleteUserFromRealm failed: %s", err)
				}
				fmt.Println(sf.Format("User: \"{0}\" successfully deleted in Realm: \"{1}\"", cfg.Resource_id, cfg.Params))
			}

		case domain_cli.RealmResource:
			// Удаляет realmClients и realmUsers и realm. Удаление самих client и user не происходит.
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
			client, err := manager.UpdateClient(cfg.Resource_id, cfg.Value)
			if err != nil {
				log.Fatalf("UpdateClient failed: %s", err)
			}
			fmt.Println(sf.Format("Client: \"{0}\" successfully updated", client.Name))

		case domain_cli.UserResource:
			user, err := manager.UpdateUser(cfg.Resource_id, cfg.Value)
			if err != nil {
				log.Fatalf("UpdateUser failed: %s", err)
			}
			fmt.Println(sf.Format("User: \"{0}\" successfully updated", user.GetUsername(), cfg.Params))

		case domain_cli.RealmResource:
			realm, err := manager.UpdateRealm(cfg.Resource_id, cfg.Value)
			if err != nil {
				log.Fatalf("UpdateRealm failed: %s", err)
			}
			fmt.Println(sf.Format("Realm: \"{0}\" successfully updated", realm.Name))
		}

		return
	case "change_password":
		fmt.Println("change_password")

	case "reset_password ":
		fmt.Println("reset_password")

	default:
		log.Fatalf("Bad Operation") // TODO
	}
}
