package rest

import (
	"encoding/base64"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/wissance/Ferrum/dto"
	"github.com/wissance/Ferrum/errors"
	"github.com/wissance/stringFormatter"
	"net/http"
	"strings"
	"time"
)

func (wCtx *WebApiContext) IssueNewToken(respWriter http.ResponseWriter, request *http.Request) {
	beforeHandle(&respWriter)
	vars := mux.Vars(request)
	realm := vars["realm"]
	var result interface{}
	status := http.StatusOK
	if len(realm) == 0 {
		// 400
		status = http.StatusBadRequest
		result = dto.ErrorDetails{Msg: errors.RealmNotProviderMsg}
	} else {
		// todo: validate ...
		realmPtr := (*wCtx.DataProvider).GetRealm(realm)
		if realmPtr == nil {
			status = http.StatusNotFound
			result = dto.ErrorDetails{Msg: stringFormatter.Format(errors.RealmDoesNotExistsTemplate, realm)}
		} else {
			// todo(UMV): think we don't have refresh strategy yet, add in v1.0 ...
			// New token issue strategy ...
			tokenGenerationData := dto.TokenGenerationData{}
			err := request.ParseForm()
			if err != nil {
				status = http.StatusBadRequest
				result = dto.ErrorDetails{Msg: errors.BadBodyForTokenGenerationMsg}
			} else {
				var decoder = schema.NewDecoder()
				err = decoder.Decode(&tokenGenerationData, request.PostForm)
				if err != nil {
					// todo(UMV): log events
					status = http.StatusBadRequest
					result = dto.ErrorDetails{Msg: errors.BadBodyForTokenGenerationMsg}
				} else {
					// 1. Validate client data: client_id, client_secret (if we have so), scope
					check := (*wCtx.Security).Validate(&tokenGenerationData, realmPtr)
					if check != nil {
						status = http.StatusBadRequest
						result = dto.ErrorDetails{Msg: check.Msg, Description: check.Description}
					} else {
						// 2. Validate user credentials
						check = (*wCtx.Security).CheckCredentials(&tokenGenerationData, realmPtr)
						if check != nil {
							status = http.StatusUnauthorized
							result = dto.ErrorDetails{Msg: check.Msg, Description: check.Description}
						} else {
							currentUser := (*wCtx.Security).GetCurrentUser(realmPtr, tokenGenerationData.Username)
							userId := (*currentUser).GetId()
							// 3. Create access token && refresh token
							// 4. Generate new token
							duration := realmPtr.TokenExpiration
							refreshDuration := realmPtr.RefreshTokenExpiration
							// 4. Save session
							sessionId := (*wCtx.Security).StartOrUpdateSession(realm, userId, duration)
							session := (*wCtx.Security).GetSession(realm, userId)
							// 5. Generate new tokens
							accessToken := wCtx.TokenGenerator.GenerateJwtAccessToken(wCtx.getRealmBaseUrl(realm), "Bearer", "profile email", session, currentUser)
							refreshToken := wCtx.TokenGenerator.GenerateJwtRefreshToken(wCtx.getRealmBaseUrl(realm), "Refresh", "profile email", session)
							(*wCtx.Security).AssignTokens(realm, userId, &accessToken, &refreshToken)
							// 6. Assign token to result
							result = dto.Token{AccessToken: accessToken, Expires: duration, RefreshToken: refreshToken,
								RefreshExpires: refreshDuration, TokenType: "Bearer", NotBeforePolicy: 0,
								Session: sessionId.String()}

						}
					}
				}
			}
		}
	}
	afterHandle(&respWriter, status, &result)
}

