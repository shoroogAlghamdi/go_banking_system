package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/banking_system/db/util"
)

func CreateRandomTransfer(t *testing.T, from_account_id int64, to_account_id int64) Transfer {

	args := CreateTransferParams{
		FromAccountID: from_account_id,
		ToAccountID: to_account_id,
		Amount:    util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.Equal(t, args.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer

}

func TestCreateTransfer(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	CreateRandomTransfer(t, account1.ID, account2.ID)
}

func TestGetTransfer(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	transfer1 := CreateRandomTransfer(t, account1.ID, account2.ID) 
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}



func TestListTransfers(t *testing.T) {
	accounts := []Account{ CreateRandomAccount(t), CreateRandomAccount(t)}

	for i:=0; i< 10; i++ {
		CreateRandomTransfer(t, accounts[0].ID, accounts[1].ID)
	}

	
	args := ListTransfersParams {
		FromAccountID: accounts[0].ID,
		ToAccountID: accounts[1].ID,
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 5)

	for _, transfer :=range transfers {
		require.NotEmpty(t, transfer)
	}
}
