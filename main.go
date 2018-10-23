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
		mysqlHost := os.Getenv("MYSQL_HOST")

		log.Println("RABBIT & DB HOST: ", rabbitHost)
		log.Println("THIS QUEUE: " , thisQueue)

		rabbit := NewRabbitClient(
			fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitHost),
			thisQueue)
		db, err := NewDBClient(fmt.Sprintf("guest:guest@tcp(%s:3306)/store", mysqlHost))
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