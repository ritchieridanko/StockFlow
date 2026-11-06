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
	CodeInternal         errCode = "INTERNAL"
	CodeInvalidPayload   errCode = "INVALID_PAYLOAD"
	CodeSKUAlreadyExists errCode = "SKU_ALREADY_EXISTS"
	CodeUnknown          errCode = "UNKNOWN"
)

// External Error Messages
const (
	MsgInternalServer   string = "Internal server error"
	MsgInvalidPayload   string = "Invalid payload"
	MsgSKUAlreadyExists string = "SKU is already registered"
)
