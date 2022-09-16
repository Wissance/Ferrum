package main

import (
	"Ferrum/config"
	"Ferrum/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testKey = []byte("qwerty1234567890")
var testServerData = data.ServerData{
	Realms: []data.Realm{
		{Name: "testrealm1", TokenExpiration: 30, RefreshTokenExpiration: 20,
			Clients: []data.Client{
				{Name: "testclient1", Type: data.Confidential, Auth: data.Authentication{Value: "atatatata"}},
			}},
	},
}

func TestApplicationOnHttp(t *testing.T) {
	appConfig := config.AppConfig{
		ServerCfg: config.ServerConfig{Schema: "http", Address: "localhost", Port: 8284},
	}
	app := CreateAppWithData(&appConfig, &testServerData, testKey)
	res, err := app.Init()
	assert.True(t, res)
	assert.Nil(t, err)

	res, err = app.Start()
	assert.True(t, res)
	assert.Nil(t, err)

	//todo(UMV): make here sets of API calls

	res, err = app.Stop()
	assert.True(t, res)
	assert.Nil(t, err)
}
