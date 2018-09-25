package main

import (
	. "github.com/waltermblair/brain/brain"
	"log"
)

func main() {

	// TODO - replace with env
	rabbit := NewRabbitClient("amqp://guest:guest@localhost:5672/", "0")
	db, err := NewDBClient("root:root@tcp(localhost:3306)/store")
	if err != nil {
		log.Fatal("failed to initialize service: ", err)
	} else {
		log.Println("initialized service successfully")
	}
	RunAPI(rabbit, db)

}