package api 

import (
	"github.com/gin-gonic/gin"
	db "github.com/shoroogAlghamdi/go_banking_system/db/sqlc"
)
type Server struct {
	store *db.Store
	router *gin.Engine
}


func NewServer(store *db.Store) Server {
	server := Server{store: store}
	router := gin.Default()
	router.POST("/accounts", server.CreateAccount)
	server.router = router 
	return server 
}

func errResponse(err error) gin.H {
	// gin.H is a map of map[string]interface{} 
	return gin.H{"error": err.Error()}
}

// public function cuz router is private
func (server *Server) Start(address string) err error {
	return server.router.run(address)
}