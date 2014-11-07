package world

import (
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/vec"
)

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
				c.Layers = make([]*Layer, c.World.GeneralInfo.Height)
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

	fErr := w.Db.SetChunkInfo(chunk)

	return fErr
}

// Fetches a chunk from the database at x,y
// Chunk is nil if not found
func (w *World) GetChunk(x, y int, getLayers bool) (*Chunk, ferr.FortiaError) {
	chunk, err := w.Db.GetChunkInfo(vec.Vec2I{x, y})
	if err != nil {
		return nil, err
	}
	chunk.World = w

	if getLayers {
		positions := make([]vec.Vec3I, w.GeneralInfo.Height)
		for i := 0; i < w.GeneralInfo.Height; i++ {
			positions[i] = vec.Vec3I{x, y, i}
		}
		layers, err := w.Db.GetLayers(positions)
		if err != nil {
			return nil, err
		}

		for k, l := range layers {

			l.World = w
			l.Chunk = chunk

			layers[k] = l
		}
		chunk.Layers = layers

	}

	return chunk, nil
}
