package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/komron-dev/bank/util"
	"github.com/stretchr/testify/require"
)

func GetAccountHelper(t *testing.T) Account {
	acc, err := testQueries.GetRandomAccount(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, acc)

	return acc
}

func TestCreateAccount(t *testing.T)  {
	arg := CreateAccountParams {
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	
	require.NoError(t, err)
	require.NotEmpty(t, account)
	
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T)  {
	acc1 := GetAccountHelper(t)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)

	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T)  {
	account1 := GetAccountHelper(t)

	arg := UpdateAccountParams{
		ID: account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T)  {
	account1 := GetAccountHelper(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T)  {
	count, err := testQueries.GetAccountsCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count)

	var randomNum = int32(util.RandomInt(1, count))
	arg := ListAccountsParams{
		Limit: randomNum,
		Offset: randomNum,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	// require.Len(t, accounts, int(randomNum))

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}