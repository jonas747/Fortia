package world

import (
	"encoding/json"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/vec"
)

type Layer struct {
	World    *World `json:"-"`
	Chunk    *Chunk `json:"-"`
	Position vec.Vec3I
	Blocks   []*Block
	Flags    int
	IsAir    bool // True if this layer is just air
}

// Saves a layer to the database
func (w *World) SetLayer(layer *Layer) ferr.FortiaError {
	return w.Db.SetLayer(layer)
}

// Returns a layer from the database
func (w *World) GetLayer(pos vec.Vec3I) (*Layer, ferr.FortiaError) {
	layer, err := w.Db.GetLayer(pos)
	if err != nil {
		return nil, err
	}
	layer.World = w
	return layer, nil
}

// Gets the block at local position lx ly, return nil if out of bounds
func (l *Layer) GetLocalBlock(lx, ly int) *Block {
	index := l.World.CoordsToIndex(vec.Vec3I{lx, ly, 0})
	if index >= len(l.Blocks) || index < 0 {
		return nil
	}
	return l.Blocks[index]
}

// Returns serialized json of the layer
func (l *Layer) Json() ([]byte, ferr.FortiaError) {
	out, err := json.Marshal(l)
	if err != nil {
		return []byte{}, ferr.Wrap(err, "")
	}
	return out, nil
}
