package usecase

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/samaita/double-book/model"
	"github.com/samaita/double-book/repository"
)

type ATCData struct {
	ProductID  int64 `json:"product_id"`
	SuccessATC bool  `json:"success_atc"`
	Amount     int   `json:"amount"`
}

func HandleAddToCart(userID, productID, flashSaleID int64, amount int) (ATCData, error) {
	var (
		data    ATCData
		err     error
		success int
	)

	data.ProductID = productID
	data.Amount = amount

	if flashSaleID > 0 && amount > 1 {
		err = errors.New("tidak bisa membeli barang flash sale lebih dari 1")
		return data, err
	}

	timestamp := time.Now().Unix()
	if err = repository.Publish("FLASH_SALE", ATCNSQPayload{
		UserID:      userID,
		ProductID:   productID,
		FlashSaleID: flashSaleID,
		Amount:      amount,
		IPOrigin:    repository.IPAddr,
		Timestamp:   timestamp,
	}); err != nil {
		log.Printf("[HandleAddToCart][Publish] Input: %+v Output: %v", data, err)
		return data, err
	}

	// using dumb solution: simple cacheMap as workaround as I can't make the channel work.
	t := time.Now()
	for {
		keyGoChannel := fmt.Sprintf("%d-%d-%d-%d", userID, flashSaleID, productID, timestamp)
		success = repository.MapGoChannel[keyGoChannel]
		time.Sleep(50 * time.Millisecond)
		if success != 0 || time.Since(t) > 2*time.Second {
			break
		}
	}

	data.SuccessATC = success == model.StatusSuccessATC
	return data, nil
}
