// Package world contains world realted stuff

package world

import (
	"github.com/jonas747/fortia/db"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/vec"
	"strconv"
)

type World struct {
	Logger    *log.LogClient
	Db        *db.GameDB
	ChunkSize int
}

func (w *World) ChunkToWorldPos(chunkPos vec.Vec3I) vec.Vec3I {
	chunkWorldPos := chunkPos.Clone()
	chunkWorldPos.MltiplyScalar(float64(w.ChunkSize))
	return chunkWorldPos
}

func (w *World) GenChunk(pos vec.Vec3I) *Chunk {
	//chunkWorldPos := w.ChunkToWorldPos(pos)

	blocks := make([]*Block, w.ChunkSize*w.ChunkSize*w.ChunkSize)

	for x := 0; x < w.ChunkSize; x++ {
		for y := 0; y < w.ChunkSize; y++ {
			for z := 0; z < w.ChunkSize; z++ {
				// worldX := x + chunkWorldPos.X
				// worldY := y + chunkWorldPos.Y
				// TODO Actual generation
				b := Block{
					LocalPosition: vec.Vec3I{X: x, Y: y, Z: z},
					ChunkPos:      pos.Clone(),
					Id:            0,
				}
				if z < 50 {
					b.Id = 1
				}
				blocks[w.CoordsToIndex(vec.Vec3I{X: x, Y: y, Z: z})] = &b
			}
		}
	}
	chunk := &Chunk{
		Blocks:        blocks,
		ChunkPosition: pos,
		World:         w,
	}
	return chunk
}

func (w *World) CoordsToIndex(pos vec.Vec3I) int {
	//Flat[x + WIDTH * (y + DEPTH * z)] = Original[x, y, z]
	return pos.X + w.ChunkSize*(pos.Y+w.ChunkSize*pos.Z)
}

func (w *World) IndexToCoords(index int) vec.Vec3I {
	/*
		z = Math.round(i / (WIDTH * HEIGHT));
		y = Math.round((i - z * WIDTH * HEIGHT) / WIDTH);
		x = i - WIDTH * (y + HEIGHT * z);
	*/
	z := index / (w.ChunkSize * w.ChunkSize)
	y := (index - z*w.ChunkSize*w.ChunkSize) / w.ChunkSize
	x := index - w.ChunkSize*(y+w.ChunkSize*z)
	return vec.Vec3I{
		X: x,
		Y: y,
		Z: z,
	}
}

func (w *World) GetChunk(pos vec.Vec3I) *Chunk {
	//TODO
	return nil
}

type Chunk struct {
	Blocks        []*Block
	ChunkPosition vec.Vec3I
	World         *World
}

func (c *Chunk) LocalToWorld(pos vec.Vec3I) vec.Vec3I {
	chunkWorldPos := c.World.ChunkToWorldPos(c.ChunkPosition)

	worldPos := pos.Clone()
	worldPos.Add(chunkWorldPos)
	return worldPos
}

func (c *Chunk) ExportRedisList() []string {
	numBlocks := c.World.ChunkSize * c.World.ChunkSize * c.World.ChunkSize
	out := make([]string, numBlocks)
	for i := 0; i < numBlocks; i++ {
		block := c.Blocks[i]
		out[i] = block.ExportRedisString()
	}
	return out
}

type Block struct {
	LocalPosition vec.Vec3I
	ChunkPos      vec.Vec3I
	Entities      []Entity
	Id            int
}

func (b *Block) ExportRedisString() string {
	str := strconv.Itoa(b.Id)
	if len(b.Entities) > 0 {
		for _, v := range b.Entities {
			str += ":"
			str += strconv.Itoa(v.GetId())
		}
	}
	return str
}

type Entity interface {
	GetPosition()
	GetId() int
}
