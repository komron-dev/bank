package db

import (
	"context"
	"testing"
	"time"

	"github.com/komron-dev/bank/util"
	"github.com/stretchr/testify/require"
)

func GetEntryHelper(t *testing.T) Entry {
	entry, err := testQueries.GetRandomEntry(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	return entry
}

func TestCreateEntry(t *testing.T)  {
	account := GetAccountHelper(t)

	arg := CreateEntryParams {
		AccountID: account.ID,
		Amount: util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func TestGetEntry(t *testing.T)  {
	entry1 := GetEntryHelper(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)

	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T)  {
	count, err := testQueries.GetEntriesCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count)

	arg := ListEntriesParams{
		Limit: int32(util.RandomInt(1, count)),
		Offset: int32(util.RandomInt(1, count)),
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	// require.Len(t, entries, int(count))

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}