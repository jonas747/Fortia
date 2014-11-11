package world

import (
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/vec"
	"sync"
)

// Represents a chunk
type Chunk struct {
	World         *World   `json:"-"`
	Layers        []*Layer `json:"-"` // No need to store layers twice...
	Position      vec.Vec2I
	Biome         Biome
	Potency       int   // The biome potency this chunk has
	VisibleLayers []int // Layers that have one or more visible blocks
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
			if x == 0 && y == 0 {
				// Dont add istelf
				continue
			}
			chunk, err := c.GetNeighbour(x, y)
			if err != nil {
				if err.GetMessage() == "404" { // continue even if the chunks was not found in the db
					continue
				}
				return out, err
			}
			out = append(out, chunk)
		}
	}
	return out, nil
}

// Flags all hidden blocks as hidden
// If provided neighbours' len is 0 then it will fetch from db
func (c *Chunk) FlagHidden(neighbours []*Chunk) {
	if len(neighbours) < 1 {
		n, err := c.GetAllNeighbours()
		if err != nil {
			c.World.Logger.Error(err)
		}
		neighbours = n
	}
	var wg sync.WaitGroup
	wg.Add(c.World.GeneralInfo.Height)
	visibleLayers := make([]bool, c.World.GeneralInfo.Height)
	for k, layer := range c.Layers {
		l := layer
		n := k
		go func() {
			l.Hidden = true
			for _, block := range l.Blocks {
				if block.Layer == nil {
					block.Layer = l
				}
				hidden, err := block.CheckHidden(neighbours)
				if err != nil {
					c.World.Logger.Error(err)
				}
				if hidden {
					block.Flags |= BlockHidden
				} else {
					if block.Flags&BlockHidden != 0 {
						block.Flags ^= BlockHidden
					}
					l.Hidden = false
				}
			}
			visibleLayers[n] = !l.Hidden
			wg.Done()
		}()
	}
	wg.Wait()
	c.VisibleLayers = make([]int, 0)
	for k, v := range visibleLayers {
		if v == true {
			c.VisibleLayers = append(c.VisibleLayers, k)
		}
	}
}

// Saves the chunk to the database
func (w *World) SetChunk(chunk *Chunk, setLayers bool) ferr.FortiaError {
	if setLayers {
		err := w.Db.SetLayers(chunk.Layers)
		if err != nil {
			return err
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
		realLayers := make([]*Layer, 200)
		for _, l := range layers {

			l.World = w
			l.Chunk = chunk

			realLayers[l.Position.Z] = l
		}
		chunk.Layers = realLayers

	}

	return chunk, nil
}
