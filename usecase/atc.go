package usecase

import (
	"errors"
	"log"

	"github.com/samaita/double-book/repository"
)

type ATCData struct {
	ProductID  int64 `json:"product_id"`
	SuccessATC bool  `json:"success_atc"`
	Amount     int   `json:"amount"`
}

func HandleAddToCart(userID, productID, flashSaleID int64, amount int) (ATCData, error) {
	var (
		data ATCData
		err  error
	)

	data.ProductID = productID
	data.Amount = amount

	if flashSaleID > 0 && amount > 1 {
		err = errors.New("tidak bisa membeli barang flash sale lebih dari 1")
		return data, err
	}

	if err = repository.Publish("FLASH_SALE", ATCPayload{
		UserID:      userID,
		ProductID:   productID,
		FlashSaleID: flashSaleID,
		Amount:      amount,
		IPOrigin:    repository.IPAddr,
	}); err != nil {
		log.Printf("[HandleAddToCart][Publish] Input: %+v Output: %v", data, err)
		return data, err
	}

	// WAIT

	return data, nil
}
