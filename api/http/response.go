package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DataResponse struct {
	Data interface{} `json:"data"`
}
type DataWithTotalResponse struct {
	DataResponse
	Total uint `json:"total"`
}

func BadRequest(ctx *gin.Context, err error)  {
	if b, err2 := json.Marshal(err); err2 == nil && string(b) != "{}" {
		var c gin.H
		_ = json.Unmarshal(b, &c)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": c})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
func ServerError(ctx *gin.Context, err error)  {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	panic(err)
}
func Unauthorized(ctx *gin.Context, err error)  {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	panic(err)
}
func NotFound(ctx *gin.Context)  {
	ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found"})
}

func OkWithTotal(ctx *gin.Context, result interface{}, total uint)  {
	ctx.JSON(http.StatusOK, &DataWithTotalResponse{
		DataResponse: DataResponse{
			Data: result,
		},
		Total:        total,
	})
}

func Ok(ctx *gin.Context, result interface{}) {
	ctx.JSON(http.StatusOK, &DataResponse{Data: result})
}

func Deleted(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, nil)
}
