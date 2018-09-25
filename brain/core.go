package brain

import (
	"fmt"
	"strconv"
)

func fetchComponentConfig(config Config) Config {
	id, nextKeys, fn := fetchConfig(config.ID)
	return Config{
		id,
		config.Status,
		fn,
		nextKeys,
	}
}

// TODO - only send input if key matches one of brain's nextKeys
// TODO - select which component gets which initial input
func selectInput(body MessageBody, config Config) bool {

	input := body.Input[0]
	return input

}

func buildMessage(body MessageBody, config Config) MessageBody {

	config = fetchComponentConfig(config)
	input := selectInput(body, config)

	return MessageBody{
		Configs: []Config{config},
		Input: []bool{input},
	}
}

func RunDemo(body MessageBody, rabbit RabbitClient) (output bool, err error){

	configs := body.Configs
	fmt.Println("number of messages to send: ", len(configs))

	//	build and publish each message
	for _, config := range configs {

		msg := buildMessage(body, config)

		// determine routing key
		nextQueue := strconv.Itoa(config.ID)

		fmt.Println("sending this message: ", msg)

		err = rabbit.Publish(msg, nextQueue)

	}

	fmt.Println("waiting for output...")
	output, err = rabbit.RunConsumer()
	fmt.Println("received output: ", output)
	fmt.Println("received err: ", err)

	return output, err

}

