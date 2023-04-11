package main

import(
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/go-worker-golden-data-rmq/internal/adapter/notification"
	"github.com/go-worker-golden-data-rmq/internal/repository/db_postgre"
	"github.com/go-worker-golden-data-rmq/internal/core"
	"github.com/go-worker-golden-data-rmq/internal/service"
)

var (
	logLevel		=	zerolog.DebugLevel // InfoLevel DebugLevel
	version			=	"go-worker-golden-data-rmq"

	configRabbitMQ 		core.ConfigRabbitMQ
	envDB	 			core.DatabaseRDS
	dataBaseHelper 		db_postgre.DatabaseHelper
	repoDB				db_postgre.WorkerRepository
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)

	configRabbitMQ.User = "guest"
	configRabbitMQ.Password = "guest"
	configRabbitMQ.Port = "localhost:5672/"
	configRabbitMQ.QueueName = "queue_person_quorum"
	configRabbitMQ.TimeDeleyQueue = 500

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
		configRabbitMQ.User = os.Getenv("RMQ_USER")
	}
	if os.Getenv("RMQ_PASS") !=  "" {
		configRabbitMQ.Password = os.Getenv("RMQ_PASS")
	}
	if os.Getenv("RMQ_PORT") !=  "" {
		configRabbitMQ.Port = os.Getenv("RMQ_PORT")
	}
	if os.Getenv("RMQ_QUEUE") !=  "" {
		configRabbitMQ.QueueName = os.Getenv("RMQ_QUEUE")
	}
	if os.Getenv("TIME_DELAY_QUEUE") !=  "" {
		intVar, _ := strconv.Atoi(os.Getenv("TIME_DELAY_QUEUE"))
		configRabbitMQ.TimeDeleyQueue = intVar
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
	log.Debug().Str("version", version).
				Msg("Enviroment Variables")
	log.Debug().Str("configRabbitMQ.User: ", configRabbitMQ.User).
				Msg("-----")
	log.Debug().Str("configRabbitMQ.Password: ", configRabbitMQ.Password).
				Msg("-----")
	log.Debug().Str("configRabbitMQ.Port :", configRabbitMQ.Port).
				Msg("----")
	log.Debug().Str("configRabbitMQ.QueueName :", configRabbitMQ.QueueName).
				Msg("----")
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

	repoDB 			= db_postgre.NewWorkerRepository(dataBaseHelper)
	workerService 	:= service.NewWorkerService(&repoDB)

	consumer, err := notification.NewConsumerService(&configRabbitMQ, workerService)
	if err != nil {
		log.Error().Err(err).Msg("error create notification.NewConsumerService")
		os.Exit(3)
	}

	consumer.ConsumerQueue()
}