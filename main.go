package main

import (
	"fmt"
	. "github.com/waltermblair/brain/brain"
	"log"
	"os"
	"time"
)

// Creates rabbit client with queue specified by env variable. Creates processor and runs services.
func main() {

	for {
		rabbitHost := os.Getenv("RABBIT_HOST")
		thisQueue := os.Getenv("THIS_QUEUE")

		fmt.Println("RABBIT HOST: ", rabbitHost)
		fmt.Println("THIS QUEUE: " , thisQueue)

		rabbit := NewRabbitClient(
			fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitHost),
			thisQueue)
		db, err := NewDBClient("root:root@tcp(localhost:3306)/store")
		service, err := NewService(db)

		if err != nil {
			log.Println("failed to initialize service: ", err)
		} else {
			log.Println("initialized service successfully, starting API")
			RunAPI(service, rabbit, db)
		}

		log.Println("attempting to re-initialize service")
		time.Sleep(2 * time.Second)
	}

}