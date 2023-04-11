package db_postgre

import (
	"context"

	_ "github.com/lib/pq"

	"github.com/go-worker-golden-data-rmq/internal/core"
	"github.com/go-worker-golden-data-rmq/internal/erro"

)

type WorkerRepository struct {
	databaseHelper DatabaseHelper
}

func NewWorkerRepository(databaseHelper DatabaseHelper) WorkerRepository {
	childLogger.Debug().Msg("NewWorkerRepository")
	return WorkerRepository{
		databaseHelper: databaseHelper,
	}
}

//---------------------------

func (w WorkerRepository) Ping() (bool, error) {
	childLogger.Debug().Msg("++++++++++++++++++++++++++++++++")
	childLogger.Debug().Msg("Ping")
	childLogger.Debug().Msg("++++++++++++++++++++++++++++++++")

	ctx, cancel := context.WithTimeout(context.Background(), 1000)
	defer cancel()

	client, _ := w.databaseHelper.GetConnection(ctx)
	err := client.Ping()
	if err != nil {
		return false, erro.ErrConnectionDatabase
	}

	return true, nil
}

func (w WorkerRepository) GetPerson(id string) (core.Person, error){
	childLogger.Debug().Msg("++++++++++++++++++++++++++++++++")
	childLogger.Debug().Msg("GetPerson")
	childLogger.Debug().Msg("++++++++++++++++++++++++++++++++")

	ctx, cancel := context.WithTimeout(context.Background(), 1000)
	defer cancel()

	client, _ := w.databaseHelper.GetConnection(ctx)
	person := core.Person{}

	rows, err := client.Query(`SELECT id, email FROM person WHERE id = $1`, id)
	defer rows.Close()
	if err != nil {
		return person, erro.ErrConnectionDatabase
	}
	for rows.Next() {
		err := rows.Scan( &person.ID, &person.Email )
		if err != nil {
			return person , erro.ErrNotFound
        }
		return person, nil
	}

	return person , erro.ErrNotFound
}