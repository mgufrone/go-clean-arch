package stats

import (
	"context"
	"recipes/domains/shared"
	"recipes/domains/user"
)

type View struct {
	shared.Model
	User        *user.User `json:"user"`
	Reference   string `json:"reference"`
	ReferenceID uint64 `json:"reference_id"`
}

type ViewWithStats struct {
	*View
	Count int64 `json:"count"`
}


type IViewRepository interface {
	GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*View) ([]*View, *shared.CommonLimiter, error)
	GroupCount(ctx context.Context, limiter *shared.CommonLimiter, filters ...*View) ([]*ViewWithStats, *shared.CommonLimiter, error)
	CountByReference(ctx context.Context, reference string, id uint64) (int64, error)
	Count(ctx context.Context, filters ...*View) (int64, error)
	Put(ctx context.Context, view *View) error
}