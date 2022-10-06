package rest

import (
	"Ferrum/dto"
	"Ferrum/errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
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
		wCtx.Logger.Debug("New token issue: realm wasn't provided")
		result = dto.ErrorDetails{Msg: errors.RealmNotProviderMsg}
	} else {
		// todo: validate ...
		realmPtr := (*wCtx.DataProvider).GetRealm(realm)
		if realmPtr == nil {
			status = http.StatusNotFound
			wCtx.Logger.Debug("New token issue: realm doesn't exist")
			result = dto.ErrorDetails{Msg: stringFormatter.Format(errors.RealmDoesNotExistsTemplate, realm)}
		} else {
			// todo(UMV): think we don't have refresh strategy yet, add in v1.0 ...
			// New token issue strategy ...
			tokenGenerationData := dto.TokenGenerationData{}
			err := request.ParseForm()
			if err != nil {
				status = http.StatusBadRequest
				wCtx.Logger.Debug("New token issue: body is bad (unable to unmarshal to dto.TokenGenerationData)")
				result = dto.ErrorDetails{Msg: errors.BadBodyForTokenGenerationMsg}
			} else {
				var decoder = schema.NewDecoder()
				err = decoder.Decode(&tokenGenerationData, request.PostForm)
				if err != nil {
					// todo(UMV): log events
					status = http.StatusBadRequest
					wCtx.Logger.Debug("New token issue: body is bad (unable to unmarshal to dto.TokenGenerationData)")
					result = dto.ErrorDetails{Msg: errors.BadBodyForTokenGenerationMsg}
				} else {
					// 1. Validate client data: client_id, client_secret (if we have so), scope
					check := (*wCtx.Security).Validate(&tokenGenerationData, realmPtr)
					if check != nil {
						status = http.StatusBadRequest
						wCtx.Logger.Debug("New token issue: client data is invalid (client_id or client_secret)")
						result = dto.ErrorDetails{Msg: check.Msg, Description: check.Description}
					} else {
						// 2. Validate user credentials
						check = (*wCtx.Security).CheckCredentials(&tokenGenerationData, realmPtr)
						if check != nil {
							wCtx.Logger.Debug("New token issue: invalid user credentials (username or password)")
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
		wCtx.Logger.Debug("Get userinfo: realm wasn't provided")
		status = http.StatusBadRequest
		result = dto.ErrorDetails{Msg: errors.RealmNotProviderMsg}
	} else {
		realmPtr := (*wCtx.DataProvider).GetRealm(realm)
		if realmPtr == nil {
			wCtx.Logger.Debug("Get userinfo: realm doesn't exist")
			status = http.StatusNotFound
			result = dto.ErrorDetails{Msg: stringFormatter.Format(errors.RealmDoesNotExistsTemplate, realm)}
		} else {
			// Just get access token,  find user + session
			authorization := request.Header.Get("Authorization")
			parts := strings.Split(authorization, " ")
			if parts[0] != "Bearer" {
				wCtx.Logger.Debug("Get userinfo: expected only Bearer authorization yet")
				status = http.StatusBadRequest
				result = dto.ErrorDetails{Msg: errors.InvalidRequestMsg, Description: errors.InvalidRequestDesc}
			} else {
				session := (*wCtx.Security).GetSessionByAccessToken(realm, &parts[1])
				if session == nil {
					wCtx.Logger.Debug("Get userinfo: invalid token")
					status = http.StatusUnauthorized
					result = dto.ErrorDetails{Msg: errors.InvalidTokenMsg, Description: errors.InvalidTokenDesc}
				} else {
					if session.Expired.Before(time.Now()) {
						status = http.StatusUnauthorized
						wCtx.Logger.Debug("Get userinfo: token expired")
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

func (wCtx *WebApiContext) getRealmBaseUrl(realm string) string {
	return stringFormatter.Format("/{0}/{1}/auth/realms/{2}", "http", "localhost:8182", realm)
}
