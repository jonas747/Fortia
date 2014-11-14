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
	WorldGen    WorldGenerator
	GeneralInfo *WorldInfo
}

func (w *World) NewGenerator() *Generator {
	return NewGenerator(w, &w.GeneralInfo.Biomes, w.GeneralInfo.BlockTypes, map[string]int64{"landscape": 8})
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
