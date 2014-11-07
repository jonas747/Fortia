package world

import (
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/vec"
)

// Info about the world
type WorldInfo struct {
	Size       int         // The size of the world in chunks
	Height     int         // The height of thw world in blocks
	LayerSize  int         // Size of a layer in blocks
	BlockTypes []BlockType // ..
	Biomes     BiomesInfo  // ..
}

// TODO
type UserInfo struct {
}

type GameDB interface {
	GetWorldInfo() (*WorldInfo, ferr.FortiaError)  // Returns info about the world
	SetWorldInfo(info *WorldInfo) ferr.FortiaError // Saves world information to the database

	GetUserInfo(user string) (*UserInfo, ferr.FortiaError)
	SetUserInfo(info *UserInfo) ferr.FortiaError

	GetLayer(pos vec.Vec3I) (*Layer, ferr.FortiaError)
	SetLayer(layer *Layer) ferr.FortiaError
	GetLayers(positions []vec.Vec3I) ([]*Layer, ferr.FortiaError)
	SetLayers(layers []*Layer) ferr.FortiaError

	GetChunkInfo(pos vec.Vec2I) (*Chunk, ferr.FortiaError)
	SetChunkInfo(chunk *Chunk) ferr.FortiaError

	GetEntity()
	SetEntity()
}
