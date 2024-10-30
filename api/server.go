package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/komron-dev/bank/db/sqlc"
	"github.com/komron-dev/bank/token"
	"github.com/komron-dev/bank/util"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency-validate", validCurrency)
		if err != nil {
			return nil, err
		}
	}
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authorizedRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authorizedRoutes.POST("/accounts", server.createAccount)
	authorizedRoutes.GET("/accounts/:id", server.getAccount)
	authorizedRoutes.GET("/accounts", server.listAccounts)
	authorizedRoutes.DELETE("/accounts/:id", server.deleteAccount)
	authorizedRoutes.PUT("/accounts", server.updateAccount)

	authorizedRoutes.POST("/transfers", server.createTransfer)
	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
