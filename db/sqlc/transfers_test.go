package db

import (
	"context"
	"testing"
	"time"

	"github.com/diegocrequi/backend-master-class-golang-postgres-kubernetes-grpc/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, accountIDs ...int64) Transfer {
	var from_account, to_account Account
	if len(accountIDs) == 0 {
		from_account = createRandomAccount(t)
		to_account = createRandomAccount(t)
	} else {
		from_account1, err := testQueries.GetAccount(context.Background(), accountIDs[0])
		require.NoError(t, err)
		require.NotEmpty(t, from_account1)
		from_account = from_account1

		to_account1, err := testQueries.GetAccount(context.Background(), accountIDs[1])
		require.NoError(t, err)
		require.NotEmpty(t, to_account1)
		to_account = to_account1
	}

	arg := CreateTransferParams{
		FromAccountID: from_account.ID,
		ToAccountID:   to_account.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.CreatedAt)
	require.NotZero(t, transfer.ID)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	transfer1, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer1)

	require.Equal(t, transfer.ID, transfer1.ID)
	require.Equal(t, transfer.FromAccountID, transfer1.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transfer1.ToAccountID)
	require.Equal(t, transfer.Amount, transfer1.Amount)
	require.WithinDuration(t, transfer.CreatedAt, transfer1.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	from_account := createRandomAccount(t)
	to_account := createRandomAccount(t)
	for range 10 {
		createRandomTransfer(t, from_account.ID, to_account.ID)
	}

	arg := ListTransfersParams{
		FromAccountID: from_account.ID,
		ToAccountID:   to_account.ID,
		Limit:         5,
		Offset:        5,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
