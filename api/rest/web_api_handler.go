package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (wCtx *WebApiContext) IssueNewToken(respWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	realm := vars["realm"]
	if len(realm) == 0 {
		// 400
	}
}

func (wCtx *WebApiContext) GetUserInfo(respWriter http.ResponseWriter, request *http.Request) {

}
