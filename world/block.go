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

const (
	BFlagCovered = 1 << iota // On if the block is souronded by blocks and cannot be seen
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
	LocalPosition vec.Vec2I              `json:"-"`
	Layer         *Layer                 `json:"-"`
	Kind          *BlockType             `json:"-"`
	Entities      []int                  `json:",omitempty"`
	Flags         byte                   `json:",omitempty"`
	Data          map[string]interface{} `json:",omitempty"`

	Id int
}

// TODO check chunks nearby
// Should we still check even if this block is air?
func (b *Block) IsSurounded() (bool, ferr.FortiaError) {
	if b.Layer == nil {
		return false, ferr.New("Layer nil")
	}
	if b.Layer.Chunk == nil {
		return false, ferr.New("Chunk nil")
	}

	pos := b.LocalPosition

	// Set chunk edges to not covered for now
	if pos.X == 0 || pos.X >= b.Layer.World.LayerSize ||
		pos.Y == 0 || pos.Y >= b.Layer.World.LayerSize {
		return false, nil
	}

	// get surounding blocks
	blocks := make([]*Block, 0)
	blocks = append(blocks, b.Layer.GetLocalBlock(pos.X+1, pos.Y))
	blocks = append(blocks, b.Layer.GetLocalBlock(pos.X-1, pos.Y))
	blocks = append(blocks, b.Layer.GetLocalBlock(pos.X, pos.Y+1))
	blocks = append(blocks, b.Layer.GetLocalBlock(pos.X, pos.Y-1))

	if b.Layer.Position.Z > 0 {
		// Check block below
		layer := b.Layer.Chunk.Layers[b.Layer.Position.Z-1]
		blocks = append(blocks, layer.GetLocalBlock(pos.X, pos.Y))
	}

	if b.Layer.Position.Z < b.Layer.World.WorldHeight-1 {
		// Check block above
		layer := b.Layer.Chunk.Layers[b.Layer.Position.Z+1]
		blocks = append(blocks, layer.GetLocalBlock(pos.X, pos.Y))
	}

	for _, v := range blocks {
		if v == nil || v.Id <= 0 {
			// air
			return false, nil
		}
	}

	return true, nil
}
