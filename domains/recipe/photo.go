package recipe

import (
	"context"
	"recipes/domains/shared"
)

type Photo struct {
	shared.Model
	Sequence uint `json:"sequence"`
	URL string `json:"url"`
}

type IPhotoRepository interface {
	GetBySequence(ctx context.Context, sequence uint) (*Photo, error)
	DeleteSequence(ctx context.Context, sequence uint) error
	Update(ctx context.Context, sequence uint, photo *Photo) error
	Create(ctx context.Context, photo *Photo) error
}