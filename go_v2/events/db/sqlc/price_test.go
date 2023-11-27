package db

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"

	"github.com/amirul-zafrin/event/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomPrice(t *testing.T) Price {
	event := CreateRandomEvent(t)
	arg := CreatePriceParams{
		Class: util.RandomClass(),
		Price: 100 * rand.Float64(),
		Event: int32(event.ID),
	}
	price, err := testQueries.CreatePrice(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, price)

	return price
}

func TestCreatePrice(t *testing.T) {
	CreateRandomPrice(t)
}

func TestGetPrice(t *testing.T) {
	price1 := CreateRandomPrice(t)

	price2, err := testQueries.GetPrice(context.Background(), price1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, price2)

	require.Equal(t, price1.Class, price2.Class)
	require.Equal(t, price1.Price, price2.Price)
	require.Equal(t, price1.Event, price2.Event)
}

func TestListAllPrice(t *testing.T) {
	for i := 0; i < 5; i++ {
		CreateRandomPrice(t)
	}
	arg := ListAllPricesParams{
		Limit:  5,
		Offset: 1,
	}
	prices, err := testQueries.ListAllPrices(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, prices)

	for _, price := range prices {
		require.NotEmpty(t, price)
	}
}

func TestUpdatePrice(t *testing.T) {
	price := CreateRandomPrice(t)

	arg := UpdatePriceParams{
		ID:    price.ID,
		Price: 100 * rand.Float64(),
	}

	price1, err := testQueries.UpdatePrice(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, price1)

	require.Equal(t, price.Event, price1.Event)
	require.Equal(t, price1.Price, arg.Price)
	require.Equal(t, price.Class, price1.Class)
}

func TestSoftDeletePrice(t *testing.T) {
	event := CreateRandomPrice(t)
	testQueries.SoftDeletePrice(context.Background(), event.ID)

	event1, err := testQueries.GetPrice(context.Background(), event.ID)
	require.NoError(t, err)
	require.NotEmpty(t, event1)
	require.NotEmpty(t, event1.DeletedAt.Time)
}

func TestPermaDeletePrice(t *testing.T) {
	event := CreateRandomPrice(t)
	testQueries.PermaDeletePrice(context.Background(), event.ID)

	_, err := testQueries.GetPrice(context.Background(), event.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}
