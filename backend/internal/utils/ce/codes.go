package ce

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

// Internal Errors
var (
	ErrDBAffectNoRows error = errors.New("no rows affected")
	ErrDBQueryNoRows  error = pgx.ErrNoRows
)

// Internal Error Codes
const (
	CodeInsufficientStock errCode = "INSUFFICIENT_STOCK"
	CodeInternal          errCode = "INTERNAL"
	CodeInvalidParam      errCode = "INVALID_PARAM"
	CodeInvalidPayload    errCode = "INVALID_PAYLOAD"
	CodeProductHasHistory errCode = "PRODUCT_HAS_HISTORY"
	CodeProductNotFound   errCode = "PRODUCT_NOT_FOUND"
	CodeSKUAlreadyExists  errCode = "SKU_ALREADY_EXISTS"
	CodeUnknown           errCode = "UNKNOWN"
)

// External Error Messages
const (
	MsgInsufficientStock string = "Insufficient stock"
	MsgInternalServer    string = "Internal server error"
	MsgInvalidParam      string = "Invalid param"
	MsgInvalidPayload    string = "Invalid payload"
	MsgProductHasHistory string = "Unable to delete product with stock movements"
	MsgProductNotFound   string = "Product is not found"
	MsgSKUAlreadyExists  string = "SKU is already registered"
)
