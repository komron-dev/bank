package grpc_api

import (
	"fmt"
	db "github.com/komron-dev/bank/db/sqlc"
	"github.com/komron-dev/bank/pb"
	"github.com/komron-dev/bank/token"
	"github.com/komron-dev/bank/util"
)

type Server struct {
	pb.UnimplementedBankServer
	store      db.Store
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

	return server, nil
}