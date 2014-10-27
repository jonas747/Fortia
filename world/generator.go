package world

import (
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/vec"
	"math/rand"
)

/*
	The fortia world generator

	Steps:
	- Assign biome to chunk
	- Generate landscape based on biomes roughness and the biomes flags (ocean will sink everything down a bit) property
	- smooth chunk edges
	- metals, soil, trees, vegetation etc..

	- Wildlife gets spawned by the ticker(also dependign on the wildlife property of this biome)
*/
type Generator struct {
	W               *World // The world
	Rng             *rand.Rand
	Biomes          []Biome             // The biomes this generator uses
	BlockTypes      map[int]BlockType   // The blocktypes this generator uses, accessed by id
	ChunkBlockTypes map[int][]BlockType // The blocktypes but accessed by chunkid instead
}

// Initialises the generator
func (g *Generator) Init() {

}

// Generates a chunk, saves chunk information and layers to db.
func (g *Generator) GenerateChunk(position vec.Vec2I) ferr.FortiaError {
	return nil
}

// First stage
func (g *Generator) getBiome(position vec.Vec2I) Biome {
	return Biome{}
}

// Second stage: Returns a chunk, at this stage it only sets the block id to one of 2, 0 for air and 1 for land
func (g *Generator) generateLandscape(position vec.Vec2I, biome Biome) *Chunk {
	return &Chunk{}
}

// Smooths the chunk edges
func (g *Generator) smoothEedges(chunk *Chunk) *Chunk {
	return &Chunk{}
}

// Assigns proper blocks to everything, stone should be stone etc...
func (g *Generator) assignBlocks(chunk *Chunk) *Chunk {
	return &Chunk{}
}
