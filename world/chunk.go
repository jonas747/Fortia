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
func (c *Chunk) GetNeighbour(x, y int, fetchLayers bool) (*Chunk, ferr.FortiaError) {
	wPos := c.Position.Clone()
	wPos.Add(vec.Vec2I{x, y})
	//c.World.Logger.Info("cpos ", c.Position, "n ", x, ",", y, " wpos ", wPos)
	chunk, err := c.World.GetChunk(wPos.X, wPos.Y, fetchLayers, false)
	return chunk, err
}

// Returns all neighbours
func (c *Chunk) GetAllNeighbours(fetchLayers bool) ([]*Chunk, ferr.FortiaError) {
	out := make([]*Chunk, 0)
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			if x == 0 && y == 0 {
				// Dont add istelf
				continue
			}
			chunk, err := c.GetNeighbour(x, y, fetchLayers)
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
func (c *Chunk) FlagHidden(neighbours map[vec.Vec2I]*Chunk) {
	if len(neighbours) < 1 {
		n := make(map[vec.Vec2I]*Chunk)
		errs := make([]ferr.FortiaError, 0)
		c1, e1 := c.GetNeighbour(1, 0, true)
		c2, e2 := c.GetNeighbour(0, 1, true)
		c3, e3 := c.GetNeighbour(-1, 0, true)
		c4, e4 := c.GetNeighbour(0, -1, true)
		if e1 == nil {
			n[vec.Vec2I{1, 0}] = c1
		}
		if e2 == nil {
			n[vec.Vec2I{0, 1}] = c2
		}
		if e3 == nil {
			n[vec.Vec2I{-1, 0}] = c3
		}
		if e4 == nil {
			n[vec.Vec2I{0, -1}] = c4
		}
		errs = append(errs, e1, e2, e3, e4)
		for _, v := range errs {
			if v != nil {
				if v.GetMessage() != "404" {
					c.World.Logger.Error(v)
				}
			}
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
			for index, block := range l.Blocks {
				coords := c.World.IndexToCoords(index)
				if block.Layer == nil {
					block.Layer = l
				}
				block.LocalPosition = vec.Vec2I{coords.X, coords.Y}
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
// If getlayers is true it will also fetch layers depending on what onlyVisible is set to
func (w *World) GetChunk(x, y int, getLayers, onlyVisible bool) (*Chunk, ferr.FortiaError) {
	if x > w.GeneralInfo.Size-1 || y > w.GeneralInfo.Size-1 ||
		x < 0 || y < 0 {
		return nil, ferr.New("404")
	}
	chunk, err := w.Db.GetChunkInfo(vec.Vec2I{x, y})
	if err != nil {
		return nil, err
	}
	chunk.World = w

	if getLayers {
		var positions []vec.Vec3I
		if onlyVisible {
			positions = make([]vec.Vec3I, len(chunk.VisibleLayers))
			for k, v := range chunk.VisibleLayers {
				positions[k] = vec.Vec3I{x, y, v}
			}
		} else {
			positions = make([]vec.Vec3I, w.GeneralInfo.Height)
			for i := 0; i < w.GeneralInfo.Height; i++ {
				positions[i] = vec.Vec3I{x, y, i}
			}
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
