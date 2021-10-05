package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/samaita/double-book/repository"
)

const (
	TOPIC_FLASH_SALE      = "FLASH_SALE"
	TOPIC_FLASH_SALE_HOOK = "FLASH_SALE_HOOK"
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

	// ATC

	if err = repository.Publish("FLASH_SALE_HOOK", ATCNSQPayload{
		UserID:      data.UserID,
		ProductID:   data.ProductID,
		FlashSaleID: data.FlashSaleID,
		Amount:      data.Amount,
		IPOrigin:    data.IPOrigin,
		Timestamp:   data.Timestamp,
		Success:     1,
	}); err != nil {
		log.Printf("[HandleAddToCartNSQ][Publish] Input: %+v Output: %v", data, err)
		return
	}
}
