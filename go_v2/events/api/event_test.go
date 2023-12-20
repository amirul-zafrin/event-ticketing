package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	mockdb "github.com/amirul-zafrin/event/db/mock"
	db "github.com/amirul-zafrin/event/db/sqlc"
	"github.com/amirul-zafrin/event/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	event := RandomEvent()

	testCases := []struct {
		name          string
		eventID       int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *http.Response)
	}{
		{
			name:    "OK",
			eventID: event.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetEvent(gomock.Any(), gomock.Eq(event.ID)).
					Times(1).
					Return(event, nil)
			},
			checkResponse: func(t *testing.T, recorder *http.Response) {
				require.Equal(t, http.StatusOK, recorder.StatusCode)
				requireBodyMatch(t, &recorder.Body, event)
			},
		},
		{
			name:    "NotFound",
			eventID: event.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetEvent(gomock.Any(), gomock.Eq(event.ID)).
					Times(1).
					Return(db.Event{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *http.Response) {
				require.Equal(t, http.StatusNotFound, recorder.StatusCode)
			},
		},
		{
			name:    "InternalError",
			eventID: event.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetEvent(gomock.Any(), gomock.Eq(event.ID)).
					Times(1).
					Return(db.Event{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *http.Response) {
				require.Equal(t, http.StatusInternalServerError, recorder.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server := NewServer(store)

			url := fmt.Sprintf("/event/%d", tc.eventID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			recorder, err := server.router.Test(request)
			require.NoError(t, err)
			tc.checkResponse(t, recorder)
		})
	}
}

func RandomEvent() db.Event {
	return db.Event{
		ID:        util.RandomInt(1, 1000),
		CreatedAt: util.RandomDate(),
		Name:      util.RandomEventName(),
		Date:      sql.NullTime{Time: util.RandomDate(), Valid: true},
		Location:  util.RandomLocation(),
		Capacity:  int32(util.RandomInt(1000, 10000)),
	}
}

type EventWrapper struct {
	Message db.Event `json:"message"`
	Status  string   `json:"status"`
}

func requireBodyMatch(t *testing.T, body *io.ReadCloser, event db.Event) {
	data, err := io.ReadAll(*body)
	require.NoError(t, err)
	var gotEvent EventWrapper
	err = json.Unmarshal(data, &gotEvent)
	if !gotEvent.Message.Seats.Valid {
		gotEvent.Message.Seats.RawMessage = nil
	}
	require.NoError(t, err)
	require.Equal(t, event, gotEvent.Message)
}
