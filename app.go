package main

import (
	"log"

	"github.com/samaita/double-book/delivery"
	"github.com/samaita/double-book/repository"
	"github.com/samaita/double-book/usecase"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	repository.InitDB("sqlite3", "db_poc.db")
	repository.InitPublisher()
	repository.GetIP()

	go usecase.InitConsumerFlashSale()

	delivery.InitHandler()
}
