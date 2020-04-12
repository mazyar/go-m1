package m1

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
	log "github.com/sirupsen/logrus"
)

type DataBase struct {
	DB      *sql.DB
	Context *context.Context
}

// Singletone instance
var (
	_instance *DataBase
	_error    error
	once      sync.Once
)

func ProvideDataBase(connectionString string, ctx *context.Context) (*DataBase, error) {

	once.Do(func() {

		var db *sql.DB
		var err error
		db, err = sql.Open("sqlserver", connectionString)
		// The server can not be connected to the database server
		if err != nil {
			log.Println(fmt.Sprintf("Error creating connection pool: %s", err.Error()))
			// returns error
			_instance = nil
			_error = err
		} else {
			_instance = &DataBase{DB: db, Context: ctx}
			_error = nil
		}
	})

	return _instance, _error
}
