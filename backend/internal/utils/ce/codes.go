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
	CodeInsufficientStock     errCode = "INSUFFICIENT_STOCK"
	CodeInternal              errCode = "INTERNAL"
	CodeInvalidParam          errCode = "INVALID_PARAM"
	CodeInvalidPayload        errCode = "INVALID_PAYLOAD"
	CodeOrderAlreadyCancelled errCode = "ORDER_ALREADY_CANCELLED"
	CodeOrderNotFound         errCode = "ORDER_NOT_FOUND"
	CodeProductHasHistory     errCode = "PRODUCT_HAS_HISTORY"
	CodeProductNotFound       errCode = "PRODUCT_NOT_FOUND"
	CodeSKUAlreadyExists      errCode = "SKU_ALREADY_EXISTS"
	CodeUnknown               errCode = "UNKNOWN"
)

// External Error Messages
const (
	MsgEmptyOrderItems       string = "Order items cannot be empty"
	MsgInsufficientStock     string = "Insufficient stock"
	MsgInternalServer        string = "Internal server error"
	MsgInvalidParam          string = "Invalid param"
	MsgInvalidPayload        string = "Invalid payload"
	MsgOrderAlreadyCancelled string = "Order is already cancelled"
	MsgOrderNotFound         string = "Order is not found"
	MsgProductHasHistory     string = "Unable to delete product with stock movements"
	MsgProductNotFound       string = "Product is not found"
	MsgSKUAlreadyExists      string = "SKU is already registered"
)
