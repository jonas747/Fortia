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
	LayersStr string

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
			if v.LayersStr == "outside" {
				btypes[i].LayerOutSide = true
			} else if v.LayersStr == "inside" {
				btypes[i].LayerOutSide = false
			} else if v.LayersStr == "*" {
				btypes[i].LayerEnd = 1000
			} else if strings.Contains(v.LayersStr, "-") {
				split := strings.Split(v.LayersStr, "-")
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

		if v.Biomes[0] == "*" {
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

type Chunk struct {
	Position vec.Vec2I
	Layers   []*Layer
	Biome    Biome
}

func (w *World) SetChunk(chunk *Chunk) {

}
