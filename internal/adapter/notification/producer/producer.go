package producer

import (
	"time"
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"

	"github.com/go-worker-golden-data-rmq/internal/core"
	//"github.com/go-worker-golden-data-rmq/internal/service"

)

var (
	childLogger = log.With().Str("notification", "rabbitMQ").Logger()
	consumer_name = "consumer.person"
)

type RMQNotification struct{
	rmqConnection 	*amqp.Connection
	configRabbitMQ 	*core.ConfigRabbitMQ
	//workerService 	*service.WorkerService
}

func NewRMQNotification(configRabbitMQ *core.ConfigRabbitMQ,
						//workerService  *service.WorkerService,
						) (*RMQNotification, error){
	childLogger.Debug().Msg("NewRMQNotification")

	rabbitmqURL := "amqp://" + configRabbitMQ.User + ":" + configRabbitMQ.Password + "@" + configRabbitMQ.Port
	childLogger.Debug().Str("rabbitmqURL :", rabbitmqURL).Msg("Rabbitmq URI")
	
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		childLogger.Error().Err(err).Msg("error connect to server message") 
		return nil, err
	}
	return &RMQNotification{
		rmqConnection: conn,
		configRabbitMQ: configRabbitMQ,
		//workerService: 	workerService,
	}, nil
}

func (n *RMQNotification) ProducerQueue(msg *core.WebHook) error {
	childLogger.Debug().Msg("ProducerQueue")

	ch, err := n.rmqConnection.Channel()
	if err != nil {
		childLogger.Error().Err(err).Msg("error channel the server message") 
		return err
	}
	defer ch.Close()

	args := amqp.Table{ // queue args
		amqp.QueueTypeArg: amqp.QueueTypeQuorum,
	}
	q, err := ch.QueueDeclare(n.configRabbitMQ.QueueName, // name
								true,         // durable
								false,        // delete when unused
								false,        // exclusive
								false,        // no-wait
								args,          // arguments
	)
	if err != nil {
		childLogger.Error().Err(err).Msg("error declare queue !!!") 
		return err
	}

	//person_mock := p.CreateDataMock(i,my_ip)
	body, _ := json.Marshal(msg)

	payloadMsg := amqp.Publishing{	
							ContentType:  "application/json",
							Timestamp:    time.Now(),
							DeliveryMode: amqp.Persistent,
							Body:         []byte(body),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
								"", // exchange
								q.Name, // routing key
								false,  // mandatory
								false,  // immediate
								payloadMsg)
	if err != nil {
		childLogger.Error().Err(err).Msg("error publish message") 
		return err
	}

	childLogger.Debug().Str("msg :", string(body)).Msg("Success Publish a message (ProducerQueue)")

	return nil	
}