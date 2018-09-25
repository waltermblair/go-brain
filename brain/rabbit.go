package brain

import (
	"encoding/json"
	"fmt"
	"github.com/assembla/cony"
	"github.com/streadway/amqp"
	"log"
)

type RabbitClient interface {
	RunConsumer() (bool, error)
	Publish(MessageBody, string) error
	InitRabbit()
}

type RabbitClientImpl struct {
	URL       string
	exc       cony.Exchange
	cli       *cony.Client
	cns       *cony.Consumer
	cnsQue    *cony.Queue
	cnsBnd    cony.Binding
	pbl       *cony.Publisher
	pblQue    *cony.Queue
	pblBnd    cony.Binding
	thisQueue string
	nextQueue string
}

func NewRabbitClient(url string, thisQueue string) RabbitClient {

	r := RabbitClientImpl{
		URL:       url,
		thisQueue: thisQueue,
	}
	r.InitRabbit()

	fmt.Println("Initialized rabbit client at ", r.URL)

	return &r

}

// TODO - ack/delete message?
// Used to run callback queue for returning final output to UI
func (r *RabbitClientImpl) RunConsumer() (res bool, err error){

	cli := cony.NewClient(
		cony.URL(r.URL),
		cony.Backoff(cony.DefaultBackoff),
	)

	cli.Declare([]cony.Declaration{
		cony.DeclareQueue(r.cnsQue),
		cony.DeclareExchange(r.exc),
		cony.DeclareBinding(r.cnsBnd),
	})

	// Declare and register a consumer
	cns := cony.NewConsumer(r.cnsQue)

	cli.Consume(cns)
	defer cli.Close()

	for cli.Loop() {
		var res MessageBody

		select {
		case msg := <-cns.Deliveries():
			log.Printf("Received body: %q\n", msg.Body)
			json.Unmarshal(msg.Body, &res)
			return res.Input[0], nil
			msg.Ack(false)
		case err := <-cns.Errors():
			fmt.Printf("Consumer error: %v\n", err)
			return false, err
		case err := <-cli.Errors():
			fmt.Printf("Client error: %v\n", err)
			return false, err
		}
	}

	return false, err
}

func (r *RabbitClientImpl) Publish(body MessageBody, nextQueue string) error {

	cli := cony.NewClient(
		cony.URL(r.URL),
		cony.Backoff(cony.DefaultBackoff),
	)

	r.pblQue = &cony.Queue{
		AutoDelete: false,
		Name:       nextQueue,
		Durable:	true,
	}
	r.pblBnd = cony.Binding{
		Queue:    r.pblQue,
		Exchange: r.exc,
		Key:      nextQueue,
	}

	cli.Declare([]cony.Declaration{
		cony.DeclareQueue(r.pblQue),
		cony.DeclareExchange(r.exc),
		cony.DeclareBinding(r.pblBnd),
	})

	pbl := cony.NewPublisher(r.exc.Name, nextQueue)
	cli.Publish(pbl)

	go func() {
		for cli.Loop() {
			select {
			case err := <-cli.Errors():
				fmt.Println(err)
			}
		}
	}()

	bytes, err := json.Marshal(body)

	if err != nil {
		fmt.Printf("Error unmarshaling MessageBody: %v\n", err)
	}

	go func() {
		err = pbl.Publish(amqp.Publishing{
			Body: bytes,
		})
		if err != nil {
			fmt.Printf("Client publish error: %v\n", err)
		}
	}()

	return err

}

func (r *RabbitClientImpl) InitRabbit() {

	r.exc = cony.Exchange{
		Name:       "myExc",
		Kind:       "topic",
		AutoDelete: false,
		Durable:	true,
	}
	r.cnsQue = &cony.Queue{
		AutoDelete: false,
		Name:       r.thisQueue,
		Durable:	true,
	}
	r.cnsBnd = cony.Binding{
		Queue:    r.cnsQue,
		Exchange: r.exc,
		Key:      r.thisQueue,
	}

}