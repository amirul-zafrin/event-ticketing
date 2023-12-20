package api

import (
	"database/sql"
	"fmt"
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	//builds stubs
	store.EXPECT().
		GetEvent(gomock.Any(), gomock.Eq(event.ID)).
		Times(1).
		Return(event, nil)

	//Start test
	server := NewServer(store)

	url := fmt.Sprintf("/event/%d", event.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	recorder, err := server.router.Test(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, recorder.StatusCode)
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

// ID        int64                 `json:"id"`
//     CreatedAt time.Time             `json:"created_at"`
//     UpdatedAt sql.NullTime          `json:"updated_at"`
//     DeletedAt sql.NullTime          `json:"deleted_at"`
//     Name      string                `json:"name"`
//     Date      sql.NullTime          `json:"date"`
//     Location  string                `json:"location"`
//     Capacity  int32                 `json:"capacity"`
//     Seats     pqtype.NullRawMessage `json:"seats"`
