package sqlite_db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/parvez0/metric-ingestion/custom_logger"
	"os"
)

// custom wrapper for sqlite db object
type SQLDB struct {
	Db *sql.DB
}

// initialize global logger
var clog = custom_logger.NewLogger()

// CreateDbConnection initializes the sqlite and
func CreateDbConnection() *SQLDB {
	dbpath := "/data"
	if path := os.Getenv("SQLITE_DB_PATH"); path != ""{
		dbpath = path
	}
	//
	db, err := sql.Open("sqlite3", dbpath + "/wa.sqlite")
	if err != nil{
		clog.Panicf("failed to initialize sqlite - %+v", err)
	}
	return &SQLDB{
		Db: db,
	}
}

// PopulateDB initializes the db with table structure and initial values
func (db *SQLDB) PopulateDB(table string) error {
	if table == ""{
		table = "infra_queue"
	}
	query := fmt.Sprintf(`create table %s(
					wa_id varchar(255), 
					number varchar(255),
					status varchar(255),
					ip varchar(255),
					tier int
					 )`, table)
	_, err := db.Db.Exec(query)
	return err
}

// DropTable drops the table after testing
func (db *SQLDB) DropTable(table string) error {
	if table == ""{
		return errors.New("table name not provided")
	}
	clog.Warn("dropping table - ", table)
	_, err := db.Db.Exec("drop table " + table)
	return err
}

