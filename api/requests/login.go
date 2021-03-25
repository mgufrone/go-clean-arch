package requests

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Login struct {
	Token string `json:"token"`
	Agent string `header:"user-agent"`
}

func (l *Login) Validate(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, l,
		validation.Field(&l.Token, validation.Required, validation.NotNil),
		validation.Field(&l.Agent, validation.Required, validation.NotNil),
	)
}
