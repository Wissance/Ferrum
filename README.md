# Ferrum

Ferrum is a **better** Authorization Server, this is a Community version.

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/wissance/Ferrum?style=plastic) 
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/wissance/Ferrum?style=plastic) 
![GitHub issues](https://img.shields.io/github/issues/wissance/Ferrum?style=plastic)
![GitHub Release Date](https://img.shields.io/github/release-date/wissance/Ferrum) 
![GitHub release (latest by date)](https://img.shields.io/github/downloads/wissance/Ferrum/v0.9.1/total?style=plastic)

![Ferrum: A better Auth Server](/img/ferrum_cover.png)

## 1. Communication

* Discord channel : https://discord.gg/9RYNYu2Mxq

## 2. General info

`Ferrum` is `OpenId-Connect` Authorization server written on GO. It has Data Contract similar to
`Keycloak` server (**minimal `Keycloak`** and we'll grow to full-fledged `KeyCloak` analog).

Today we are having **following features**:

1. Issue new tokens.
2. Refresh tokens.
2. Control user sessions (token expiration).
3. Get UserInfo.
4. Token Introspect.
4. Managed from external code (`Start` and `Stop`) making them an ***ideal candidate*** for using in ***integration
   tests*** for WEB API services that uses `Keycloak` as authorization server;
5. Ability to use different data storage:
   * `FILE` data storage for small Read only systems
   * `REDIS` data storage for systems with large number of users and small response time;
6. Ability to use any user data and attributes (any valid JSON but with some requirements), if you have to
   properly configure your users just add what user have to `data.json` or in memory
7. Ability to ***become high performance enterprise level Authorization server***.

it has `endpoints` SIMILAR to `Keycloak`, at present time we are having following:

1. Issue and Refresh tokens: `POST ~/auth/realms/{realm}/protocol/openid-connect/token`
2. Get UserInfo `GET  ~/auth/realms/{realm}/protocol/openid-connect/userinfo`
3. Introspect tokens `POST ~/auth/realms/{realm}/protocol/openid-connect/token/introspect`

## 3. How to use

### 3.1 Build

First of all build is simple run `go build` from application root directory. Additionally it is possible
to generate self signed certificates - run `go generate` from command line

If you don't specify the name of executable (by passing -o {execName} to go build) than name of executable = name of project

### 3.2 Run application as Standalone

Run is simple (`Ferrum` starts with default config - `config.json`):
```ps1
./Ferrum
```

To run `Ferrum` with selected config i.e. `config_w_redis.json` :

```ps1
./Ferrum --config ./config_w_redis.json
```

### 3.3 Run application in docker

It is possible to start app in docker with already installed `REDIS` and with initial data (see python
data insert script):

```ps1
    docker-compose up --build 
```

### 3.4 Run with direct configuration && data pass from code (embedding Authorization server in you applications)

There are 2 ways to use `Ferrum`:
1. Start with config file (described above)
2. Start with direct pass `config.AppConfig` and `data.ServerData` in application, i.e.
   ```go
    app := CreateAppWithData(appConfig, &testServerData, testKey)
	res, err := app.Init()
	assert.True(t, res)
	assert.Nil(t, err)

	res, err = app.Start()
	assert.True(t, res)
	assert.Nil(t, err)
	// do what you should ...
	app.Stop()
   ```

### Test
At present moment we have 2 fully integration tests, and number of them continues to grow. To run test execute from cmd:
```ps1
go test
```
For running Manager tests on `Redis` you must have redis on `127.0.0.1:6379` with `ferrum_db` / `FeRRuM000` `auth` `user+password`
pair, it is possible to start docker_compose and test on compose `ferrum_db` container 

## 4. Configure

### 4.1 Server configuration

Configuration splitted onto several sections:

        ```json
        "server": {
            "schema": "https",
            "address": "localhost",
            "port": 8182,
            "security": {
                "key_file": "./certs/server.key",
                "certificate_file": "./certs/server.crt"
            }
        }
        ```
      - data file: `realms`, `clients` and `users` application takes from this data file and stores in 
        app memory, data file name - `data.json`
      - key file that is using for `JWT` tokens generation (`access_token` && `refresh_token`), 
        name `keyfile` (without extensions).

### 4.2 Configure user data as you wish

Users does not have any specific structure, you could add whatever you want, but for compatibility
with keycloak and for ability to check password minimal user looks like:
```json
{
    "info": {
        "sub": "" // <-- THIS PROPERTY USED AS ID, PROBABLY WE SHOULD CHANGE THIS TO ID
        "preferred_username": "admin", // <-- THIS IS REQUIRED
        ...
    },
    "credentials": {
        "password": "1s2d3f4g90xs" // <-- TODAY WE STORE PASSWORDS AS OPENED
    }
}
```

in this minimal user example you could expand `info` structure as you want, `credentials` is a service structure,
there are NO SENSES in modifying it.

### 4.3 Server embedding into application (use from code)

Minimal full example of how to use coud be found in `application_test.go`, here is a minimal snippet:

```go
var testKey = []byte("qwerty1234567890")
var testServerData = data.ServerData{
	Realms: []data.Realm{
		{Name: "testrealm1", TokenExpiration: 10, RefreshTokenExpiration: 5,
			Clients: []data.Client{
				{Name: "testclient1", Type: data.Confidential, Auth: data.Authentication{Type: data.ClientIdAndSecrets,
					Value: "fb6Z4RsOadVycQoeQiN57xpu8w8wplYz"}},
			}, Users: []interface{}{
				map[string]interface{}{"info": map[string]interface{}{"sub": "667ff6a7-3f6b-449b-a217-6fc5d9ac0723",
					"name": "vano", "preferred_username": "vano",
					"given_name": "vano ivanov", "family_name": "ivanov", "email_verified": true},
					"credentials": map[string]interface{}{"password": "1234567890"}},
			}},
	},
}
var httpsAppConfig = config.AppConfig{ServerCfg: config.ServerConfig{Schema: config.HTTPS, Address: "127.0.0.1", Port: 8672,
	Security: config.SecurityConfig{KeyFile: "./certs/server.key", CertificateFile: "./certs/server.crt"}}}
	
app := CreateAppWithData(appConfig, &testServerData, testKey)
res, err := app.Init()
if err != nil {
	// handle ERROR
}

res, err = app.Start() 

if err != nil {
	// handle ERROR
}

// do whatever you want

app.Stop()
```

## 5. Server administer

Since version `0.9.1` it is possible to use `CLI Admin` [See](api/admin/cli/README.md)

### 5.1 Use CLI admin in a docker

1. Run docker compose - `docker compose up --build`
2. List running containers - `docker ps -a`
3. Attach to running container using listed hash `docker exec -it 060cfb8dd84c sh`
4. Run admin interface providing a valid config `ferrum-admin --config=config_docker_w_redis.json ...`, see picture

![Use CLI Admin from docker](/img/additional/cli_from_docker.png)

## 6. Contributors

<a href="https://github.com/Wissance/Ferrum/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Wissance/Ferrum" />
</a>
