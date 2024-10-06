package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/komron-dev/bank/db/sqlc"
)

type createTransferRequest struct {
	SenderID    int64  `json:"sender_id" binding:"required"`
	RecipientID int64  `json:"recipient_id" binding:"required"`
	Amount      int64  `json:"amount" binding:"required"`
	Currency    string `json:"currency" binding:"required,currency-validate,gt=0"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var request createTransferRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.isValidAccount(ctx, request.RecipientID, request.Currency) {
		return
	}
	if !server.isValidAccount(ctx, request.SenderID, request.Currency) {
		return
	}

	arg := db.TransferTxParams{
		RecipientID: request.RecipientID,
		SenderID:    request.SenderID,
		Amount:      request.Amount,
	}

	transfer, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

func (server *Server) isValidAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account %d currency mismatch: %s and %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
