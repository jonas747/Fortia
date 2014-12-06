package world

import (
	"code.google.com/p/goprotobuf/proto"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/messages"
	"github.com/jonas747/fortia/vec"
)

// Represents a chunk
type Chunk struct {
	World    *World `json:"-"`
	RawChunk *messages.Chunk
}

// returns chunk at x y, local to current chunk
func (c *Chunk) GetNeighbour(x, y int) (*Chunk, ferr.FortiaError) {
	wPos := vec.Vec2I{int(c.RawChunk.GetX()), int(c.RawChunk.GetY())}
	wPos.Add(vec.Vec2I{x, y})
	chunk, err := c.World.GetChunk(wPos)
	return chunk, err
}

// Returns all neighbours
func (c *Chunk) GetAllNeighbours() ([]*Chunk, ferr.FortiaError) {
	out := make([]*Chunk, 0)
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			if x == 0 && y == 0 {
				// Dont add itself
				continue
			}
			chunk, err := c.GetNeighbour(x, y)
			if err != nil {
				if err.GetCode() == 404 { // continue even if the chunks was not found in the db
					continue
				}
				return out, err
			}
			out = append(out, chunk)
		}
	}
	return out, nil
}

func (c *Chunk) GetBlock(pos vec.Vec3I) *Block {
	index := c.World.CoordsToIndex(pos)
	if index > len(c.RawChunk.Blocks) {
		return nil
	}
	rawBlock := c.RawChunk.Blocks[index]
	b := &Block{
		RawBlock:      rawBlock,
		Chunk:         c,
		LocalPosition: pos,
	}
	return b
}

// Flags all hidden blocks as hidden
// If provided neighbours' len is 0 then it will fetch from db
func (c *Chunk) FlagHidden(neighbours map[vec.Vec2I]*Chunk) {
	if len(neighbours) < 1 {
		n := make(map[vec.Vec2I]*Chunk)
		errs := make([]ferr.FortiaError, 0)
		c1, e1 := c.GetNeighbour(1, 0)
		c2, e2 := c.GetNeighbour(0, 1)
		c3, e3 := c.GetNeighbour(-1, 0)
		c4, e4 := c.GetNeighbour(0, -1)
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
				if v.GetCode() != 404 {
					c.World.Logger.Error(v)
				}
			}
		}
		neighbours = n
	}
	visibleLayers := make([]bool, c.World.GeneralInfo.ChunkHeight)
	for k, v := range c.RawChunk.Blocks {
		b := &Block{
			RawBlock:      v,
			Chunk:         c,
			LocalPosition: c.World.IndexToCoords(k),
		}
		hidden := b.CheckHidden(neighbours)
		flags := v.GetFlags()
		if hidden {
			flags |= int32(BlockHidden)
		} else {
			if flags&int32(BlockHidden) != 0 {
				flags ^= int32(BlockHidden)
			}
			visibleLayers[b.LocalPosition.Z] = true
		}
		v.Flags = proto.Int32(flags)
	}
	c.RawChunk.VisibleLayers = visibleLayers
}

func (w *World) GetChunk(pos vec.Vec2I) (*Chunk, ferr.FortiaError) {
	raw, err := w.Db.GetChunk(pos)
	if err != nil {
		return nil, err
	}

	chunk := &Chunk{
		World:    w,
		RawChunk: raw,
	}

	return chunk, nil
}

func (w *World) SetChunk(chunk *Chunk) ferr.FortiaError {
	return w.Db.SetChunk(chunk.RawChunk)
}
