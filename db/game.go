package db

import (
	"github.com/jonas747/fortia/errors"
	"github.com/jonas747/fortia/messages"
	"github.com/jonas747/fortia/vec"
)

// Represents a user
type GameUserInfo struct {
	Name string // Name of the user
}

// The way world information is stored may be changed in the future
type GameDB interface {
	GetWorldSettings() (*messages.WorldSettings, errors.FortiaError)  // Returns info about the world
	SetWorldSettings(info *messages.WorldSettings) errors.FortiaError // Saves world information to the database

	GetUserInfo(user string) (*GameUserInfo, errors.FortiaError)
	SetUserInfo(info *GameUserInfo) errors.FortiaError

	GetUserEntities(user string) ([]int, errors.FortiaError)         // Returns the users owned entities
	EditUserEntities(user string, add, del []int) errors.FortiaError // Adds and removes entities from the users owned list

	GetChunk(pos vec.Vec2I) (*messages.Chunk, errors.FortiaError)
	SetChunk(chunk *messages.Chunk) errors.FortiaError

	PushAction(action *messages.Action, tick int) errors.FortiaError // Pushes a new action to the db to be processed at tick
	PopAction(tick, kind int) (*messages.Action, errors.FortiaError) // Returns a action with the given type, errcodeempty if none
	IncrTick() (int, errors.FortiaError)                             // Increases the tick counter, and returns the new ticknumber
	GetCurrentTick() (int, errors.FortiaError)                       // Returns the current tick
	// Todo
	//GetEntity()
	//SetEntity()
}
