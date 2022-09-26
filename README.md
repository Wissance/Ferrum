# Ferrum

Ferrum is a **better** Authorization Server.

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/wissance/Ferrum?style=plastic) 
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/wissance/Ferrum?style=plastic) 
![GitHub issues](https://img.shields.io/github/issues/wissance/Ferrum?style=plastic)
![GitHub Release Date](https://img.shields.io/github/release-date/wissance/Ferrum) 
![GitHub release (latest by date)](https://img.shields.io/github/downloads/wissance/Ferrum/v0.1.0/total?style=plastic)

![Ferrum: A better Auth Server](/img/ferrum_cover.png)

## Communication

* Discord channel : https://discord.gg/9RYNYu2Mxq

## General info

`Ferrum` is `OpenId-Connect` Authorization server written on GO. It behaves like `Keycloak` server (**minimal `Keycloak`**
 but we'll grow to full-fledged `KeyCloak`) and could:

1. issue new tokens:
2. control user sessions;
3. get useinfo;
4. managed from external code (`Start` and `Stop`) making them an ***ideal candidate*** for using in ***integration
   tests*** for WEB API services that uses `Keycloak` as authorization server;
5. ability to use any user data and attributes (any valid JSON but with some requirements), if you have to
   properly configure your users just add what user have to `data.json` or in memory
6. ability to ***become serious enterprise level Authorization server (GIVE US 100+ STARS on Github
   and we'll make them enterprise)***.

it has `endpoints` SIMILAR to `Keycloak`, at present time we are having following:
1. `POST "/auth/realms/{realm}/protocol/openid-connect/token/"``
2. `GET "/auth/realms/{realm}/protocol/openid-connect/userinfo/"`

## How to use

### Build

First of all build is simple run `go build` from application root directory. Additionally it is possible
to generate self signed certificates - run `go generate` from command line (if you are going to generate 
new certs **remove certs and key file from ./certs directory** prior to generate new)

If you don't specify the name of executable (by passing -o {execName} to go build) than name of executable = name of project

### Run application

run is simple:
```ps1
./Ferrum
```

There are 2 ways to use `Ferrum`:
1. It could be started as traditional application and depends from following files:
      - application config file - file with name `config.json` in repo we have config that starts
        application as `HTTP` if you need `HTTPS` change server section like this:
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
   
   Names are standard, in future we allow to pass own files from cmd
2. It could be started without reading files, all data could be passed as arguments (see 
   `application_test.go` for details):
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

## Configure user data as you wish

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

## Use from code

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
