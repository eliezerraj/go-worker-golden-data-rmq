package db_postgre

import (
	"context"
	"time"
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/go-worker-golden-data-rmq/internal/core"

)

var childLogger = log.With().Str("repository/db_postgre", "NewDatabaseHelper").Logger()

type DatabaseHelper interface {
	GetConnection(ctx context.Context) (*sql.DB, error)
	CloseConnection()
}

type DatabaseHelperImpl struct {
	client   	*sql.DB
	timeout		time.Duration
}

func (d DatabaseHelperImpl) GetConnection(ctx context.Context) (*sql.DB, error) {
	childLogger.Debug().Msg("GetConnection")
	return d.client, nil
}

func (d DatabaseHelperImpl) CloseConnection()  {
	childLogger.Debug().Msg("CloseConnection")
	defer d.client.Close()
}

func NewDatabaseHelper(databaseRDS core.DatabaseRDS) (DatabaseHelper, error) {
	childLogger.Debug().Msg("NewDatabaseHelper")

	_ , cancel := context.WithTimeout(context.Background(), time.Duration(databaseRDS.Db_timeout)*time.Second)
	defer cancel()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", 
							databaseRDS.User, 
							databaseRDS.Password, 
							databaseRDS.Host, 
							databaseRDS.Port, 
							databaseRDS.DatabaseName) 
	
	//fmt.Println("==========>", databaseRDS.Postgres_Driver, connStr)

	client, err := sql.Open(databaseRDS.Postgres_Driver, connStr)
	if err != nil {
		return DatabaseHelperImpl{}, err
	}
	err = client.Ping()
	if err != nil {
		return DatabaseHelperImpl{}, err
	}

	return DatabaseHelperImpl{
		client: client,
		timeout:  time.Duration(databaseRDS.Db_timeout) * time.Second,
	}, nil
}