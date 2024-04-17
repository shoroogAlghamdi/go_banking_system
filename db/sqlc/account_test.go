package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/banking_system/db/util"
)

func CreateRandomAccount(t *testing.T) Account {
	// Could generate random data 
	// args := CreateAccountParams{
	// 	Owner:    "Tom",
	// 	Balance:  100,
	// 	Currency: "USD",
	// }

	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account

}
func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetaccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}


func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	updated_account := UpdateAccountParams{
		ID: account1.ID,
		Balance: util.RandomMoney(),
	}
	_, err := testQueries.UpdateAccount(context.Background(), updated_account)
	require.NoError(t, err)
	account2, _ := testQueries.GetAccount(context.Background(), account1.ID)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.NotEqual(t, account1.Balance, account2.Balance)
	require.Equal(t, updated_account.Balance, account2.Balance)
	
}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	_, err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

}

func TestListAccounts(t *testing.T) {
	for i:=0; i< 10; i++ {
		CreateRandomAccount(t)
	}

	args := ListAccountsParams {
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Equal(t, int(len(accounts)), int(args.Limit))
	require.Len(t, accounts, 5)

	for _, account :=range accounts {
		require.NotEmpty(t, account)
	}
}