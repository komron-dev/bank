package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// fmt.Printf("--- Sender: %v ---\n", account1.Owner)
	// fmt.Printf("--- Recipient: %v ---\n", account2.Owner)

	// fmt.Printf(">>> BEFORE: account1 balance: %v, account2 balance: %v\n", account1.Balance, account2.Balance)

	// run concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				SenderID:    account1.ID,
				RecipientID: account2.ID,
				Amount:      amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer result
		transfer := result.Transfer

		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.SenderID)
		require.Equal(t, account2.ID, transfer.RecipientID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, account1.ID)
		require.NotZero(t, account1.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, account2.ID)
		require.NotZero(t, account2.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// accounts

		sender := result.Sender
		require.NotEmpty(t, sender)
		require.Equal(t, account1.ID, sender.ID)

		recipient := result.Recipient
		require.NotEmpty(t, recipient)
		require.Equal(t, account2.ID, recipient.ID)

		// fmt.Printf("----- tx: %v,  %v\n", sender.Balance, recipient.Balance)
		// accounts' balace
		diff1 := account1.Balance - sender.Balance
		diff2 := recipient.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

	// fmt.Printf(">>> AFTER: account1 balance: %v, account2 balance: %v\n", updatedAccount1.Balance, updatedAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// fmt.Printf("--- Sender: %v ---\n", account1.Owner)
	// fmt.Printf("--- Recipient: %v ---\n", account2.Owner)

	// fmt.Printf(">>> BEFORE: account1 balance: %v, account2 balance: %v\n", account1.Balance, account2.Balance)
	// run concurrent transfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		senderID := account1.ID
		recipientID := account2.ID

		if i%2 == 1 {
			senderID = account2.ID
			recipientID = account1.ID
		}

		go func() {
			_, err := testStore.TransferTx(context.Background(), TransferTxParams{
				SenderID:    senderID,
				RecipientID: recipientID,
				Amount:      amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	// fmt.Printf(">>> AFTER: account1 balance: %v, account2 balance: %v\n", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
