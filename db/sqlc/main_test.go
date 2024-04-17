package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	// this package is not used directly but we need it for the driver, if we save the formatter will remove it because it is not used, so we need to add - to keep it
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:shoroog@localhost:5432/banking_system?sslmode=disable"
)

// global var
var testQueries *Queries
var testDB *sql.DB
var err error 
// main entry point for all tests
func TestMain(m *testing.M) {
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Couldn't connect to DB!")
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
