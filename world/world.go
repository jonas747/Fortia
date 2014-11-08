// Package world contains world realted stuff
// TODO: Check if layer is air and mark it if so
package world

import (
	//"github.com/jonas747/fortia/rdb"
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
	Db          GameDB
	WorldGen    WorldGenerator
	GeneralInfo *WorldInfo
}

func (w *World) NewGenerator() *Generator {
	return NewGenerator(w, &w.GeneralInfo.Biomes, w.GeneralInfo.BlockTypes, map[string]int64{"landscape": 123089})
}

func (w *World) LoadSettingsFromDb() ferr.FortiaError {
	info, err := w.Db.GetWorldInfo()
	if err != nil {
		return err
	}
	w.GeneralInfo = info

	return nil
}

func (w *World) SaveSettingsToDb() ferr.FortiaError {
	err := w.Db.SetWorldInfo(w.GeneralInfo)
	return err
}

func (w *World) LayerToWorldPos(layePos vec.Vec3I) vec.Vec3I {
	lw := layePos.Clone()
	lw.Multiply(vec.Vec3I{X: w.GeneralInfo.LayerSize, Y: w.GeneralInfo.LayerSize})
	return lw
}

func (w *World) GenLayer(pos vec.Vec3I) *Layer {
	//chunkWorldPos := w.ChunkToWorldPos(pos)

	layer := &Layer{
		Position: pos,
	}

	blocks := make([]*Block, w.GeneralInfo.LayerSize*w.GeneralInfo.LayerSize)

	for x := 0; x < w.GeneralInfo.LayerSize; x++ {
		for y := 0; y < w.GeneralInfo.LayerSize; y++ {

			id := rand.Intn(2) + 1
			if layer.Position.Z > w.GeneralInfo.Height/2 {
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
	return w.GeneralInfo.LayerSize*pos.X + pos.Y
}

// Return a blocks x and y from the index in the layer slice
// x = index / size
// y = index - (x * size)
func (w *World) IndexToCoords(index int) vec.Vec3I {
	x := index / w.GeneralInfo.LayerSize
	y := index - (x * w.GeneralInfo.LayerSize)
	return vec.Vec3I{x, y, 0}
}

type Entity interface {
	GetPosition()
	GetId() int
}

type WorldGenerator interface {
	GenerateBLock(x, y, z int) int
}
