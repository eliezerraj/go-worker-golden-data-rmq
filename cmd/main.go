package main

import(
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/go-worker-golden-data-rmq/internal/adapter/notification/consumer"
	"github.com/go-worker-golden-data-rmq/internal/adapter/notification/producer"
	"github.com/go-worker-golden-data-rmq/internal/repository/db_postgre"
	"github.com/go-worker-golden-data-rmq/internal/core"
	"github.com/go-worker-golden-data-rmq/internal/service"
)

var (
	logLevel		=	zerolog.DebugLevel // InfoLevel DebugLevel
	version			=	"go-worker-golden-data-rmq"

	configRMQConsumer 	core.ConfigRabbitMQ
	configRMQProducer	core.ConfigRabbitMQ
	envDB	 			core.DatabaseRDS
	dataBaseHelper 		db_postgre.DatabaseHelper
	repoDB				db_postgre.WorkerRepository
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)

	configRMQConsumer.User = "guest"
	configRMQConsumer.Password = "guest"
	configRMQConsumer.Port = "localhost:5672/"
	configRMQConsumer.QueueName = "queue_person_quorum"
	configRMQConsumer.TimeDeleyQueue = 500

	configRMQProducer.User = "guest"
	configRMQProducer.Password = "guest"
	configRMQProducer.Port = "localhost:5672/"
	configRMQProducer.QueueName = "queue_webhook_quorum"
	configRMQProducer.TimeDeleyQueue = 500

	envDB.Host = "host.docker.internal"
	envDB.Port = "5432"
	envDB.Schema = "public"
	envDB.DatabaseName = "postgres"
	envDB.User  = "admin"
	envDB.Password  = "admin"
	envDB.Db_timeout = 90
	envDB.Postgres_Driver = "postgres"

	getEnv()
}

func getEnv() {
	if os.Getenv("LOG_LEVEL") !=  "" {
		if (os.Getenv("LOG_LEVEL") == "DEBUG"){
			logLevel = zerolog.DebugLevel
		}else if (os.Getenv("LOG_LEVEL") == "INFO"){
			logLevel = zerolog.InfoLevel
		}else if (os.Getenv("LOG_LEVEL") == "ERROR"){
				logLevel = zerolog.ErrorLevel
		}else {
			logLevel = zerolog.DebugLevel
		}
	}
	if os.Getenv("VERSION") !=  "" {
		version = os.Getenv("VERSION")
	}
	if os.Getenv("RMQ_USER") !=  "" {
		configRMQConsumer.User = os.Getenv("RMQ_USER")
	}
	if os.Getenv("RMQ_PASS") !=  "" {
		configRMQConsumer.Password = os.Getenv("RMQ_PASS")
	}
	if os.Getenv("RMQ_PORT") !=  "" {
		configRMQConsumer.Port = os.Getenv("RMQ_PORT")
	}
	if os.Getenv("RMQ_QUEUE") !=  "" {
		configRMQConsumer.QueueName = os.Getenv("RMQ_QUEUE")
	}
	if os.Getenv("TIME_DELAY_QUEUE") !=  "" {
		intVar, _ := strconv.Atoi(os.Getenv("TIME_DELAY_QUEUE"))
		configRMQConsumer.TimeDeleyQueue = intVar
	}

	if os.Getenv("DB_HOST") !=  "" {
		envDB.Host = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") !=  "" {
		envDB.Port = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_USER") !=  "" {
		envDB.User = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PASSWORD") !=  "" {	
		envDB.Password = os.Getenv("DB_PASSWORD")
	}
	if os.Getenv("DB_NAME") !=  "" {	
		envDB.DatabaseName = os.Getenv("DB_NAME")
	}
	if os.Getenv("DB_SCHEMA") !=  "" {	
		envDB.Schema = os.Getenv("DB_SCHEMA")
	}

	if os.Getenv("DB_DRIVER") !=  "" {	
		envDB.Postgres_Driver = os.Getenv("DB_DRIVER")
	}
	if os.Getenv("DB_TIMEOUT") !=  "" {	
		intVar, _ := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
		envDB.Db_timeout = intVar
	}

}

func main()  {
	log.Debug().Msg("main")
	log.Debug().Msg("-------------------")
	log.Debug().Str("version", version).Msg("Enviroment Variables")
	log.Debug().Str("configRabbitMQ.User: ", configRMQConsumer.User).Msg("")
	log.Debug().Str("configRabbitMQ.Password: ", configRMQConsumer.Password).Msg("")
	log.Debug().Str("configRabbitMQ.Port :", configRMQConsumer.Port).Msg("")
	log.Debug().Str("configRabbitMQ.QueueName :", configRMQConsumer.QueueName).Msg("")
	log.Debug().Msg("--------------------")
	
	count := 1
	var err error
	for {
		dataBaseHelper, err = db_postgre.NewDatabaseHelper(envDB)
		if err != nil {
			if count < 3 {
				log.Error().Err(err).Msg("Erro na abertura do Database")
			} else {
				log.Error().Err(err).Msg("EERRO FATAL na abertura do Database aborting")
				panic(err)	
			}
			time.Sleep(3 * time.Second)
			count = count + 1
			continue
		}
		break
	}
	repoDB = db_postgre.NewWorkerRepository(dataBaseHelper)

	producer, err := producer.NewRMQNotification(&configRMQProducer)
	if err != nil {
		log.Error().Err(err).Msg("error create producer notification.NewRMQNotification")
		os.Exit(3)
	}
	workerService := service.NewWorkerService(&repoDB, producer)

	consumer, err := consumer.NewConsumerService(&configRMQConsumer, workerService)
	if err != nil {
		log.Error().Err(err).Msg("error create notification.NewConsumerService")
		os.Exit(3)
	}

	consumer.ConsumerQueue()

}