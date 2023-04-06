package db_postgre

type WorkerRepository struct {
	DatabaseHelper DatabaseHelper
}

func NewWorkerRepository(databaseHelper DatabaseHelper) WorkerRepository {
	return WorkerRepository{
		DatabaseHelper: databaseHelper,
	}
}
