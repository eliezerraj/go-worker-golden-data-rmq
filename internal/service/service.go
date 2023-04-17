package service

import (
	"github.com/rs/zerolog/log"

	"github.com/go-worker-golden-data-rmq/internal/repository/db_postgre"
	"github.com/go-worker-golden-data-rmq/internal/core"
	"github.com/go-worker-golden-data-rmq/internal/adapter/notification/producer"

)

var childLogger = log.With().Str("service", "DataService").Logger()

type WorkerService struct {
	workerRepository 		*db_postgre.WorkerRepository
	producerRMQNotification	*producer.RMQNotification
}

func NewWorkerService(	workerRepository 		*db_postgre.WorkerRepository,
						producerRMQNotification *producer.RMQNotification ) *WorkerService{

	childLogger.Debug().Msg("NewWorkerService")
	return &WorkerService{
		workerRepository: workerRepository,
		producerRMQNotification: producerRMQNotification,
	}
}

func (s *WorkerService) DataEnrichment(id string) error{
	childLogger.Debug().Msg("--------------------------")
	childLogger.Debug().Msg("DataEnrichment")
	childLogger.Debug().Str("id : ", id).Msg("")
	childLogger.Debug().Msg("---------------------------")

	result, err := s.workerRepository.GetPerson(id)
	if err != nil {
		log.Error().Err(err).Msg("error workerRepository.GetPerson")
		return err
	}
	//log.Debug().Interface("result",result).Msg("")

	url := "https://my-webhook.com.br/" + id

	webHook := core.NewWebHook(id, result.Email, url)
	log.Debug().Interface("webHook : ",webHook).Msg("")

	err = s.producerRMQNotification.ProducerQueue(webHook)
	if err != nil {
		log.Error().Err(err).Msg("error producerRMQNotification.ProducerQueue")
		return err
	}

	err = s.workerRepository.AddWebHook(*webHook)
	if err != nil {
		log.Error().Err(err).Msg("error workerRepository.AddWebHook")
		return err
	}

	return nil
}