package rest

import (
	"Ferrum/data"
	"Ferrum/dto"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/wissance/stringFormatter"
	"net/http"
)

const (
	realmNotProviderMsg          = "You does not provided any realm"
	realmDoesNotExistsTemplate   = "Realm \"{0}\" does not exists"
	badBodyForTokenGenerationMsg = "Bad body for token generation, see documentations"
	invalidClientMsg             = "Invalid client"
	invalidClientCredentialDesc  = "Invalid client credentials"
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
		result = dto.ErrorDetails{Msg: realmNotProviderMsg}
	} else {
		// todo: validate ...
		realmPtr := (*wCtx.DataProvider).GetRealm(realm)
		if realmPtr == nil {
			status = http.StatusNotFound
			result = dto.ErrorDetails{Msg: stringFormatter.Format(realmDoesNotExistsTemplate, realm)}
		} else {
			// todo(UMV): think we don't have refresh strategy yet, add in v1.0 ...
			// New token issue strategy ...
			tokenGenerationData := dto.TokenGenerationData{}
			err := request.ParseForm()
			if err != nil {
				status = http.StatusBadRequest
				result = dto.ErrorDetails{Msg: badBodyForTokenGenerationMsg}
			} else {
				var decoder = schema.NewDecoder()
				err = decoder.Decode(&tokenGenerationData, request.PostForm)
				if err != nil {
					// todo(UMV): log events
					status = http.StatusBadRequest
					result = dto.ErrorDetails{Msg: badBodyForTokenGenerationMsg}
				} else {
					// 1. Validate client data: client_id, client_secret (if we have so), scope
					check := wCtx.validate(&tokenGenerationData, realmPtr)
					if check != nil {
						status = http.StatusBadRequest
						result = *check
					} else {
						// 2. Validate user credentials
						// 3. If all steps were passed return new dto.Token
						result = dto.Token{AccessToken: "123445"}
					}
				}
			}
		}
	}
	afterHandle(&respWriter, status, &result)
}

func (wCtx *WebApiContext) GetUserInfo(respWriter http.ResponseWriter, request *http.Request) {

}

// todo(UMV): this is temporary, MUST move to some service, not related to web
func (wCtx *WebApiContext) validate(tokenGenData *dto.TokenGenerationData, realm *data.Realm) *dto.ErrorDetails {
	for _, c := range realm.Clients {
		if c.Name == tokenGenData.ClientId {
			if c.Type == data.Public {
				return nil
			}

			// here we make deal with confidential client
			if c.Auth.Type == data.ClientIdAndSecrets && c.Auth.Value == tokenGenData.ClientSecret {
				return nil
			}

		}
	}
	return &dto.ErrorDetails{Msg: invalidClientMsg, Description: invalidClientCredentialDesc}
}
