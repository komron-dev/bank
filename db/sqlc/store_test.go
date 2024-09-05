package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := GetAccountHelper(t)
	account2 := GetAccountHelper(t) 

	fmt.Printf("--- Sender: %v ---\n", account1.Owner)
	fmt.Printf("--- Reciepent: %v ---\n", account2.Owner)

	// run concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), CreateTransferParams{
				SenderID: account1.ID,
				ReciepentID: account2.ID,
				Amount: amount,
			}) 

			errs <- err
			results <- result 
		}()
	}

	// check results

	for i := 0; i < n; i++ {
		err := <- errs 
		require.NoError(t, err)

		result := <- results
		require.NotEmpty(t, result)

		// check transfer result
		transfer := result.Transfer

		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.SenderID)
		require.Equal(t, account2.ID, transfer.ReciepentID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, account1.ID)
		require.NotZero(t, account1.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
	
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, account2.ID)
		require.NotZero(t, account2.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
	}

}