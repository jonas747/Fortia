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
	"strconv"
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
	WorldHeight int
	WorldGen    WorldGenerator

	BlockTypes map[int]*BlockType
	Biomes     []*Biome
}

func (w *World) LoadSettingsFromDb() ferr.FortiaError {
	settings, err := w.Db.GetWorldInfo()
	if err != nil {
		return err
	}
	w.WorldHeight, _ = strconv.Atoi(settings["worldHeight"])
	w.LayerSize, _ = strconv.Atoi(settings["layerSize"])
	btypesJson := settings["blockTypes"]
	biomesJson := settings["biomes"]

	// Decode the json
	var blocktypes map[int]*BlockType
	nErr := json.Unmarshal([]byte(btypesJson), &blocktypes)
	if nErr != nil {
		return ferr.Wrap(nErr, "")
	}
	w.BlockTypes = blocktypes

	var biomes []*Biome
	nErr = json.Unmarshal([]byte(biomesJson), &biomes)
	if nErr != nil {
		return ferr.Wrap(err, "")
	}
	return nil
}

func (w *World) SaveSettingsToDb() ferr.FortiaError {
	btypesJson, err := json.Marshal(w.BlockTypes)
	if err != nil {
		return ferr.Wrap(err, "")
	}
	biomesJson, err := json.Marshal(w.Biomes)
	if err != nil {
		return ferr.Wrap(err, "")
	}

	infoMap := map[string]interface{}{
		"layerSize":   w.LayerSize,
		"worldHeight": w.WorldHeight,
		"blockTypes":  btypesJson,
		"biomes":      biomesJson,
	}
	return w.Db.SetWorldInfo(infoMap)
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
			if layer.Position.Z > w.WorldHeight/2 {
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

type Entity interface {
	GetPosition()
	GetId() int
}

type WorldGenerator interface {
	GenerateBLock(x, y, z int) int
}
