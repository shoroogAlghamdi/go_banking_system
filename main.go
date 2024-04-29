package main
import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"github.com/shoroogAlghamdi/banking_system/api"
	"github.com/shoroogAlghamdi/banking_system/util"
	db "github.com/shoroogAlghamdi/banking_system/db/sqlc"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Couldn't read config vars!")
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Couldn't connect to DB!")
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}