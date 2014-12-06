// Package world contains world realted stuff
// TODO: Check if layer is air and mark it if so
package world

import (
	//"github.com/jonas747/fortia/rdb"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/vec"
)

type World struct {
	Logger      *log.LogClient
	Db          GameDB
	GeneralInfo *WorldInfo
}

func (w *World) NewGenerator() *Generator {
	return NewGenerator(w, &w.GeneralInfo.Biomes, w.GeneralInfo.BlockTypes, 1)
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

func (w *World) ChunkToWorldPos(chunkPos vec.Vec2I) vec.Vec2I {
	chunkPos.Multiply(vec.Vec2I{X: w.GeneralInfo.ChunkWidth, Y: w.GeneralInfo.ChunkWidth})
	return chunkPos
}

func (w *World) CoordsToIndex(pos vec.Vec3I) int {
	return pos.X + w.GeneralInfo.ChunkWidth*(pos.Y+w.GeneralInfo.ChunkWidth*pos.Z)
}

func (w *World) IndexToCoords(index int) vec.Vec3I {
	x := index % w.GeneralInfo.ChunkWidth
	y := (index / w.GeneralInfo.ChunkWidth) % w.GeneralInfo.ChunkWidth
	z := index / (w.GeneralInfo.ChunkWidth * w.GeneralInfo.ChunkWidth)
	return vec.Vec3I{x, y, z}
}

type Entity interface {
	GetPosition()
	GetId() int
}
