package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/amirul-zafrin/event/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomEvent(t *testing.T) Event {
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
	CreateRandomEvent(t)
}

func TestGetEvent(t *testing.T) {
	event1 := CreateRandomEvent(t)
	event2, err := testQueries.GetEvent(context.Background(), event1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, event2)

	require.Equal(t, event1.Name, event2.Name)
	require.Equal(t, event1.Location, event1.Location)
	require.Equal(t, event1.Date, event2.Date)

	require.WithinDuration(t, event1.CreatedAt, event2.CreatedAt, time.Second)
}

func TestListAllEvent(t *testing.T) {
	for i := 0; i < 5; i++ {
		CreateRandomEvent(t)
	}

	arg := ListAllEventsParams{
		Limit:  5,
		Offset: 0,
	}

	events, err := testQueries.ListAllEvents(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, events, 5)

	for _, event := range events {
		require.NotEmpty(t, event)
	}

}

func TestUpdateEvent(t *testing.T) {
	event1 := CreateRandomEvent(t)

	arg := UpdateEventParams{
		ID:   event1.ID,
		Name: util.RandomEventName(),
	}
	event2, err := testQueries.UpdateEvent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, event2)
	require.True(t, event2.UpdatedAt.Valid)

	require.Equal(t, arg.Name, event2.Name)
	require.Equal(t, event1.Location, event1.Location)
	require.Equal(t, event1.Date, event2.Date)

	require.WithinDuration(t, event1.CreatedAt, event2.CreatedAt, time.Second)
}

func TestSoftDelete(t *testing.T) {
	event := CreateRandomEvent(t)

	testQueries.SoftDeleteEvent(context.Background(), event.ID)

	event1, err := testQueries.GetEvent(context.Background(), event.ID)
	require.NoError(t, err)
	require.NotEmpty(t, event1.Name)
	require.NotEmpty(t, event1.DeletedAt)
}

func TestPermaDelete(t *testing.T) {
	event := CreateRandomEvent(t)

	testQueries.PermaDeleteEvent(context.Background(), event.ID)

	_, err := testQueries.GetEvent(context.Background(), event.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}
