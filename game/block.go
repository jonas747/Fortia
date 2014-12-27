package game

import (
	"github.com/jonas747/fortia/messages"
	"github.com/jonas747/fortia/vec"
)

type BlockFlag int32

const (
	BlockConnectedGround BlockFlag = 1 << iota // Wether this block is connected to the grond or not
	BlockOccupiedFull                          // If set, no units or anything can pass
	BlockOccupiedHalf                          // Only small units can pass
	BlockHidden                                // Wether this block is visible or not
)

type Block struct {
	Chunk         *Chunk
	LocalPosition vec.Vec3I
	RawBlock      *messages.Block
}

// Checks surouding blocks and returns wether this block is hidden or not
func (b *Block) CheckHidden(neighbours map[vec.Vec2I]*Chunk) bool {

	if b.Chunk == nil {
		return false
	}
	if *b.RawBlock.Kind <= 0 { // Air, obviously visible
		return false
	}

	//chunk := b.Layer.Chunk
	pos := b.LocalPosition
	chunkWidth := int(b.Chunk.World.Settings.GetChunkWidth())
	chunkHeight := int(b.Chunk.World.Settings.GetChunkHeight())

	blocks := make([]*Block, 0)

	// get surounding blocks on same layer
	if pos.X > 0 {
		blocks = append(blocks, b.Chunk.GetBlock(vec.Vec3I{pos.X - 1, pos.Y, pos.Z}))
	}
	if pos.X < chunkWidth-1 {
		blocks = append(blocks, b.Chunk.GetBlock(vec.Vec3I{pos.X + 1, pos.Y, pos.Z}))
	}

	if pos.Y > 0 {
		blocks = append(blocks, b.Chunk.GetBlock(vec.Vec3I{pos.X, pos.Y - 1, pos.Z}))
	}
	if pos.Y < chunkWidth-1 {
		blocks = append(blocks, b.Chunk.GetBlock(vec.Vec3I{pos.X, pos.Y + 1, pos.Z}))
	}

	if pos.Z > 0 {
		blocks = append(blocks, b.Chunk.GetBlock(vec.Vec3I{pos.X, pos.Y, pos.Z - 1}))
	}
	if pos.Z < chunkHeight-1 {
		blocks = append(blocks, b.Chunk.GetBlock(vec.Vec3I{pos.X, pos.Y, pos.Z + 1}))
	}

	// Check chunk neighbours
	if pos.X == 0 || pos.X >= chunkWidth-1 ||
		pos.Y == 0 || pos.Y >= chunkWidth-1 {
		// map out the chunks to make it easier
		x1c := neighbours[vec.Vec2I{1, 0}]
		x_1c := neighbours[vec.Vec2I{-1, 0}]
		y1c := neighbours[vec.Vec2I{0, 1}]
		y_1c := neighbours[vec.Vec2I{0, -1}]

		if pos.X == 0 {
			if x_1c != nil {
				blocks = append(blocks, x_1c.GetBlock(vec.Vec3I{chunkWidth - 1, pos.Y, pos.Z}))
			}
		} else if pos.X >= chunkWidth-1 {
			if x1c != nil {
				blocks = append(blocks, x1c.GetBlock(vec.Vec3I{0, pos.Y, pos.Z}))
			}
		}

		if pos.Y == 0 {
			if y_1c != nil {
				blocks = append(blocks, y_1c.GetBlock(vec.Vec3I{pos.X, chunkWidth - 1, pos.Z}))
			}
		} else if pos.Y >= chunkWidth-1 {
			if y1c != nil {
				blocks = append(blocks, y1c.GetBlock(vec.Vec3I{pos.X, 0, pos.Z}))
			}
		}

	}

	for _, v := range blocks {
		if v == nil {
			continue
		}
		if *v.RawBlock.Kind <= 0 {
			// air
			return false
		}
	}

	return true
}
