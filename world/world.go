// Package world contains world realted stuff
// TODO: Check if layer is air and mark it if so
package world

import (
	"encoding/json"
	"github.com/jonas747/fortia/db"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/vec"
	"math/rand"
)

type BlockFlag byte

const (
	BlockConnectedGround BlockFlag = 1 << iota
	BlockOccupiedFull
	BlockOccupiedHalf
)

type World struct {
	Logger      *log.LogClient
	Db          *db.GameDB
	LayerSize   int
	LayerHeight int // How many layers high the world is
	WorldGen    WorldGenerator
}

func (w *World) LayerToWorldPos(layePos vec.Vec3I) vec.Vec3I {
	lw := layePos.Clone()
	lw.Multiply(vec.Vec3I{X: w.LayerSize, Y: w.LayerSize})
	return lw
}

func (w *World) GenLayer(pos vec.Vec3I) *Layer {
	//chunkWorldPos := w.ChunkToWorldPos(pos)

	layer := &Layer{
		Position: pos,
	}

	blocks := make([]*Block, w.LayerSize*w.LayerSize)

	for x := 0; x < w.LayerSize; x++ {
		for y := 0; y < w.LayerSize; y++ {

			id := rand.Intn(2) + 1
			if layer.Position.Z > w.LayerHeight/2 {
				id = 0
				layer.IsAir = true
			}

			b := Block{
				LocalPosition: vec.Vec2I{x, y},
				Layer:         layer,
				Id:            id,
			}
			blocks[w.CoordsToIndex(vec.Vec3I{x, y, 0})] = &b
		}
	}

	layer.Blocks = blocks
	layer.World = w
	return layer
}

// index = size * x + y
func (w *World) CoordsToIndex(pos vec.Vec3I) int {
	return w.LayerSize*pos.X + pos.Y
}

// Return a blocks x and y from the index in the layer slice
// x = index / size
// y = index - (x * size)
func (w *World) IndexToCoords(index int) vec.Vec3I {
	x := index / w.LayerSize
	y := index - (x * w.LayerSize)
	return vec.Vec3I{x, y, 0}
}

func (w *World) SetLayer(layer *Layer) ferr.FortiaError {
	raw, err := layer.Json()
	if err != nil {
		return err
	}
	return w.Db.SetLayer(layer.Position.X, layer.Position.Y, layer.Position.Z, raw)
}

func (w *World) GetLayer(pos vec.Vec3I) (*Layer, ferr.FortiaError) {
	rawLayer, err := w.Db.GetLayer(pos.X, pos.Y, pos.Z)
	if err != nil {
		return nil, err
	}

	layer := &Layer{
		Position: pos,
	}

	nErr := json.Unmarshal(rawLayer, layer)
	if nErr != nil {
		return nil, ferr.Wrap(nErr, "")
	}

	return layer, nil
}

type Layer struct {
	Position vec.Vec3I
	Flags    int
	World    *World `json:"-"`
	Blocks   []*Block
	IsAir    bool // True if this layer is just air
}

func (l *Layer) Json() ([]byte, ferr.FortiaError) {
	out, err := json.Marshal(l)
	if err != nil {
		return []byte{}, ferr.Wrap(err, "")
	}
	return out, nil
}

type Block struct {
	LocalPosition vec.Vec2I `json:"-"`
	Layer         *Layer    `json:"-"`
	Entities      []int     `json:",omitempty"`
	Id            int
	Flags         byte                   `json:",omitempty"`
	Properties    map[string]interface{} `json:",omitempty"`
}

type Entity interface {
	GetPosition()
	GetId() int
}

type WorldGenerator interface {
	GenerateBLock(x, y, z int) int
}
