package rest

import (
	"Ferrum/dto"
	"github.com/gorilla/mux"
	"github.com/wissance/stringFormatter"
	"net/http"
)

const (
	realmDoesNotExistsTemplate = "Realm \"{0}\" does not exists"
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
	} else {
		// todo: validate ...
		realmPtr := (*wCtx.DataProvider).GetRealm(realm)
		if realmPtr == nil {
			status = http.StatusNotFound
			result = dto.ErrorDetails{Msg: stringFormatter.Format(realmDoesNotExistsTemplate, realm)}
		} else {
			// todo(UMV): think we don't have refresh strategy yet, add in v1.0 ...
			// 1. Validate client data: client_id, client_secret (if we have so), scope
			// 2. Validate user credentials
			// 3. If all steps were passed return new dto.Token
			result = dto.Token{AccessToken: "123445"}
		}
	}
	afterHandle(&respWriter, status, &result)
}

func (wCtx *WebApiContext) GetUserInfo(respWriter http.ResponseWriter, request *http.Request) {

}
