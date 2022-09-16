package main

import (
	"Ferrum/config"
	"Ferrum/data"
	"Ferrum/dto"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/wissance/stringFormatter"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
)

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

func TestApplicationOnHttp(t *testing.T) {
	appConfig := config.AppConfig{
		ServerCfg: config.ServerConfig{Schema: "http", Address: "127.0.0.1", Port: 8284},
	}
	app := CreateAppWithData(&appConfig, &testServerData, testKey)
	res, err := app.Init()
	assert.True(t, res)
	assert.Nil(t, err)

	res, err = app.Start()
	assert.True(t, res)
	assert.Nil(t, err)
	baseUrl := "http://127.0.0.1:8284"
	realm := "testrealm1"
	username := "vano"
	response := issueNewToken(t, baseUrl, realm, "testclient1", "fb6Z4RsOadVycQoeQiN57xpu8w8wplYz",
		username, "1234567890")
	token := getDataFromResponse[dto.Token](t, response)
	assert.True(t, len(token.AccessToken) > 0)
	assert.True(t, len(token.RefreshToken) > 0)
	// check token by query username
	userInfo := getUserInfo(t, baseUrl, realm, token.AccessToken, "200 OK")
	assert.True(t, len(userInfo) > 0)
	assert.Equal(t, username, userInfo["preferred_username"])
	time.Sleep(time.Second * time.Duration(10))
	userInfo = getUserInfo(t, baseUrl, realm, token.AccessToken, "401 Unauthorized")
	// wait token expiration and call one more, got 401

	// try with bad client credentials
	// try with bad user credentials

	res, err = app.Stop()
	assert.True(t, res)
	assert.Nil(t, err)
}

func issueNewToken(t *testing.T, baseUrl string, realm string, clientId string, clientSecret string,
	userName string, password string) *http.Response {
	tokenUrlTemplate := "{0}/auth/realms/{1}/protocol/openid-connect/token/"
	tokenUrl := stringFormatter.Format(tokenUrlTemplate, baseUrl, realm)
	getTokenData := url.Values{}
	getTokenData.Set("client_id", clientId)
	getTokenData.Set("client_secret", clientSecret)
	getTokenData.Set("scope", "profile")
	getTokenData.Set("grant_type", "password")
	getTokenData.Set("username", userName)
	getTokenData.Set("password", password)
	response, err := http.PostForm(tokenUrl, getTokenData)
	assert.Nil(t, err)
	return response
}

func getDataFromResponse[TR dto.Token | dto.ErrorDetails](t *testing.T, response *http.Response) TR {
	responseBody, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	var result TR
	err = json.Unmarshal(responseBody, &result)
	assert.Nil(t, err)
	return result
}

func getUserInfo(t *testing.T, baseUrl string, realm string, token string, expectedStatus string) map[string]interface{} {
	userInfoUrlTemplate := "{0}/auth/realms/{1}/protocol/openid-connect/userinfo/"
	userInfoUrl := stringFormatter.Format(userInfoUrlTemplate, baseUrl, realm)
	client := http.Client{}
	request, err := http.NewRequest("GET", userInfoUrl, nil)
	request.Header.Set("Authorization", "Bearer "+token)
	assert.Nil(t, err)
	response, err := client.Do(request)
	assert.Equal(t, expectedStatus, response.Status)
	assert.Nil(t, err)
	responseBody, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	assert.Nil(t, err)
	return result
}
