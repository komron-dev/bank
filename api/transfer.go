package api

import (
	"errors"
	"fmt"
	"github.com/komron-dev/bank/token"
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

	sender, isValid := server.isValidAccount(ctx, request.SenderID, request.Currency)
	if !isValid {
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if sender.Owner != authPayload.Username {
		err := errors.New("sender account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, isValid = server.isValidAccount(ctx, request.RecipientID, request.Currency)
	if !isValid {
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

func (server *Server) isValidAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
