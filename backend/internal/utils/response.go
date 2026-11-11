package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

type (
	response[T any] struct {
		Status     int         `json:"status"`
		Message    string      `json:"message"`
		Data       T           `json:"data,omitempty"`
		Metadata   metadata    `json:"metadata"`
		Pagination *Pagination `json:"pagination,omitempty"`
	}

	errResponse struct {
		Status   int      `json:"status"`
		Error    err      `json:"error"`
		Metadata metadata `json:"metadata"`
	}

	metadata struct {
		RequestID string    `json:"request_id"`
		Timestamp time.Time `json:"timestamp"`
	}

	Pagination struct {
		Page  int   `json:"page"`
		Size  int   `json:"size"`
		Total int64 `json:"total"`
	}

	err struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
)

func SetResponse[T any](ctx *gin.Context, status int, message string, data T, p *Pagination) {
	ctx.JSON(
		status,
		response[T]{
			Status:  status,
			Message: message,
			Data:    data,
			Metadata: metadata{
				RequestID: CtxRequestID(ctx.Request.Context()),
				Timestamp: time.Now().UTC(),
			},
			Pagination: p,
		},
	)
}

func SetErrResponse(ctx *gin.Context, status int, code, message string) {
	ctx.JSON(
		status,
		errResponse{
			Status: status,
			Error: err{
				Code:    code,
				Message: message,
			},
			Metadata: metadata{
				RequestID: CtxRequestID(ctx.Request.Context()),
				Timestamp: time.Now().UTC(),
			},
		},
	)
}
