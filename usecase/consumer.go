package usecase

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/samaita/double-book/model"
	"github.com/samaita/double-book/repository"
)

const (
	TOPIC_FLASH_SALE      = "FLASH_SALE"
	TOPIC_FLASH_SALE_HOOK = "FLASH_SALE_HOOK"

	REQUEUE = 1
	FINISH  = 0
)

type ATCNSQPayload struct {
	IPOrigin    string `json:"ip_origin"`
	FlashSaleID int64  `json:"flashsale_id"`
	ProductID   int64  `json:"product_id"`
	UserID      int64  `json:"user_id"`
	Amount      int    `json:"amount"`
	Timestamp   int64  `json:"timestamp"`
	Success     int    `json:"success,omitempty"`
}

func InitConsumerFlashSale() {

	listFlashSaleToday, err := HandleGetFlashSaleByDate(time.Now().Format("2006-01-02"))
	if err != nil {
		log.Fatalln("[InitConsumerFlashSale][HandleGetFlashSaleByDate] Consumer cant't start, err:", err)
	}

	var listChannelToCreate []string
	for _, flashSaleToday := range listFlashSaleToday {
		for _, product := range flashSaleToday.ProductList {
			listChannelToCreate = append(listChannelToCreate, fmt.Sprintf("flashsale_%d_product_%d", flashSaleToday.FlashSale.FlashSaleID, product.ProductID))
		}
	}

	wg := &sync.WaitGroup{}

	for _, channelToCreate := range listChannelToCreate {
		wg.Add(1)

		config := nsq.NewConfig()
		q, _ := nsq.NewConsumer(TOPIC_FLASH_SALE, channelToCreate, config)
		q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
			HandleAddToCartNSQ(message.Body)
			return nil
		}))
		err := q.ConnectToNSQD("127.0.0.1:4150")
		if err != nil {
			log.Panic("Could not connect")
		}
	}

	wg.Wait()
}

func InitConsumerFlashSaleHook() {
	// Use service discover to populate the list
	listInstance := []string{repository.IPAddr}

	var listChannelToCreate []string
	for _, instance := range listInstance {
		listChannelToCreate = append(listChannelToCreate, fmt.Sprintf("%s", instance))
	}

	wg := &sync.WaitGroup{}

	for _, channelToCreate := range listChannelToCreate {
		wg.Add(1)

		config := nsq.NewConfig()
		q, _ := nsq.NewConsumer(TOPIC_FLASH_SALE_HOOK, channelToCreate, config)
		q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
			HandleFlashSaleHookNSQ(message.Body)
			return nil
		}))
		err := q.ConnectToNSQD("127.0.0.1:4150")
		if err != nil {
			log.Panic("Could not connect")
		}
	}

	wg.Wait()
}

func HandleFlashSaleHookNSQ(b []byte) {
	var (
		data ATCNSQPayload
		err  error
	)

	if err = json.Unmarshal(b, &data); err != nil {
		log.Printf("[HandleFlashSaleHookNSQ][Unmarshal] Input: %v, Output %s", string(b), err.Error())
		return
	}

	if data.IPOrigin != repository.IPAddr {
		return
	}

	keyGoChannel := fmt.Sprintf("%d-%d-%d-%d", data.UserID, data.FlashSaleID, data.ProductID, data.Timestamp)
	repository.MapGoChannel[keyGoChannel] = data.Success

}

func HandleAddToCartNSQ(b []byte) {
	var (
		data ATCNSQPayload
		err  error
	)

	if err = json.Unmarshal(b, &data); err != nil {
		log.Printf("[HandleAddToCartNSQ][Unmarshal] Input: %v, Output %s", string(b), err.Error())
		return
	}

	product := model.NewProduct(data.ProductID)
	if err = product.GetStock(); err != nil {
		log.Printf("[HandleAddToCartNSQ][Cart][LoadByUser] Input: %v, Output %s", data.UserID, err.Error())
		return
	}

	if product.Remaining <= 0 {
		return
	}

	cart := model.NewCart(0)
	if err = cart.LoadByUser(data.UserID); err != nil && err != sql.ErrNoRows {
		log.Printf("[HandleAddToCartNSQ][Cart][LoadByUser] Input: %v, Output %s", data.UserID, err.Error())
		return
	}

	if err == sql.ErrNoRows {
		if err = cart.Create(data.UserID, model.StatusCartActive); err != nil {
			log.Printf("[HandleAddToCartNSQ][Cart][Create] Input: %v, Output %s", data.UserID, err.Error())
			return
		}
	}

	isExist, err := cart.IsExist(data.ProductID)
	if err != nil || isExist {
		return
	}

	if err = cart.Add(data.ProductID, data.Amount); err == nil {
		data.Success = model.StatusSuccessATC
	}

	if err = repository.Publish("FLASH_SALE_HOOK", ATCNSQPayload{
		UserID:      data.UserID,
		ProductID:   data.ProductID,
		FlashSaleID: data.FlashSaleID,
		Amount:      data.Amount,
		IPOrigin:    data.IPOrigin,
		Timestamp:   data.Timestamp,
		Success:     data.Success,
	}); err != nil {
		log.Printf("[HandleAddToCartNSQ][Publish] Input: %+v Output: %v", data, err)
		return
	}

	return
}