func (wCtx *WebApiContext) GetUserInfo(respWriter http.ResponseWriter, request *http.Request) {
	beforeHandle(&respWriter)
	vars := mux.Vars(request)
	realm := vars["realm"]
	var result interface{}
	status := http.StatusOK
	if len(realm) == 0 {
		// 400
		status = http.StatusBadRequest
		result = dto.ErrorDetails{Msg: errors.RealmNotProviderMsg}
	} else {
		realmPtr := (*wCtx.DataProvider).GetRealm(realm)
		if realmPtr == nil {
			status = http.StatusNotFound
			result = dto.ErrorDetails{Msg: stringFormatter.Format(errors.RealmDoesNotExistsTemplate, realm)}
		} else {
			// Just get access token,  find user + session
			authorization := request.Header.Get("Authorization")
			parts := strings.Split(authorization, " ")
			if parts[0] != "Bearer" {
				status = http.StatusBadRequest
				result = dto.ErrorDetails{Msg: errors.InvalidRequestMsg, Description: errors.InvalidRequestDesc}
			} else {
				session := (*wCtx.Security).GetSessionByAccessToken(realm, &parts[1])
				if session == nil {
					status = http.StatusUnauthorized
					result = dto.ErrorDetails{Msg: errors.InvalidTokenMsg, Description: errors.InvalidTokenDesc}
				} else {
					if session.Expired.Before(time.Now()) {
						status = http.StatusUnauthorized
						result = dto.ErrorDetails{Msg: errors.InvalidTokenMsg, Description: errors.InvalidTokenDesc}
					} else {
						user := (*wCtx.DataProvider).GetUserById(realmPtr, session.UserId)
						status = http.StatusOK
						if user != nil {
							result = (*user).GetUserInfo()
						}
					}
				}
			}
		}
	}
	afterHandle(&respWriter, status, &result)
}
func (wCtx *WebApiContext) Introspect(respWriter http.ResponseWriter, request *http.Request) {
	beforeHandle(&respWriter)
	vars := mux.Vars(request)
	realm := vars["realm"]
	if len(realm) == 0 {
		// 400
		status := http.StatusBadRequest
		result := dto.ErrorDetails{Msg: errors.RealmNotProviderMsg}
		afterHandle(&respWriter, status, &result)
		return
	}
	realmPtr := (*wCtx.DataProvider).GetRealm(realm)
	if realmPtr == nil {
		status := http.StatusNotFound
		result := dto.ErrorDetails{Msg: stringFormatter.Format(errors.RealmDoesNotExistsTemplate, realm)}
		afterHandle(&respWriter, status, &result)
		return
	}
	authorization := request.Header.Get("Authorization")
	parts := strings.Split(authorization, " ")
	if parts[0] != "Basic" {
		status := http.StatusBadRequest
		result := dto.ErrorDetails{Msg: errors.InvalidRequestMsg, Description: errors.InvalidRequestDesc}
		afterHandle(&respWriter, status, &result)
		return
	}
	basicString, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		status := http.StatusBadRequest
		result := dto.ErrorDetails{Msg: errors.InvalidClientMsg, Description: errors.InvalidClientCredentialDesc}
		afterHandle(&respWriter, status, &result)
		return
	}
	secretPair := strings.Split(string(basicString), ":")
	checkResult := (*wCtx.Security).Validate(&dto.TokenGenerationData{
		ClientSecret: secretPair[1],
		ClientId:     secretPair[0],
	}, realmPtr)
	if checkResult != nil {
		status := http.StatusUnauthorized
		result := dto.ErrorDetails{Msg: errors.InvalidClientMsg, Description: errors.InvalidClientCredentialDesc}
		afterHandle(&respWriter, status, &result)
		return
	}
	token := request.FormValue("token")
	session := (*wCtx.Security).GetSessionByAccessToken(realm, &token)
	if session == nil {
		status := http.StatusUnauthorized
		result := dto.ErrorDetails{Msg: errors.InvalidTokenMsg, Description: errors.InvalidTokenDesc}
		afterHandle(&respWriter, status, &result)
		return
	}
	active := !session.Expired.Before(time.Now())
	status := http.StatusOK
	tokenType := "Bearer"
	result := dto.IntrospectTokenResult{
		Active: active,
		Type:   tokenType,
		Exp:    realmPtr.TokenExpiration,
	}
	afterHandle(&respWriter, status, &result)
}

//func checkSecret(realmPtr *data.Realm, secretPair []string) bool {
//	clientId := secretPair[0]
//	clientSecret := secretPair[1]
//	for _, c := range realmPtr.Clients {
//		if c.Name == clientId {
//			if c.Type == data.Public {
//				return true
//			}
//
//			// here we make deal with confidential client
//			if c.Auth.Type == data.ClientIdAndSecrets && c.Auth.Value == clientSecret {
//				return true
//			}
//
//		}
//	}
//	return false
//}
func (wCtx *WebApiContext) getRealmBaseUrl(realm string) string {
	return stringFormatter.Format("/{0}/{1}/auth/realms/{2}", "http", "localhost:8182", realm)
}
