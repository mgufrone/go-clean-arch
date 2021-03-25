package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"recipes/api/constants"
	"recipes/app/common"
)

type Request interface {
	Validate(ctx context.Context) error
}

func Sanitize(out Request) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		err := common.Try(func() error {
			return ctx.BindHeader(out)
		}, func() error {
			return ctx.BindQuery(out)
		}, func() error {
			return ctx.BindUri(out)
		},func() error {
			return ctx.Bind(out)
		})
		if err != nil {
			BadRequest(ctx, err)
			return
		}
		ctx.Set(constants.RequestKey, out)
	}
}
