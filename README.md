# Ferrum

![Ferrum: A better Auth Server](/img/ferrum_cover_sm.png)

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

First of all build is simple run `go build` from application root directory. Additionally it is possible
to generate self signed certificates - run `go generate` from command line (if you are going to generate 
new certs **remove certs and key file from ./certs directory** prior to generate new)

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
