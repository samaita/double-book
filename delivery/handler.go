package delivery

import (
	"log"
	"strconv"
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

func handleAddToCart(c *gin.Context) {
	var (
		err  error
		date string
		data usecase.ATCData
	)

	userIDString := c.GetString(CTX_USERID)
	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		log.Printf("[handleAddToCart][ParseInt] Input: %s, Output %s", userIDString, err.Error())
		APIResponseBadRequest(c, nil, err.Error())
		return
	}

	productIDString := c.Request.FormValue("product_id")
	productID, err := strconv.ParseInt(productIDString, 10, 64)
	if err != nil {
		log.Printf("[handleAddToCart][ParseInt] Input: %s, Output %s", productIDString, err.Error())
		APIResponseBadRequest(c, nil, err.Error())
		return
	}

	flashSaleIDString := c.Request.FormValue("flashsale_id")
	flashSaleID, err := strconv.ParseInt(flashSaleIDString, 10, 64)
	if err != nil {
		log.Printf("[handleAddToCart][ParseInt] Input: %s, Output %s", flashSaleIDString, err.Error())
		APIResponseBadRequest(c, nil, err.Error())
		return
	}

	amountString := c.Request.FormValue("amount")
	amount, err := strconv.Atoi(amountString)
	if err != nil {
		log.Printf("[handleAddToCart][ParseInt] Input: %s, Output %s", amountString, err.Error())
		APIResponseBadRequest(c, nil, err.Error())
		return
	}

	if data, err = usecase.HandleAddToCart(userID, productID, flashSaleID, amount); err != nil {
		log.Printf("[handleGetUserInfo][getUserInfo] Input: %s, Output %s", date, err.Error())
		APIResponseInternalServerError(c, nil, err.Error())
		return
	}

	response := map[string]interface{}{
		"data": data,
	}

	APIResponseOK(c, response)
}
