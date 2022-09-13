package rest

import (
	"Ferrum/data"
	"Ferrum/dto"
	"Ferrum/errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/wissance/stringFormatter"
	"net/http"
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
							// 3. Create access token && refresh token
							// 4. Generate new token
							duration := 300
							refreshDuration := 120
							// 4. Save session
							session := (*wCtx.Security).StartOrUpdateSession(realm, (*currentUser).GetId(), duration)
							// 5. Generate new token
							result = dto.Token{AccessToken: "123445", Expires: duration, RefreshToken: "123",
								RefreshExpires: refreshDuration,
								TokenType:      "Bearer", NotBeforePolicy: 0, Session: session.String()}
							// todo(UMV): create JWT ...
							// token := jwt.NewWithClaims(jwt.SigningMethodHS256, currentUser)
						}
					}
				}
			}
		}
	}
	afterHandle(&respWriter, status, &result)
}

func (wCtx *WebApiContext) GetUserInfo(respWriter http.ResponseWriter, request *http.Request) {
	// Just get access token,  find user + session
}

func (wCtx *WebApiContext) generateAccessToken(realm string, tokenType string, scope string, sessionData *data.UserSession, userData *data.User) *data.AccessTokenData {
	// todo(UMV): store schema and pair address:port in the wCtx
	issuer := stringFormatter.Format("/{0}/{1}/auth/realms/{2}", "http", "localhost:8182", realm)
	jwtCommon := data.JwtTokenData{Issuer: issuer, Type: tokenType, Audience: "account", Scope: scope, JwtId: uuid.New(),
		IssuedAt: sessionData.Started, ExpiredAt: sessionData.Expired, Subject: sessionData.UserId,
		SessionId: sessionData.Id, SessionState: sessionData.Id}
	accessToken := data.CreateAccessToken(&jwtCommon, userData)
	return accessToken
}
