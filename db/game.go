package db

import (
	ferr "github.com/jonas747/fortia/error"
)

type GameDB struct {
	*Database
}

// Returns the specified user's info
func (g *GameDB) GetUserInfo(user string) (map[string]string, ferr.FortiaError) {
	// TODO
	return make(map[string]string), nil
}

// Sets the specified users info fields from info map to whatever is in the info map
func (g *GameDB) SetUserInfo(user string, info map[string]interface{}) ferr.FortiaError {
	// TODO
	return nil
}
