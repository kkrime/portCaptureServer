package app

import (
	"fmt"
	"portCaptureServer/app/config"
	"portCaptureServer/app/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/logrusadapter"
)

func (a *app) ConnectToDB(config config.DBConfig) (*sqlx.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)

	database, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	dbLog := logger.CreateNewLogger()

	database.DB = sqldblogger.OpenDriver(dsn, database.DB.Driver(), logrusadapter.New(dbLog),
		sqldblogger.WithTimeFormat(sqldblogger.TimeFormatRFC3339),
		sqldblogger.WithLogDriverErrorSkip(true),
		sqldblogger.WithSQLQueryAsMessage(true))

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	return database, nil
}
