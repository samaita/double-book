package delivery

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samaita/double-book/usecase"
)

func handleGetFlashSaleList(c *gin.Context) {
	var (
		err  error
		date string
		data []usecase.FlashSaleListData
	)

	if date = c.Query("date"); date == "" {
		date = time.Now().Format("2006-01-02")
	}

	if data, err = usecase.HandleGetFlashSaleByDate(date); err != nil {
		log.Printf("[handleGetUserInfo][getUserInfo] Input: %s, Output %s", date, err.Error())
		APIResponseInternalServerError(c, nil, err.Error())
		return
	}

	response := map[string]interface{}{
		"data": data,
	}

	APIResponseOK(c, response)
}
