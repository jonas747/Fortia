package resterrors

import (
	"fmt"
)

// Returns json with an error
func ApiError(code int, msg string) []byte {
	return []byte(fmt.Sprintf("{\"error\":\"%s\", \"errorcode\": %d}", msg, code))
}

// Not using iota because when the order changes
// all the clients would break
var (
	// Auth
	// Login
	ErrWrongLoginDetails = 1 // Username or password is incorrect
	// Misc
	ErrInvalidSessionCookie = 2 // The session has expired
	ErrServerError          = 7 // Seomthing went wrong on the server side
	// Register
	ErrInvalidUsername = 3 // The username is invalid
	ErrInvalidEmail    = 4 // ...
	ErrInavlidPassword = 5 // ...
	ErrUsernameTaken   = 6 // The username is allready taken

	// Game
	ErrNotRegisteredWorld = 8 // User not registered to specified world
)
