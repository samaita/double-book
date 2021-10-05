package repository

import (
	"encoding/json"
	"log"

	"github.com/nsqio/go-nsq"
)

var Producer *nsq.Producer
var MapGoChannel map[string]int

func InitPublisher() {
	var (
		err error
	)

	config := nsq.NewConfig()
	Producer, err = nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal("[Producer][MQ] Err:", err)
	}

	MapGoChannel = make(map[string]int)
}

func Publish(topic string, payload interface{}) error {
	var (
		err     error
		message []byte
	)

	if message, err = json.Marshal(payload); err != nil {
		log.Printf("[Producer][Marshal] Input: %s Output: %v", topic, err)
		return err
	}

	err = Producer.Publish(topic, message)
	if err != nil {
		log.Printf("[Producer][Publish] Input: %s Output: %v", topic, err)
		return err
	}

	return nil
}
