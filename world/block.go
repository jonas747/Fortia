package world

import (
	"errors"
	"github.com/jonas747/fortia/vec"
	"strconv"
)

var (
	ErrPropertyNotFound = errors.New("Property not found")
)

type BlockType struct {
	Id         int
	Name       string
	Flags      []string
	Properties map[string]string
	Biomes     []string
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

func (w *World) SetChunk(chunk *Chunk){
	
}