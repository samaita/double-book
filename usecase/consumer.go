package usecase

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
)

const (
	TOPIC_FLASH_SALE = "FLASH_SALE"
)

type ATCPayload struct {
	IPOrigin    string `json:"ip_origin"`
	FlashSaleID int64  `json:"flashsale_id"`
	ProductID   int64  `json:"product_id"`
	UserID      int64  `json:"user_id"`
	Amount      int    `json:"amount"`
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
			log.Printf("Got a message: string(%v)", string(message.Body))
			return nil
		}))
		err := q.ConnectToNSQD("127.0.0.1:4150")
		if err != nil {
			log.Panic("Could not connect")
		}
	}

	wg.Wait()
}
