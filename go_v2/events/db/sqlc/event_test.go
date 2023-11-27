package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/amirul-zafrin/event/util"
	"github.com/stretchr/testify/require"
)

func createRandomEvent(t *testing.T) Event {
	arg := CreateEventParams{
		Name:     util.RandomEventName(),
		Date:     sql.NullTime{Time: util.RandomDate(), Valid: true},
		Location: util.RandomLocation(),
		Capacity: int32(util.RandomInt(1000, 10000)),
	}

	event, err := testQueries.CreateEvent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, event)

	require.Equal(t, arg.Name, event.Name)
	require.Equal(t, arg.Location, event.Location)

	require.NotZero(t, event.ID)
	require.NotZero(t, event.CreatedAt)
	require.NotZero(t, event.Date)

	return event
}
func TestCreateEvent(t *testing.T) {
	createRandomEvent(t)
}

func TestGetAccount(t *testing.T) {
	event1 := createRandomEvent(t)
	event2, err := testQueries.GetEvent(context.Background(), event1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, event2)

	require.Equal(t, event1.Name, event2.Name)
	require.Equal(t, event1.Location, event1.Location)
	require.Equal(t, event1.Date, event2.Date)

	require.WithinDuration(t, event1.CreatedAt, event2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	event1 := createRandomEvent(t)

	arg := UpdateEventParams{
		ID:   event1.ID,
		Name: util.RandomEventName(),
	}
	event2, err := testQueries.UpdateEvent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, event2)

	require.Equal(t, arg.Name, event2.Name)
	require.Equal(t, event1.Location, event1.Location)
	require.Equal(t, event1.Date, event2.Date)

	require.WithinDuration(t, event1.CreatedAt, event2.CreatedAt, time.Second)
}
