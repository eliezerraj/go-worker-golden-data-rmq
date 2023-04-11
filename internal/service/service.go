package service

import (
	"github.com/rs/zerolog/log"

	"github.com/go-worker-golden-data-rmq/internal/repository/db_postgre"

)

var childLogger = log.With().Str("service", "DataService").Logger()

type WorkerService struct {
	workerRepository 	*db_postgre.WorkerRepository
}

func NewWorkerService(workerRepository *db_postgre.WorkerRepository ) *WorkerService{
	childLogger.Debug().Msg("NewWorkerService")
	return &WorkerService{
		workerRepository: workerRepository,
	}
}

func (s *WorkerService) DataEnrichment() error{
	childLogger.Debug().Msg("--------------------------")
	childLogger.Debug().Msg("DataEnrichment")
	childLogger.Debug().Msg("---------------------------")

	id := "1"
	result, err := s.workerRepository.GetPerson(id)
	if err != nil {
		log.Error().Err(err).Msg("error workerRepository.GetPerson")
	}

	log.Debug().Interface("result",result).Msg("")

	return nil
}