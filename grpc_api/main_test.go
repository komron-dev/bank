package grpc_api

import (
	db "github.com/komron-dev/bank/db/sqlc"
	"github.com/komron-dev/bank/util"
	"github.com/komron-dev/bank/worker"
	"github.com/stretchr/testify/require"
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
