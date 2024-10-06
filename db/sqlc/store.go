package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	SenderID    int64 `json:"sender_id"`
	RecipientID int64 `json:"recipient_id"`
	Amount      int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer  Transfer `json:"transfer"`
	Sender    Account  `json:"sender"`
	Recipient Account  `json:"recipient"`
	FromEntry Entry    `json:"from_entry"`
	ToEntry   Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	account, err := store.GetAccount(ctx, arg.SenderID)
	if account.Balance < arg.Amount {
		return TransferTxResult{}, fmt.Errorf("not enough funds in sender account (ID: %d), current balance: %d, required: %d", arg.SenderID, account.Balance, arg.Amount)
	}

	var result TransferTxResult
	err = store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			SenderID:    arg.SenderID,
			RecipientID: arg.RecipientID,
			Amount:      arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.SenderID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.RecipientID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// update account's balance

		if arg.SenderID < arg.RecipientID {
			result.Sender, result.Recipient, err = addMoney(ctx, q, arg.SenderID, -arg.Amount, arg.RecipientID, arg.Amount)
		} else {
			result.Recipient, result.Sender, err = addMoney(ctx, q, arg.RecipientID, arg.Amount, arg.SenderID, -arg.Amount)
		}

		return err
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})
	if err != nil {
		return
	}

	return
}
