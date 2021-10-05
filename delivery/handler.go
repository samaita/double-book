package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/samaita/double-book/model"
)

func handleGetUserInfo(c *gin.Context) {
	var (
	// err error
	)

	UID := c.GetString(model.CtxUID)

	// if err = newUser.getUserInfo(); err != nil {
	// 	log.Printf("[handleGetUserInfo][getUserInfo] Input: %v, Output %v", newUser.UID, err.Error())
	// 	APIResponseInternalServerError(c, nil, err.Error())
	// 	return
	// }

	response := map[string]interface{}{
		"uid": UID,
	}

	APIResponseOK(c, response)
}
