package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">>>>>> Before ", account1.Balance, account2.Balance)
	n := 5
	amount := int64(10)

	// Handle concurrent with routines and channels

	errors := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.trasnferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errors <- err
			results <- result
		}()

	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fEntry := result.FromEntry
		tEntry := result.ToEntry
		require.NotEmpty(t, fEntry)
		require.Equal(t, fEntry.AccountID, account1.ID)
		require.Equal(t, fEntry.Amount, -amount)

		require.NotEmpty(t, tEntry)
		require.Equal(t, tEntry.AccountID, account2.ID)
		require.Equal(t, tEntry.Amount, amount)

		_, err = store.GetEntry(context.Background(), fEntry.ID)
		require.NoError(t, err)

		_, err = store.GetEntry(context.Background(), tEntry.ID)
		require.NoError(t, err)

		fAccount := result.FromAccount
		require.NotEmpty(t, fAccount)
		require.Equal(t, fAccount.ID, account1.ID)

		tAccount := result.ToAccount
		require.NotEmpty(t, tAccount)
		require.Equal(t, tAccount.ID, account2.ID)
		diff1 := account1.Balance - fAccount.Balance
		diff2 := tAccount.Balance - account2.Balance
		fmt.Println(">>>>>> tx ", fAccount.Balance, tAccount.Balance)
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff2 > 0)
		require.True(t, diff1%amount == 0)
		k := int(diff1 / amount) // num of transactions
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	// check the final balance of the accounts
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
	fmt.Println(">>>>>> After ", updateAccount1.Balance, updateAccount2.Balance)
}



// this is a deadlock that might occur because of transferring money from account 1 to 2 and transferring money from 2 to 1
// the solution is to make sure the application acquires lock in a consistent order, so we may have a lock at first but evententually it will be released
func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">>>>>> Before ", account1.Balance, account2.Balance)
	n := 10
	amount := int64(10)

	// Handle concurrent with routines and channels

	errors := make(chan error)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		fromAccountID := account1.ID
			ToAccountID := account2.ID

			if i%2 == 1 {
				fromAccountID = account2.ID
				ToAccountID = account1.ID
			}
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			
			_, err := store.trasnferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   ToAccountID,
				Amount:        amount,
			})
			errors <- err
		}()

	}

	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)
		
	}

	// check the final balance of the accounts
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
	fmt.Println(">>>>>> After ", updateAccount1.Balance, updateAccount2.Balance)
}
