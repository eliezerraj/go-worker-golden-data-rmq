package notification

import (
	"time"

	"github.com/rs/zerolog/log"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/go-worker-golden-data-rmq/internal/core"
	"github.com/go-worker-golden-data-rmq/internal/service"

)

var (
	childLogger = log.With().Str("notification", "rabbitMQ").Logger()
	consumer_name = "consumer.person"
)

type ConsumerService struct{
	consumer 		*amqp.Connection
	configRabbitMQ 	*core.ConfigRabbitMQ
	workerService	*service.WorkerService
}

func NewConsumerService(configRabbitMQ *core.ConfigRabbitMQ,
						workerService  *service.WorkerService) (*ConsumerService, error){
	childLogger.Debug().Msg("NewConsumerService")

	rabbitmqURL := "amqp://" + configRabbitMQ.User + ":" + configRabbitMQ.Password + "@" + configRabbitMQ.Port
	childLogger.Debug().Str("rabbitmqURL :", rabbitmqURL).Msg("Rabbitmq URI")
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		childLogger.Error().Err(err).Msg("error connect to server message") 
		return nil, err
	}
	return &ConsumerService{
		consumer: 		conn,
		configRabbitMQ: configRabbitMQ,
		workerService: 	workerService,
	}, nil
}

func (c *ConsumerService) ConsumerQueue() error {
	childLogger.Debug().Msg("ConsumerQueue")

	ch, err := c.consumer.Channel()
	if err != nil {
		childLogger.Error().Err(err).Msg("error channel the server message") 
		return err
	}
	defer ch.Close()

	args := amqp.Table{ // queue args
		amqp.QueueTypeArg: amqp.QueueTypeQuorum,
	}
	q, err := ch.QueueDeclare(	c.configRabbitMQ.QueueName, // name
								true,         // durable
								false,        // delete when unused
								false,        // exclusive
								false,        // no-wait
								args,          // arguments
	)
	if err != nil {
		childLogger.Error().Err(err).Msg("error declare queue !!!!") 
		return err
	}

	msgs, err := ch.Consume(	q.Name, // queue
								consumer_name,    // consumer
								true,   // auto-ack
								false,  // exclusive
								false,  // no-local
								false,  // no-wait
								nil,    // args
	)
	if err != nil {
		childLogger.Error().Err(err).Msg("error consume message") 
		return err
	}
	
	var forever chan struct{}

	childLogger.Debug().Msg("Starting Consumer...")
	go func() {
		for d := range msgs {
			childLogger.Debug().Msg("++++++++++++++++++++++++++++")
			childLogger.Debug().Str("msg.Body:", string(d.Body)).Msg(" Success Receive a message (ConsumerQueue) ") 
			
			c.workerService.DataEnrichment()

			time.Sleep(time.Duration(c.configRabbitMQ.TimeDeleyQueue) * time.Millisecond)
		}
	}()
	<-forever

	return nil
}