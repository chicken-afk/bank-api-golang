package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	arg := CreateEntryParams{
		AccountID: 12,
		Amount:    100000,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.EqualValues(t, 12, entry.AccountID)
}

func TestGetEntry(t *testing.T) {
	entry, err := testQueries.GetEntry(context.Background(), 3)
	require.NotEmpty(t, entry)
	require.NoError(t, err)
	require.EqualValues(t, entry.ID, 3)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 5; i++ {
		argCreate := CreateEntryParams{
			AccountID: 22,
			Amount:    100000,
		}
		entry, err := testQueries.CreateEntry(context.Background(), argCreate)
		require.NoError(t, err)
		require.NotEmpty(t, entry)
	}
	params := ListEntriesParams{
		AccountID: 22,
		Limit:     5,
		Offset:    0,
	}
	entries, err := testQueries.ListEntries(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, 5)

}
