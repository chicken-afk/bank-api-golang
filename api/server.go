package api

import (
	db "bank-api/db/sqlc"

	"github.com/gin-gonic/gin"
)

//Server serves HTTP request for service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

//New Server HTTP Request
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	//Account API
	router.POST("/account", server.CreateAccount)
	router.PUT("/account", server.EditAccount)
	router.GET("/account/:id", server.GetAccount)
	router.GET("/accounts", server.ListAccount)
	router.DELETE("/account/:id", server.DeleteAccount)
	router.POST("/account-balance/:id", server.EditBalance)

	server.router = router
	return server
}

//Start HTTP server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
