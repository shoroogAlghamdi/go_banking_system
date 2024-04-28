package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"github.com/shoroogAlghamdi/banking_system/util"

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
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Couldn't read config vars!")
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Couldn't connect to DB!")
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
