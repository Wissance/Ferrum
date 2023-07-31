package rest

import (
	"encoding/json"
	"net/http"
)

const (
	authorizationHeader = "Authorization"
)

type tokenType string

const (
	BearerToken  tokenType = "Bearer"
	RefreshToken tokenType = "Refresh"
)

// beforeHandle
/* This function prepare response headers prior to response handle. It sets content-type and CORS headers.
 * Parameters:
 *     - respWriter - gorilla/mux response writer
 * Returns nothing
 */
func beforeHandle(respWriter *http.ResponseWriter) {
	(*respWriter).Header().Set("Content-Type", "application/json")
	(*respWriter).Header().Set("Accept", "application/json")
}

// afterHandle
/* This function finalize response handle: serialize (json) and write object and set status code. If error occur during object serialization status code sets to 500
 * Parameters:
 *     - respWriter - gorilla/mux response writer
 *     - statusCode - http response status
 *     - data - object (json) could be empty
 * Returns nothing
 */
func afterHandle(respWriter *http.ResponseWriter, statusCode int, data interface{}) {
	(*respWriter).WriteHeader(statusCode)
	if data != nil {
		err := json.NewEncoder(*respWriter).Encode(data)
		if err != nil {
			(*respWriter).WriteHeader(500)
		}
	}
}
