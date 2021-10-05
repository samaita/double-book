package main

import (
	"github.com/samaita/double-book/delivery"
	"github.com/samaita/double-book/repository"
)

func main() {
	repository.InitDB("sqlite3", "db_poc.db")
	delivery.InitHandler()
}
