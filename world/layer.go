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
	raw, err := layer.Json()
	if err != nil {
		return err
	}
	return w.Db.SetLayer(layer.Position.X, layer.Position.Y, layer.Position.Z, raw)
}

// Returns a layer from the database
func (w *World) GetLayer(pos vec.Vec3I) (*Layer, ferr.FortiaError) {
	rawLayer, err := w.Db.GetLayer(pos.X, pos.Y, pos.Z)
	if err != nil {
		return nil, err
	}

	layer := &Layer{
		Position: pos,
		World:    w,
	}

	nErr := json.Unmarshal(rawLayer, layer)
	if nErr != nil {
		return nil, ferr.Wrap(nErr, "")
	}

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

// Represents a chunk
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

// Flags all surounded blocks as surounded
// TODO remove flag if not covered and flagged allready
func (c *Chunk) FlagSurounded() {
	for _, layer := range c.Layers {
		for _, block := range layer.Blocks {
			surounded, _ := block.IsSurounded()
			if surounded {
				block.Flags |= BFlagCovered
			}
		}
	}
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
// Chunk is nil if not found
func (w *World) GetChunk(x, y int, getLayers bool) (*Chunk, ferr.FortiaError) {
	raw, found, err := w.Db.GetChunkInfo(x, y)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}

	chunk := &Chunk{
		World: w,
	}
	nErr := json.Unmarshal(raw, chunk)
	if nErr != nil {
		return nil, ferr.Wrap(nErr, "")
	}

	if getLayers {
		positions := make([]*vec.Vec3I, w.WorldHeight)
		for i := 0; i < w.WorldHeight; i++ {
			positions[i] = &vec.Vec3I{x, y, i}
		}
		layersRaw, err := w.Db.GetLayers(positions)
		if err != nil {
			return nil, err
		}

		layers := make([]*Layer, w.WorldHeight)
		for _, lRaw := range layersRaw {
			layer := &Layer{
				World: w,
				Chunk: chunk,
			}

			nErr := json.Unmarshal(lRaw, &layer)
			if nErr != nil {
				return nil, ferr.Wrap(nErr, "")
			}
			layers[layer.Position.Z] = layer
		}
		chunk.Layers = layers

	}

	return chunk, nil
}
