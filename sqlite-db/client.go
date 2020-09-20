package sqlite_db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/parvez0/metric-ingestion/custom_logger"
	"github.com/parvez0/metric-ingestion/objects"
	"os"
)

// custom wrapper for sqlite db object
type SQLDB struct {
	Db *sql.DB
	Table string
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
		table = "metrics"
	} else {
		db.Table = table
	}
	query := fmt.Sprintf(`create table %s(
					percentage_cpu_used varchar(255), 
					percentage_memory_used varchar(255),
					ip varchar(255),
					date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
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

// Insert create a metric record in sql database
func (db *SQLDB) Insert(table string, data *objects.Metrics) (sql.Result, error) {
	if table == ""{
		clog.Debugf("table name not provided for insert query, using default %s", db.Table)
		table = db.Table
	}
	clog.Debugf("inserting data into table %s --- %+v", table, data)
	stm, err := db.Db.Prepare(fmt.Sprintf(`INSERT INTO %s(percentage_cpu_used, percentage_memory_used, ip) values(?, ?, ?)`, table))
	if err != nil{
		return nil, err
	}
	return stm.Exec(data.CpuUsed, data.MemoryUsed, data.Ip)
}

// Select returns rows from the sqlite
func (db *SQLDB) Select(table string, query string) (*sql.Rows, error) {
	if table == ""{
		clog.Debugf("table name not provided for select query, using default %s", db.Table)
		table = db.Table
	}
	if query == ""{
		query = fmt.Sprintf("select * from %s limit 10", table)
	}
	return db.Db.Query(query)
}