package main

import (
	"log"

	"github.com/samaita/double-book/delivery"
	"github.com/samaita/double-book/repository"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	repository.InitDB("sqlite3", "db_poc.db")
	delivery.InitHandler()
}
