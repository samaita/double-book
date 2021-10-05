package delivery

import (
	"github.com/gin-gonic/gin"
)

func InitHandler() {
	router := gin.Default()

	authorized := router.Group("/api/flashsale/")
	authorized.Use()
	{
		authorized.GET("/list", handleGetFlashSaleList)
	}

	router.Run(":3000")
}
