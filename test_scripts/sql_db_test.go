package test_scripts

import (
	"github.com/parvez0/metric-ingestion/objects"
	sqlitedb "github.com/parvez0/metric-ingestion/sqlite-db"
	"math/rand"
	"testing"
)

var db *sqlitedb.SQLDB
var testTable = "test_metrics"

// TestDBConnection establishes a connection with sqlite
func TestDBConnection(t *testing.T)  {
	db = sqlitedb.CreateDbConnection()
}

// TestPopulateDB creates initial tables and populates the db
func TestPopulateDB(t *testing.T) {
	err := db.PopulateDB(testTable)
	if err != nil{
		t.Errorf("failed to populate data - %+v", err)
		return
	}
	t.Logf("db populated successfully.")
}

// TestInsertRecords
func TestInsertRecords(t *testing.T)  {
	data := objects.Metrics{
		CpuUsed:    rand.Intn(100),
		MemoryUsed: rand.Intn(100),
		Ip: "10.0.0.1",
	}
	res, err := db.Insert(testTable, &data)
	if err != nil{
		t.Fatalf("failed to insert record in table - %s -- error -- %+v", testTable, err)
	}
	t.Logf("data inserted successfully, last inserted record id - %v", res)
}

// TestFetchData selects all records and prints them in loop
func TestFetchData(t *testing.T)  {
	rows, err := db.Select(testTable, "")
	if err != nil{
		t.Fatalf("failed to fetch data from db - %+v", err)
	}
	for rows.Next(){
		row := objects.Metrics{}
		err := rows.Scan(&row.CpuUsed, &row.MemoryUsed, &row.Ip, &row.Date)
		if err != nil{
			t.Fatalf("failed to parse fetched data - %+v", err)
		}
		t.Logf("data fetched successfully - %+v", row)
	}
}

// TestDropTable drops the test table after testing is done
func TestDropTable(t *testing.T) {
	err := db.DropTable(testTable)
	if err != nil{
		t.Errorf("failed to drop table - %+v", err)
		return
	}
	t.Logf("table dropped successfully.")
}