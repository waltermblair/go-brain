package brain

import (
	"log"
	"strconv"
)

type Service interface {
	FetchComponentConfig(Config, DBClient) (Config, error)
	SelectInput([]bool, Config) bool
	BuildMessage([]bool, Config, DBClient) MessageBody
	RunDemo(MessageBody, RabbitClient, DBClient) (bool, error)
}

type ServiceImpl struct {
	config		Config
}

func NewService(db DBClient) (Service, error) {
	nextKeys, _, err := db.FetchConfig(0)

	if err != nil {
		log.Println("error creating new service")
		return nil, err
	}
	cfg := Config {
		0,
		"",
		"",
		nextKeys,
	}
	s := ServiceImpl{
		cfg,
	}
	return &s, err
}

func (s *ServiceImpl) FetchComponentConfig(config Config, db DBClient) (Config, error) {
	log.Println("fetching component config for routing key: ", config.ID)
	nextKeys, fn, err := db.FetchConfig(config.ID)

	if err != nil {
		log.Println("error fetching component config")
		return Config{}, nil
	}

	return Config{
		config.ID,
		config.Status,
		fn,
		nextKeys,
	}, nil
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

	config, _ = s.FetchComponentConfig(config, db)
	input := s.SelectInput(inputs, config)

	return MessageBody{
		Configs: []Config{config},
		Input: []bool{input},
	}
}

func (s *ServiceImpl) RunDemo(body MessageBody, rabbit RabbitClient, db DBClient) (output bool, err error){

	configs := body.Configs
	log.Println("number of messages to send: ", len(configs))

	//	build and publish each message
	for _, config := range configs {
		msg := s.BuildMessage(body.Input, config, db)
		// determine routing key
		nextQueue := strconv.Itoa(config.ID)
		log.Println("sending this message: ", msg)
		err = rabbit.Publish(msg, nextQueue)
	}

	log.Println("waiting for output...")
	output, err = rabbit.RunConsumer()
	log.Println("received output: ", output)
	log.Println("received err: ", err)

	return output, err

}

