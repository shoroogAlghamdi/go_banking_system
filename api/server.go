package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/shoroogAlghamdi/banking_system/db/sqlc"
	"github.com/shoroogAlghamdi/banking_system/token"
	"github.com/shoroogAlghamdi/banking_system/util"
)
type Server struct {
	config util.Config
	store db.Store
	tokenMaker token.Maker
	router *gin.Engine
}


func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPaseto(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}
	server := Server{config: config, store: store, tokenMaker: tokenMaker}
	server.setupRouters()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}
	
	return &server, nil
}

func (server *Server) setupRouters() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts/", server.listAccounts)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)
	authRoutes.PUT("/accounts/", server.updateAccount)
	authRoutes.POST("/transfers", server.transferMoney)
	server.router = router 
}
func errResponse(err error) gin.H {
	// gin.H is a map of map[string]interface{} 
	return gin.H{"error": err.Error()}
}

// public function cuz router is private
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}