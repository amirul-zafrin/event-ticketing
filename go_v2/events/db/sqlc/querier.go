// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"context"
)

type Querier interface {
	CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error)
	CreatePrice(ctx context.Context, arg CreatePriceParams) (Price, error)
	GetEvent(ctx context.Context, id int64) (Event, error)
	GetPrice(ctx context.Context, id int64) (Price, error)
	ListAllEvents(ctx context.Context, arg ListAllEventsParams) ([]Event, error)
	ListAllPrices(ctx context.Context, arg ListAllPricesParams) ([]Price, error)
	PermaDeleteEvent(ctx context.Context, id int64) error
	PermaDeletePrice(ctx context.Context, id int64) error
	SoftDeleteEvent(ctx context.Context, id int64) error
	SoftDeletePrice(ctx context.Context, id int64) error
	UpdateEvent(ctx context.Context, arg UpdateEventParams) (Event, error)
	UpdateEventSeat(ctx context.Context, arg UpdateEventSeatParams) (Event, error)
	UpdatePrice(ctx context.Context, arg UpdatePriceParams) (Price, error)
}

var _ Querier = (*Queries)(nil)
