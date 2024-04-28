package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/shoroogAlghamdi/banking_system/db/sqlc"
	"github.com/shoroogAlghamdi/banking_system/token"
)

type transferMoneyRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency`
}

func (server *Server) transferMoney(ctx *gin.Context) {
	var req transferMoneyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	fromAccount, valid := server.validateAccount(ctx, req.FromAccountID, req.Currency) 
	if !valid {
		return
	}

	if authPayload.Username != fromAccount.Owner {
		err := errors.New("cannot send money from account that does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return 
	}

	_, valid = server.validateAccount(ctx, req.ToAccountID, req.Currency) 
	if !valid {
		return
	}

	
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	res, err := server.store.TrasnferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res)

}

func (server *Server) validateAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err = fmt.Errorf("currency mismatch for account %d , got %s, but actually account currency is %s", accountId, currency, account.Currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return account, false
	}

	return account, true

}
