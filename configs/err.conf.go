package configs

import "fmt"

type ErrorCode int

const (
	// Lỗi do user (số dương)
	ErrorCode_SIGN_IN_MISSING_FIELDS     ErrorCode = 1000
	ErrorCode_MISSING_REQUIRED_FIELDS    ErrorCode = 1001
	ErrorCode_USER_ALREADY_EXISTS        ErrorCode = 1002
	ErrorCode_USER_NOT_FOUND             ErrorCode = 1003
	ErrorCode_APP_NOT_FOUND              ErrorCode = 2001
	ErrorCode_USER_NOT_ACTIVE            ErrorCode = 1004
	ErrorCode_AUTH_FAILED                ErrorCode = 1005
	ErrorCode_INVALID_TOKEN              ErrorCode = 1006
	ErrorCode_INVALID_USER_ID_IN_TOKEN   ErrorCode = 1007
	ErrorCode_SIGN_UP_MISSING_FIELDS     ErrorCode = 1008
	ErrorCode_INVALID_REQUEST            ErrorCode = 1009
	ErrorCode_MISSING_AUTH_HEADER        ErrorCode = 1010
	ErrorCode_INVALID_AUTH_HEADER_FORMAT ErrorCode = 1011

	// Lỗi hệ thống (số âm)
	ErrorCode_FAILED_TO_HASH_PASSWORD  ErrorCode = -1001
	ErrorCode_DATABASE_ERROR           ErrorCode = -1002
	ErrorCode_FAILED_TO_CREATE_USER    ErrorCode = -1003
	ErrorCode_FAILED_TO_GENERATE_TOKEN ErrorCode = -1004
	ErrorCode_FAILED_TO_INSERT_USER    ErrorCode = -1005
)

var errorMessages = map[ErrorCode]string{
	// User errors
	ErrorCode_SIGN_IN_MISSING_FIELDS:     "SIGN_IN_MISSING_FIELDS",
	ErrorCode_MISSING_REQUIRED_FIELDS:    "MISSING_REQUIRED_FIELDS",
	ErrorCode_USER_ALREADY_EXISTS:        "USER_ALREADY_EXISTS",
	ErrorCode_USER_NOT_FOUND:             "USER_NOT_FOUND",
	ErrorCode_APP_NOT_FOUND:              "APP_NOT_FOUND",
	ErrorCode_USER_NOT_ACTIVE:            "USER_NOT_ACTIVE",
	ErrorCode_AUTH_FAILED:                "AUTH_FAILED",
	ErrorCode_INVALID_TOKEN:              "INVALID_TOKEN",
	ErrorCode_INVALID_USER_ID_IN_TOKEN:   "INVALID_USER_ID_IN_TOKEN",
	ErrorCode_SIGN_UP_MISSING_FIELDS:     "SIGN_UP_MISSING_FIELDS",
	ErrorCode_INVALID_REQUEST:            "INVALID_REQUEST",
	ErrorCode_MISSING_AUTH_HEADER:        "MISSING_AUTH_HEADER",
	ErrorCode_INVALID_AUTH_HEADER_FORMAT: "INVALID_AUTH_HEADER_FORMAT",

	// System errors
	ErrorCode_FAILED_TO_HASH_PASSWORD:  "FAILED_TO_HASH_PASSWORD",
	ErrorCode_DATABASE_ERROR:           "DATABASE_ERROR",
	ErrorCode_FAILED_TO_CREATE_USER:    "FAILED_TO_CREATE_USER",
	ErrorCode_FAILED_TO_GENERATE_TOKEN: "FAILED_TO_GENERATE_TOKEN",
	ErrorCode_FAILED_TO_INSERT_USER:    "FAILED_TO_INSERT_USER",
}

func GetErrString(code ErrorCode) string {
	if msg, ok := errorMessages[code]; ok {
		return fmt.Sprintf("%d:%s", code, msg)
	}
	return fmt.Sprintf("%d:UNKNOWN_ERROR", code)
}

// Nếu vẫn cần struct Error cho JSON response
type Error struct {
	Code   string `json:"code"`
	Number int    `json:"number"`
}

func (e ErrorCode) ToError() Error {
	if msg, ok := errorMessages[e]; ok {
		return Error{Code: msg, Number: int(e)}
	}
	return Error{Code: "UNKNOWN_ERROR", Number: int(e)}
}
