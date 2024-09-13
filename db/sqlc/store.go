package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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
	SenderID int64 `json:"sender_id"`
	ReciepentID int64 `json:"reciepent_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer  Transfer `json:"transfer"`
	Sender  Account `json:"sender"`
	Reciepent  Account `json:"reciepent"`
	FromEntry  Entry `json:"from_entry"`
	ToEntry  Entry `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			SenderID: arg.SenderID,
			ReciepentID: arg.ReciepentID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.SenderID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ReciepentID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		
		// update account's balance

		if arg.SenderID < arg.ReciepentID {
			result.Sender, result.Reciepent, err =  addMoney(ctx, q, arg.SenderID, -arg.Amount, arg.ReciepentID, arg.Amount)
		} else {
			result.Reciepent, result.Sender, err =  addMoney(ctx, q, arg.ReciepentID, arg.Amount, arg.SenderID, -arg.Amount)
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
	account1, err = q.AddAcountBalance(ctx, AddAcountBalanceParams{
		Amount: amount1,
		ID: accountID1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAcountBalance(ctx, AddAcountBalanceParams{
		Amount: amount2,
		ID: accountID2,
	})
	if err != nil {
		return
	}

	return
}