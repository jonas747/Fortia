package world

import (
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/vec"
)

type BlockFlag byte

const (
	BlockConnectedGround BlockFlag = 1 << iota // Wether this block is connected to the grond or not
	BlockOccupiedFull                          // If set, no units or anything can pass
	BlockOccupiedHalf                          // Only small units can pass
	BlockHidden                                // Wether this block is visible or not
)

type Block struct {
	LocalPosition vec.Vec2I              `json:"-"` // Position relative to layer
	Layer         *Layer                 `json:"-"`
	Kind          *BlockType             `json:"-"`
	Entities      []int                  `json:",omitempty"`
	Flags         BlockFlag              `json:",omitempty"`
	Data          map[string]interface{} `json:",omitempty"`

	Id int
}

// TODO check chunks nearby
// Should we still check even if this block is air?
func (b *Block) CheckHidden(neighbours map[vec.Vec2I]*Chunk) (bool, ferr.FortiaError) {

	if b.Layer == nil {
		return false, ferr.New("Layer nil")
	}
	if b.Layer.Chunk == nil {
		return false, ferr.New("Chunk nil")
	}
	if b.Id <= 0 { // Air so this is visible
		return false, nil
	}

	//chunk := b.Layer.Chunk
	pos := b.LocalPosition
	layerSize := b.Layer.World.GeneralInfo.LayerSize

	blocks := make([]*Block, 0)

	// get surounding blocks on same layer
	if pos.X > 0 {
		blocks = append(blocks, b.Layer.GetLocalBlock(pos.X-1, pos.Y))
	}
	if pos.X < layerSize-1 {
		blocks = append(blocks, b.Layer.GetLocalBlock(pos.X+1, pos.Y))
	}

	if pos.Y > 0 {
		blocks = append(blocks, b.Layer.GetLocalBlock(pos.X, pos.Y-1))
	}
	if pos.Y < layerSize-1 {
		blocks = append(blocks, b.Layer.GetLocalBlock(pos.X, pos.Y+1))
	}

	// Below and above
	if b.Layer.Position.Z > 0 {
		// Check block below
		layer := b.Layer.Chunk.Layers[b.Layer.Position.Z-1]
		blocks = append(blocks, layer.GetLocalBlock(pos.X, pos.Y))
	}

	if b.Layer.Position.Z < b.Layer.World.GeneralInfo.Height-1 {
		// Check block above
		layer := b.Layer.Chunk.Layers[b.Layer.Position.Z+1]
		blocks = append(blocks, layer.GetLocalBlock(pos.X, pos.Y))
		bl := layer.GetLocalBlock(pos.X, pos.Y)
		if bl == nil {
			b.Layer.World.Logger.Info("shouldnt be nil")
		}
	}

	// Check chunk neighbours
	if pos.X == 0 || pos.X >= layerSize-1 ||
		pos.Y == 0 || pos.Y >= layerSize-1 {
		// map out the chunks to make it easier
		x1c := neighbours[vec.Vec2I{1, 0}]
		x_1c := neighbours[vec.Vec2I{-1, 0}]
		y1c := neighbours[vec.Vec2I{0, 1}]
		y_1c := neighbours[vec.Vec2I{0, -1}]

		if pos.X == 0 {
			if x_1c != nil {
				cl := x_1c.Layers[b.Layer.Position.Z]
				blocks = append(blocks, cl.GetLocalBlock(layerSize-1, pos.Y))
			}
		} else if pos.X >= layerSize-1 {
			if x1c != nil {
				cl := x1c.Layers[b.Layer.Position.Z]
				blocks = append(blocks, cl.GetLocalBlock(0, pos.Y))
			}
		}

		if pos.Y == 0 {
			if y_1c != nil {
				cl := y_1c.Layers[b.Layer.Position.Z]
				blocks = append(blocks, cl.GetLocalBlock(pos.X, layerSize-1))
			}
		} else if pos.Y >= layerSize-1 {
			if y1c != nil {
				cl := y1c.Layers[b.Layer.Position.Z]
				blocks = append(blocks, cl.GetLocalBlock(pos.X, 0))
			}
		}

	}

	for _, v := range blocks {
		if v == nil {
			b.Layer.World.Logger.Debug("pop")
			continue
		}
		if v.Id <= 0 {
			// air
			return false, nil
		}
	}

	return true, nil
}
