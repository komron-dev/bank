package db

import (
	"context"
	"fmt"
 	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t) 

	fmt.Printf("--- Sender: %v ---\n", account1.Owner)
	fmt.Printf("--- Reciepent: %v ---\n", account2.Owner)

	fmt.Printf(">>> BEFORE: account1 balance: %v, account2 balance: %v\n", account1.Balance, account2.Balance)
	// run concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				SenderID: account1.ID,
				ReciepentID: account2.ID,
				Amount: amount,
			}) 

			errs <- err
			results <- result 
		}()
	}

	// check results
	existed := make(map[int]bool)

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

		// accounts

		sender := result.Sender
		require.NotEmpty(t, sender)
		require.Equal(t, account1.ID, sender.ID)

		reciepent := result.Reciepent
		require.NotEmpty(t, reciepent)
		require.Equal(t, account2.ID, reciepent.ID)

		fmt.Printf("----- tx: %v,  %v\n", sender.Balance, reciepent.Balance)
		// accounts' balace
		diff1 := account1.Balance - sender.Balance
		diff2 := reciepent.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % amount == 0)

		k := int(diff1/amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	
	require.Equal(t, account1.Balance - int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance + int64(n)*amount, updatedAccount2.Balance)

	fmt.Printf(">>> AFTER: account1 balance: %v, account2 balance: %v\n", updatedAccount1.Balance, updatedAccount2.Balance)

}