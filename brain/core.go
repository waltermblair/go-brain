package brain

import (
	"fmt"
	"strconv"
)

type Service interface {
	FetchComponentConfig(Config, DBClient) Config
	SelectInput([]bool, Config) bool
	BuildMessage([]bool, Config, DBClient) MessageBody
	RunDemo(MessageBody, RabbitClient, DBClient) (bool, error)
}

type ServiceImpl struct {
	config		Config
}

func NewService(db DBClient) Service {
	nextKeys, _ := db.FetchConfig(0)
	cfg := Config {
		0,
		"",
		"",
		nextKeys,
	}
	s := ServiceImpl{
		cfg,
	}
	return &s
}

func (s *ServiceImpl) FetchComponentConfig(config Config, db DBClient) Config {
	fmt.Println("fetching component config for routing key: ", config.ID)
	nextKeys, fn := db.FetchConfig(config.ID)
	return Config{
		config.ID,
		config.Status,
		fn,
		nextKeys,
	}
}

func (s *ServiceImpl) SelectInput(inputs []bool, config Config) (input bool) {

	for i, nextKey := range s.config.NextKeys {
		if config.ID == nextKey {
			input = inputs[i]
		}
	}
	return input

}

func (s *ServiceImpl) BuildMessage(inputs []bool, config Config, db DBClient) MessageBody {

	config = s.FetchComponentConfig(config, db)
	input := s.SelectInput(inputs, config)

	return MessageBody{
		Configs: []Config{config},
		Input: []bool{input},
	}
}

func (s *ServiceImpl) RunDemo(body MessageBody, rabbit RabbitClient, db DBClient) (output bool, err error){

	configs := body.Configs
	fmt.Println("number of messages to send: ", len(configs))

	//	build and publish each message
	for _, config := range configs {
		msg := s.BuildMessage(body.Input, config, db)
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

