package grpc_api

import (
	"context"
	"fmt"
	db "github.com/komron-dev/bank/db/sqlc"
	"github.com/komron-dev/bank/token"
	"github.com/komron-dev/bank/util"
	"github.com/komron-dev/bank/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)
	require.NotEmpty(t, server)

	return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}
	return metadata.NewIncomingContext(context.Background(), md)
}
