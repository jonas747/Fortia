package world

import (
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/messages"
	"github.com/jonas747/fortia/vec"
)

const (
	ErrCodeEmpty = iota
	ErrCodeConversionErr
	ErrCodeNotFound
)

// Info about the world
type WorldInfo struct {
	Size        int         // The size of the world in chunks
	ChunkHeight int         // The height of thw world in blocks
	ChunkWidth  int         // Size of a layer in blocks
	BlockTypes  []BlockType // ..
	Biomes      BiomesInfo  // ..
}

// Represents a user
type UserInfo struct {
	Name string // Name of the user
}

// The way world information is stored may be changed in the future
type GameDB interface {
	GetWorldInfo() (*WorldInfo, ferr.FortiaError)  // Returns info about the world
	SetWorldInfo(info *WorldInfo) ferr.FortiaError // Saves world information to the database

	GetUserInfo(user string) (*UserInfo, ferr.FortiaError)
	SetUserInfo(info *UserInfo) ferr.FortiaError

	GetUserEntities(user string) ([]int, ferr.FortiaError)         // Returns the users owned entities
	EditUserEntities(user string, add, del []int) ferr.FortiaError // Adds and removes entities from the users owned list

	GetChunk(pos vec.Vec2I) (*messages.Chunk, ferr.FortiaError)
	SetChunk(chunk *messages.Chunk) ferr.FortiaError

	PopAction(tick int) (*Action, ferr.FortiaError) // Returns a action, errcodeempty if none
	IncrTick() (int, ferr.FortiaError)              // Increases the tick counter, and returns the new ticknumber
	// Todo
	//GetEntity()
	//SetEntity()
}
