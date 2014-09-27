package main

import (
	"fmt"
)

// Returns json with an error
func ApiError(code int, msg string) []byte {
	return []byte(fmt.Sprintf("{\"error\":\"%s\", \"errorcode\": %d}", msg, code))
}

// Not using iota incase the order is changed
var (
	// Login
	ErrCodeWrongLoginDetails = 1 // Username or password is incorrect
	// Misc
	ErrCodeInvalidSessionCookie = 2 // The session has expired
	ErrCodeServerError          = 7 // Seomthing went wrong on the server side
	// Register
	ErrCodeInvalidUsername = 3 // The username is invalid
	ErrCodeInvalidEmail    = 4
	ErrCodeInavlidPassword = 5
	ErrCodeUsernameTaken   = 6
)

var (
	ErrWrongLoginDetails    = ApiError(ErrCodeWrongLoginDetails, "Bad password or username")
	ErrInvalidSessionCookie = ApiError(ErrCodeInvalidSessionCookie, "Session has expired")
	ErrServerError          = ApiError(ErrCodeServerError, "Something went wrong on the server side")
)
