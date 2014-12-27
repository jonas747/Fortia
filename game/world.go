// Package world contains world realted stuff
// TODO: Check if layer is air and mark it if so
package game

import (
	"github.com/jonas747/fortia/db"
	"github.com/jonas747/fortia/errors"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/messages"
	"github.com/jonas747/fortia/vec"
)

type World struct {
	Logger   *log.LogClient
	Db       db.GameDB
	Settings *messages.WorldSettings
}

func (w *World) NewGenerator() *Generator {
	return NewGenerator(w, 1)
}

func (w *World) LoadSettingsFromDb() errors.FortiaError {
	info, err := w.Db.GetWorldSettings()
	if err != nil {
		return err
	}
	w.Settings = info

	return nil
}

func (w *World) SaveSettingsToDb() errors.FortiaError {
	err := w.Db.SetWorldSettings(w.Settings)
	return err
}

func (w *World) ChunkToWorldPos(chunkPos vec.Vec2I) vec.Vec2I {
	chunkPos.Multiply(vec.Vec2I{X: int(w.Settings.GetChunkWidth()), Y: int(w.Settings.GetChunkWidth())})
	return chunkPos
}

func (w *World) CoordsToIndex(pos vec.Vec3I) int {
	return pos.X + int(w.Settings.GetChunkWidth())*(pos.Y+int(w.Settings.GetChunkWidth())*pos.Z)
}

func (w *World) IndexToCoords(index int) vec.Vec3I {
	x := index % int(w.Settings.GetChunkWidth())
	y := (index / int(w.Settings.GetChunkWidth())) % int(w.Settings.GetChunkWidth())
	z := index / (int(w.Settings.GetChunkWidth()) * int(w.Settings.GetChunkWidth()))
	return vec.Vec3I{x, y, z}
}

func (w *World) GetBiomeFromId(id int) *messages.Biome {
	for _, v := range w.Settings.Biomes.GetBiomes() {
		if int(v.GetId()) == id {
			return v
		}
	}
	return &messages.Biome{}
}

type Entity interface {
	GetPosition()
	GetId() int
}
