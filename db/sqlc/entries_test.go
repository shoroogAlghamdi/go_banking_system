package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/shoroogAlghamdi/banking_system/util"
)

func CreateRandomEntry(t *testing.T, account_id int64) Entry {

	args := CreateEntryParams{
		AccountID: account_id,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry

}
func TestCreateEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	CreateRandomEntry(t, account.ID)
}

func TestGetEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	entry1 := CreateRandomEntry(t, account.ID)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}



func TestListEntries(t *testing.T) {
	account := CreateRandomAccount(t)
	for i:=0; i< 10; i++ {
		CreateRandomEntry(t, account.ID)
	}

	
	args := ListEntriesParams {
		AccountID: account.ID,
		Limit: 5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, 5)

	for _, entry :=range entries {
		require.NotEmpty(t, entry)
	}
}
