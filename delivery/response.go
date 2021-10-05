package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIDefaultResponse() map[string]interface{} {
	return map[string]interface{}{
		"success": true,
	}
}

func APIDefaultErrResponse() map[string]interface{} {
	return map[string]interface{}{
		"success": false,
	}
}

func APIResponseOK(c *gin.Context, data interface{}) {
	response := APIDefaultResponse()
	response["data"] = data
	c.JSON(http.StatusOK, response)
}

func APIResponseBadRequest(c *gin.Context, data interface{}, errMsg interface{}) {
	response := APIDefaultErrResponse()
	response["data"] = data
	response["message_error"] = errMsg
	c.JSON(http.StatusBadRequest, response)
}

func APIResponseUnauthorized(c *gin.Context, data interface{}, errMsg interface{}) {
	response := APIDefaultErrResponse()
	response["data"] = data
	response["message_error"] = errMsg
	c.JSON(http.StatusUnauthorized, response)
}

func APIResponseInternalServerError(c *gin.Context, data interface{}, errMsg interface{}) {
	response := APIDefaultErrResponse()
	response["data"] = data
	response["message_error"] = errMsg
	c.JSON(http.StatusInternalServerError, response)
}
