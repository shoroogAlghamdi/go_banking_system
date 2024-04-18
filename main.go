package main
import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"github.com/shoroogAlghamdi/go_banking_system/api"
	db "github.com/shoroogAlghamdi/go_banking_system/db/sqlc"
)
const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:shoroog@localhost:5432/banking_system?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Couldn't connect to DB!")
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}