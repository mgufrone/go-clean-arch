package recipe

import (
	"context"
	"recipes/domains/shared"
)

type Step struct {
	shared.Model
	Sequence uint `json:"sequence"`
	Step string `json:"step"`
}

type IStepRepository interface {
	GetBySequence(ctx context.Context, sequence uint) (*Step, error)
	DeleteSequence(ctx context.Context, sequence uint) error
	Update(ctx context.Context, sequence uint, step *Step) error
	Create(ctx context.Context, photo *Step) error
}
