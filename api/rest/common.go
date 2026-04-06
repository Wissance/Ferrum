package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
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
 *     - respWriter - gin response writer
 * Returns nothing
 */
func beforeHandle(respWriter *gin.ResponseWriter) {
	(*respWriter).Header().Set("Content-Type", "application/json")
	(*respWriter).Header().Set("Accept", "application/json")
}

// afterHandle
/* This function finalize response handle: serialize (json) and write object and set status code. If error occur during object serialization status code sets to 500
 * Parameters:
 *     - respWriter - gin response writer
 *     - statusCode - http response status
 *     - data - object (json) could be empty
 * Returns nothing
 */
func afterHandle(respWriter *gin.ResponseWriter, statusCode int, data interface{}) {
	(*respWriter).WriteHeader(statusCode)
	if data != nil {
		err := json.NewEncoder(*respWriter).Encode(data)
		if err != nil {
			(*respWriter).WriteHeader(500)
		}
	}
}
