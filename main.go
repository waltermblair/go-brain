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

		log.Println("RABBIT & DB HOST: ", rabbitHost)
		log.Println("THIS QUEUE: " , thisQueue)

		rabbit := NewRabbitClient(
			fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitHost),
			thisQueue)
		db, err := NewDBClient(fmt.Sprintf("root:root@tcp(%s:3306)/store", rabbitHost))
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