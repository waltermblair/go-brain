package main

import (
	. "github.com/waltermblair/brain/brain"
	"log"
	"os"
)

// Creates rabbit client with queue specified by env variable. Creates processor and runs services.
func main() {

	rabbit := NewRabbitClient("amqp://guest:guest@localhost:5672/", os.Getenv("THIS_QUEUE"))
	db, err := NewDBClient("root:root@tcp(localhost:3306)/store")
	service := NewService(db)
	if err != nil {
		log.Fatal("failed to initialize service: ", err)
	} else {
		log.Println("initialized service successfully")
	}
	RunAPI(service, rabbit, db)

}