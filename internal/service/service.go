package service

import (
	"github.com/rs/zerolog/log"

	"github.com/go-worker-golden-data-rmq/internal/repository/db_postgre"

)

var childLogger = log.With().Str("service", "DataService").Logger()

type WorkerService struct {
	workerRepository 	db_postgre.WorkerRepository
	//cardNotification notification.CardNotification
}

func NewWorkerService(workerRepository db_postgre.WorkerRepository,
					//cardNotification notification.CardNotification,
					) *WorkerService{
	childLogger.Debug().Msg("NewWorkerService")
	return &WorkerService{
		workerRepository: workerRepository,
	//	cardNotification: cardNotification,
	}
}