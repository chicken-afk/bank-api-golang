package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("Before TX : ", account1.Balance, account2.Balance)

	//run n councurrent transfer transactions
	n := 10
	amount := int64(50)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {

		FromAccountId := account1.ID
		ToAccountId := account2.ID

		if i%2 == 1 {
			FromAccountId = account2.ID
			ToAccountId = account1.ID
		}

		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: FromAccountId,
				ToAccountId:   ToAccountId,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//Check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		// require.Equal(t, account1.ID, transfer.FromAccountID)
		// require.Equal(t, account2.ID, transfer.ToAccountID)
		// require.Equal(t, amount, transfer.Amount)
		// require.NotZero(t, transfer.ID)
		// require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		//Check Transer Tabel
		FromEntry := result.FromEntry
		require.NotEmpty(t, FromEntry)
		// require.Equal(t, account1.ID, FromEntry.AccountID)
		// require.Equal(t, -amount, FromEntry.Amount)
		// require.NotZero(t, FromEntry.ID)
		// require.NotZero(t, FromEntry.CreatedAt)

		/** Check Entries To Account*/
		_, err = store.GetEntry(context.Background(), GetEntryParams{
			Limit:     1,
			AccountID: FromEntry.AccountID,
		})
		require.NoError(t, err)
		toEntry := result.ToEntry
		// require.NotEmpty(t, toEntry)
		// require.Equal(t, account2.ID, toEntry.AccountID)
		// require.Equal(t, amount, toEntry.Amount)
		// require.NotZero(t, toEntry.ID)
		// require.NotZero(t, toEntry.CreatedAt)

		// Get Account
		_, err = store.GetEntry(context.Background(), GetEntryParams{
			Limit:     1,
			AccountID: toEntry.ID,
		})
		require.NoError(t, err)

		FromAccount := result.FromAccount
		ToAccount := result.ToAccount

		fmt.Println("After TX : ", FromAccount.Balance, ToAccount.Balance)

		//Check update balance
		// diff1 := account1.Balance - FromAccount.Balance
		// diff2 := ToAccount.Balance - account2.Balance
		// require.Equal(t, diff1, diff2)
		// require.True(t, diff1 > 0)
		// require.True(t, diff1%amount == 0)

	}

}
