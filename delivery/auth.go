package delivery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samaita/double-book/model"
)

const (
	CTX_USERID = "USER_ID"
)

// IsLoggedIn mock an user logged in state, everyone will be treated as has login as long header filled with valid userID
func IsLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Request.Header.Get("USER_ID")
		userIDInt, _ := strconv.ParseInt(userID, 10, 64)

		defaultError := APIDefaultErrResponse()
		defaultError["message_error"] = "anda tidak memiliki akses"

		if userIDInt > 0 {
			var err error
			user := model.NewUser(userIDInt)
			if err = user.Load(); err == nil {
				c.Set(CTX_USERID, c.Request.Header.Get("USER_ID"))
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, defaultError)
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, defaultError)
		}

	}
}
