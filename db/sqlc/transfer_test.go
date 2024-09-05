package db

import (
	"context"
	"testing"
	"time"

	"github.com/komron-dev/bank/util"
	"github.com/stretchr/testify/require"
)

func GetTransferHelper(t *testing.T) Transfer {
	transfer, err := testQueries.GetRandomTransfer(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	return transfer
}

func TestCreateTransfer(t *testing.T)  {
	account1 := GetAccountHelper(t)
	account2 := GetAccountHelper(t)

	require.NotEqual(t, account1, account2)

	arg := CreateTransferParams {
		SenderID: account1.ID,
		ReciepentID: account2.ID,
		Amount: util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	
	require.Equal(t, arg.SenderID, transfer.SenderID)
	require.Equal(t, arg.ReciepentID, transfer.ReciepentID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}

func TestGetTransfer(t *testing.T)  {
	entry1 := GetEntryHelper(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)

	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T)  {
	count, err := testQueries.GetTransfersCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count)

	arg := ListTransfersParams{
		Limit: int32(util.RandomInt(1, count)),
		Offset: int32(util.RandomInt(1, count)),
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	// require.Len(t, transfers, int(count))

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}