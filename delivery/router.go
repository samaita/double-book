package delivery

import (
	"github.com/gin-gonic/gin"
)

func InitHandler() {
	router := gin.Default()

	flashsale := router.Group("/api/flashsale/")
	flashsale.Use()
	{
		flashsale.GET("/list", handleGetFlashSaleList)
	}

	cart := router.Group("/api/cart/")
	cart.Use(IsLoggedIn())
	{
		cart.POST("/add", handleAddToCart)
		// atc.GET("/checkout", handleCheckout)
	}

	router.Run(":3000")
}
