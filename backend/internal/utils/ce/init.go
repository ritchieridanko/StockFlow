package ce

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/stockflow/backend/internal/infra/logger"
)

type errCode string

type Error struct {
	code    errCode
	message string
	err     error
	fields  []logger.Field
}

func NewError(ec errCode, message string, err error, fields ...logger.Field) *Error {
	return &Error{
		code:    ec,
		message: message,
		err:     err,
		fields:  fields,
	}
}

func (e *Error) Code() errCode {
	return e.code
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Error() string {
	if e.err != nil {
		return e.message + ": " + e.err.Error()
	}
	return e.message
}

func (e *Error) Fields() []logger.Field {
	return e.fields
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Bind(ctx *gin.Context) {
	ctx.Error(e)
}

func (e *Error) ToHTTPErrCode() int {
	switch e.code {
	case
		CodeInvalidParam,
		CodeInvalidPayload:
		return http.StatusBadRequest
	case CodeProductNotFound:
		return http.StatusNotFound
	case
		CodeInsufficientStock,
		CodeProductHasHistory,
		CodeSKUAlreadyExists:
		return http.StatusConflict
	case
		CodeInternal,
		CodeUnknown:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
