package test_scripts

import (
	sqlitedb "github.com/parvez0/whatsapp-provisioner/sqlite-db"
	"testing"
)

var db *sqlitedb.SQLDB

// TestDBConnection establishes a connection with sqlite
func TestDBConnection(t *testing.T)  {
	db = sqlitedb.CreateDbConnection()
}

// TestPopulateDB creates initial tables and populates the db
func TestPopulateDB(t *testing.T) {
	testTable := "test_infra_queue"
	err := db.PopulateDB(testTable)
	if err != nil{
		t.Errorf("failed to populate data - %+v", err)
		return
	}
	err = db.DropTable(testTable)
	if err != nil{
		t.Errorf("failed to drop table - %+v", err)
		return
	}
	t.Logf("db populated successfully.")
}