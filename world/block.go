package world

import (
	"encoding/json"
	"errors"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/vec"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	ErrPropertyNotFound = errors.New("Property not found")
)

type BlockFlag byte

const (
	BlockConnectedGround BlockFlag = 1 << iota // Wether this block is connected to the grond or not
	BlockOccupiedFull                          // If set, no units or anything can pass
	BlockOccupiedHalf                          // Only small units can pass
	BlockHidden                                // Wether this block is visible or not
)

type BlockProbability struct {
	Everywhere int
	Outside    int
	Inside     int
	Biomes     map[string]int
}

type BlockType struct {
	Id        int
	Name      string
	Flags     []string
	Biomes    []string
	AllBiomes bool
	Type      string
	Layer     string

	Probability BlockProbability

	LayerStart   int
	LayerEnd     int
	LayerOutSide bool

	// Additional properties
	Properties map[string]string
}

func BlockTypesFromJson(data []byte) ([]BlockType, ferr.FortiaError) {
	// Decode the json
	var btypes []BlockType
	err := json.Unmarshal(data, &btypes)
	if err != nil {
		return []BlockType{}, ferr.Wrap(err, "")
	}
	for i, v := range btypes {
		if v.LayerStart == 0 && v.LayerEnd == 0 {
			if v.Layer == "outside" {
				btypes[i].LayerOutSide = true
			} else if v.Layer == "inside" {
				btypes[i].LayerOutSide = false
			} else if v.Layer == "*" {
				btypes[i].LayerEnd = 1000
			} else if strings.Contains(v.Layer, "-") {
				split := strings.Split(v.Layer, "-")
				start, err := strconv.Atoi(split[0])
				if err != nil {
					return []BlockType{}, ferr.Wrap(err, "")
				}
				end, err := strconv.Atoi(split[1])
				if err != nil {
					return []BlockType{}, ferr.Wrap(err, "")
				}
				btypes[i].LayerStart = start
				btypes[i].LayerEnd = end
			}
		}

		if len(v.Biomes) == 0 || v.Biomes[0] == "*" {
			btypes[i].AllBiomes = true
		}
	}
	return btypes, nil
}

func BlockTypesFromFile(file string) ([]BlockType, ferr.FortiaError) {
	data, nErr := ioutil.ReadFile(file)
	if nErr != nil {
		return []BlockType{}, ferr.Wrap(nErr, "")
	}
	btypes, err := BlockTypesFromJson(data)
	return btypes, err
}

func (j *BlockType) GetPropertyInt(key string) (value int, err error) {
	value = -1
	err = nil

	strVal, ok := j.Properties[key]
	if !ok {
		err = ErrPropertyNotFound
		return
	}

	value, err = strconv.Atoi(strVal)
	return
}

type Block struct {
	LocalPosition vec.Vec2I              `json:"-"` // Position relative to layer
	Layer         *Layer                 `json:"-"`
	Kind          *BlockType             `json:"-"`
	Entities      []int                  `json:",omitempty"`
	Flags         BlockFlag              `json:",omitempty"`
	Data          map[string]interface{} `json:",omitempty"`

	Id int
}

// TODO check chunks nearby
// Should we still check even if this block is air?
func (b *Block) CheckHidden(neighbours []*Chunk) (bool, ferr.FortiaError) {

	if b.Layer == nil {
		return false, ferr.New("Layer nil")
	}
	if b.Layer.Chunk == nil {
		return false, ferr.New("Chunk nil")
	}
	if b.Id == 0 { // Air so this is visible
		return false, nil
	}

	chunk := b.Layer.Chunk
	pos := b.LocalPosition
	layerSize := b.Layer.World.GeneralInfo.LayerSize

	// Set chunk edges to not covered for now
	// if pos.X == 0 || pos.X >= b.Layer.World.GeneralInfo.LayerSize ||
	// 	pos.Y == 0 || pos.Y >= b.Layer.World.GeneralInfo.LayerSize {
	// 	return false, nil
	// }

	blocks := make([]*Block, 0)

	// get surounding blocks on same layer
	blocks = append(blocks, b.Layer.GetLocalBlock(pos.X+1, pos.Y))
	blocks = append(blocks, b.Layer.GetLocalBlock(pos.X-1, pos.Y))
	blocks = append(blocks, b.Layer.GetLocalBlock(pos.X, pos.Y+1))
	blocks = append(blocks, b.Layer.GetLocalBlock(pos.X, pos.Y-1))

	// Below and above
	if b.Layer.Position.Z > 0 {
		// Check block below
		layer := b.Layer.Chunk.Layers[b.Layer.Position.Z-1]
		blocks = append(blocks, layer.GetLocalBlock(pos.X, pos.Y))
	}

	if b.Layer.Position.Z < b.Layer.World.GeneralInfo.Height-1 {
		// Check block above
		layer := b.Layer.Chunk.Layers[b.Layer.Position.Z+1]
		blocks = append(blocks, layer.GetLocalBlock(pos.X, pos.Y))
	}

	// Check chunk neighbours
	if pos.X == 0 || pos.X >= layerSize ||
		pos.Y == 0 || pos.Y >= layerSize {
		// map out the chunks to make it easier
		var x1c *Chunk  // Chunk x + 1
		var x_1c *Chunk // x - 1
		var y1c *Chunk  // y + 1
		var y_1c *Chunk // y -1 All relative

		for _, v := range neighbours {
			diff := v.Position.Clone()
			diff.Sub(chunk.Position)

			if diff.X == 1 && diff.Y == 0 {
				x1c = v
			} else if diff.X == -1 && diff.Y == 0 {
				x_1c = v
			} else if diff.X == 0 && diff.Y == 1 {
				y1c = v
			} else if diff.X == 0 && diff.Y == -1 {
				y_1c = v
			}
		}

		if pos.X == 0 {
			if x1c != nil {
				cl := x_1c.Layers[b.Layer.Position.Z]
				blocks = append(blocks, cl.GetLocalBlock(layerSize-1, pos.Y))
			}
		} else if pos.X >= layerSize {
			if x_1c != nil {
				cl := x1c.Layers[b.Layer.Position.Z]
				blocks = append(blocks, cl.GetLocalBlock(0, pos.Y))
			}
		}

		if pos.Y == 0 {
			if y1c != nil {
				cl := y_1c.Layers[b.Layer.Position.Z]
				blocks = append(blocks, cl.GetLocalBlock(pos.X, layerSize-1))
			}
		} else if pos.Y >= layerSize {
			if y_1c != nil {
				cl := y1c.Layers[b.Layer.Position.Z]
				blocks = append(blocks, cl.GetLocalBlock(pos.X, 0))
			}
		}

	}

	for _, v := range blocks {
		if v == nil {
			continue
		}
		if v.Id <= 0 {
			// air
			return false, nil
		}
	}

	return true, nil
}
