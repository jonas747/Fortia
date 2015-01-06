package db

import (
	"fmt"
	"github.com/jonas747/fortia/errorcodes"
	"github.com/jonas747/fortia/errors"
)

// To be used when there is an error with the database
type DBError struct {
	*errors.FortiaBaseError
	Query string
}

func NewDBError(query string, code errorcodes.ErrorCode, format string, a ...interface{}) errors.FortiaError {
	base := errors.New(code, format, a...).(*errors.FortiaBaseError)
	return &DBError{
		FortiaBaseError: base,
		Query:           query,
	}
}

func (e *DBError) Error() string {
	header := fmt.Sprintf("Database Error on query \"%s\"\n", e.Query)
	body := errors.DefaultError(e)
	return header + body
}

type AuthError struct {
	*errors.FortiaBaseError
	User string
}

func NewAuthError(user string, code errorcodes.ErrorCode, format string, a ...interface{}) errors.FortiaError {
	base := errors.New(code, format, a...).(*errors.FortiaBaseError)
	return &AuthError{
		FortiaBaseError: base,
		User:            user,
	}
}

// Auth error is to be used on authentication errors (invalid credentials, expired session cookies etc...)
func (e *AuthError) Error() string {
	msg := fmt.Sprintf("Auth Error on user [%s] Code [%d] Message: [%s]", e.User, e.Code, e.GetMessage())
	return msg
}
