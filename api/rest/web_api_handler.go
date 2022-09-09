package rest

import (
	"Ferrum/dto"
	"Ferrum/errors"
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
							// 3. If all steps were passed return new dto.Token
							result = dto.Token{AccessToken: "123445"}
						}
					}
				}
			}
		}
	}
	afterHandle(&respWriter, status, &result)
}

func (wCtx *WebApiContext) GetUserInfo(respWriter http.ResponseWriter, request *http.Request) {

}
