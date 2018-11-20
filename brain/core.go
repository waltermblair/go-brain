package brain

import (
	"log"
	"strconv"
)

type Service interface {
	GetConfig() Config
	FetchComponentConfig(Config, DBClient) (Config, error)
	BuildInputMessage(bool) MessageBody
	BuildConfigMessage(Config, DBClient) MessageBody
	RunDemo(MessageBody, RabbitClient, DBClient) (bool, error)
}

type ServiceImpl struct {
	config		Config
}

func NewService(db DBClient) (Service, error) {
	cfg := Config{
		0,
		"",
		"",
		0,
		[]int{0},
	}

	numInputs, nextKeys, _, err := db.FetchConfig(cfg)

	if err != nil {
		log.Println("error creating new service: ", err.Error())
		return nil, err
	}

	cfg.NumInputs = numInputs
	cfg.NextKeys = nextKeys

	s := ServiceImpl{
		cfg,
	}
	return &s, err
}

func (s *ServiceImpl) GetConfig() Config {
	return s.config
}

func (s *ServiceImpl) FetchComponentConfig(config Config, db DBClient) (Config, error) {
	log.Println("fetching component config for routing key: ", config.ID)
	numInputs, nextKeys, fn, err := db.FetchConfig(config)

	if err != nil {
		log.Println("error fetching component config")
		return Config{}, nil
	}

	return Config{
		config.ID,
		config.Status,
		fn,
		numInputs,
		nextKeys,
	}, nil
}

func (s *ServiceImpl) BuildInputMessage(input bool) MessageBody {

	return MessageBody{
		Input: []bool{input},
	}

}

func (s *ServiceImpl) BuildConfigMessage(config Config, db DBClient) MessageBody {

	config, _ = s.FetchComponentConfig(config, db)
	msgBody := MessageBody{
		Configs: []Config{config},
	}

	return msgBody
}

// TODO - refactor DRY
// RunDemo takes the entire GUI message body containing the user-selected configuration and inputs and sends one message to each component that contains both a config and also an input if necessary.
func (s *ServiceImpl) RunDemo(body MessageBody, rabbit RabbitClient, db DBClient) (output bool, err error){

	configs := body.Configs
	log.Println("number of messages to send: ", len(configs))

	//	build and publish each config message
	for _, config := range configs {
		msg := s.BuildConfigMessage(config, db)
		// determine routing key
		nextQueue := strconv.Itoa(config.ID)
		log.Println("sending this config message: ", msg)
		err = rabbit.Publish(msg, nextQueue)
	}

	// build and publish each direct input message to nextKeys
	for i, nextKey := range s.config.NextKeys {
		input := body.Input[i]
		msg := s.BuildInputMessage(input)
		log.Println("sending this input message: ", msg)
		nextQueue := strconv.Itoa(nextKey)
		err = rabbit.Publish(msg, nextQueue)
	}

	log.Println("waiting for output...")
	output, err = rabbit.RunConsumer()
	log.Println("received output: ", output)
	log.Println("received err: ", err)

	return output, err

}

