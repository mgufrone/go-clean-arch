package stats

import (
	"context"
	"recipes/domains/shared"
	"recipes/domains/user"
	"time"
)

type Like struct {
	shared.Model
	User        *user.User `json:"user"`
	CreatedAt   time.Time `json:"created_at"`
	Reference   string `json:"reference"`
	ReferenceID uint64 `json:"reference_id"`
}

type ILikeRepository interface {
	CountByReference(ctx context.Context, reference string, id uint64) (int64, error)
	Count(ctx context.Context, filters ...*Like) (int64, error)
	Put(ctx context.Context, like *Like) error
	GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*Like) ([]*Like, *shared.CommonLimiter, error)
	Delete(ctx context.Context, like *Like) error
}
