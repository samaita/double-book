package delivery

import (
	"github.com/gin-gonic/gin"
)

func InitHandler() {
	router := gin.Default()

	authorized := router.Group("/api")
	authorized.Use()
	{
		authorized.GET("/user/info", handleGetUserInfo)
	}

	router.Run(":3000")
}
