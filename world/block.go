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
	World    *World `json:"-"`
	Position vec.Vec3I
	Blocks   []*Block
	Flags    int
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
	World    *World   `json:"-"`
	Layers   []*Layer `json:"-"` // No need to store layers twice...
	Position vec.Vec2I
	Biome    Biome
	Potency  int // The biome potency this chunk has
}

// Rerturns a layer of the chunk, if it is in the chunk's cache then it will return that
// If fetch is true then it will fetch even if it is in the cache
// If store is true it will store the layer int he chunk's cache after fetching it
func (c *Chunk) GetLayer(layer int, fetch, store bool) (*Layer, ferr.FortiaError) {
	if len(c.Layers) == 0 || fetch {
		l, err := c.World.GetLayer(vec.Vec3I{c.Position.X, c.Position.Y, layer})
		if err != nil {
			return nil, err
		}
		if store {
			if len(c.Layers) == 0 {
				c.Layers = make([]*Layer, c.World.WorldHeight)
			}
			c.Layers[layer] = l
		}
		return l, nil
	}

	l := c.Layers[layer]
	return l, nil
}

// returns chunk at x y, local to current chunk
func (c *Chunk) GetNeighbour(x, y int) (*Chunk, ferr.FortiaError) {
	wPos := c.Position.Clone()
	wPos.Add(vec.Vec2I{x, y})

	chunk, err := c.World.GetChunk(wPos.X, wPos.Y, false)
	return chunk, err
}

// Returns all neighbours
func (c *Chunk) GetAllNeighbours() ([]*Chunk, ferr.FortiaError) {
	out := make([]*Chunk, 0)
	for x := -1; x < 1; x++ {
		for y := -1; y < 1; y++ {
			chunk, err := c.GetNeighbour(x, y)
			if err != nil {
				return out, err
			}
			out = append(out, chunk)
		}
	}
	return out, nil
}

// Saves the chunk to the database
func (w *World) SetChunk(chunk *Chunk, setLayers bool) ferr.FortiaError {
	if setLayers {
		for _, l := range chunk.Layers {
			err := w.SetLayer(l)
			if err != nil {
				return err
			}
		}
	}

	serialised, err := json.Marshal(chunk)
	if err != nil {
		return ferr.Wrap(err, "")
	}

	fErr := w.Db.SetChunkInfo(chunk.Position.X, chunk.Position.Y, serialised)

	return fErr
}

// Fetches a chunk from the database at x,y
func (w *World) GetChunk(x, y int, getLayers bool) (*Chunk, ferr.FortiaError) {
	if getLayers {
		// TODO (see general tasks)
	}

	raw, err := w.Db.GetChunkInfo(x, y)
	if err != nil {
		return nil, err
	}

	var cinfo Chunk
	nErr := json.Unmarshal(raw, &cinfo)
	if nErr != nil {
		return nil, ferr.Wrap(nErr, "")
	}

	return &cinfo, nil
}
