package world

import (
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/simplex"
	"github.com/jonas747/fortia/vec"
	//"math/rand"
	"strconv"
)

/*
	The fortia world generator

	Steps:
	- Assign biome to chunk
	- Generate landscape based on biomes by roughness and the biomes flags (ocean will sink everything down a bit) property
	- smooth chunk edges
	- metals, soil, trees, vegetation etc..

	- Wildlife gets spawned by the ticker(also dependign on the wildlife property of this biome)
*/
type Generator struct {
	W               *World              // The world
	Biomes          *BiomesInfo         // The biomes this generator uses
	BlockTypes      []BlockType         // The blocktypes this generator uses, accessed by id
	BiomeBlockType  map[int][]BlockType // The blocktypes but accessed by biome id instead
	NoiseGenerators map[string]*simplex.Noise
}

func NewGenerator(world *World, biomes *BiomesInfo, blockTypes []BlockType, seeds map[string]int64) *Generator {
	// map the blocktypes to biomes for ease of access
	biomeBTypes := make(map[int][]BlockType)
	for _, v := range blockTypes {
		if v.AllBiomes {
			for _, b := range biomes.Biomes {
				cur, ok := biomeBTypes[b.Id]
				if !ok {
					cur = make([]BlockType, 0)
				}
				cur = append(cur, v)
				biomeBTypes[b.Id] = cur
			}
		} else {
			for _, bid := range v.Biomes {
				id, _ := strconv.Atoi(bid)
				cur, ok := biomeBTypes[id]
				if !ok {
					cur = make([]BlockType, 0)
				}
				cur = append(cur, v)
				biomeBTypes[id] = cur
			}
		}
	}
	// Init the noise generators
	ngens := make(map[string]*simplex.Noise)
	for k, v := range seeds {
		noiseGen := simplex.NewNoise(v)
		ngens[k] = noiseGen
	}

	generator := &Generator{
		W:               world,
		Biomes:          biomes,
		BlockTypes:      blockTypes,
		BiomeBlockType:  biomeBTypes,
		NoiseGenerators: ngens,
	}

	return generator
}

// Generates a chunk, saves chunk information and layers to db.
func (g *Generator) GenerateChunk(position vec.Vec2I) ferr.FortiaError {
	biome := g.getBiome(position)

	chunk := g.generateLandscape(position, biome)

	chunk, err := g.smoothEedges(chunk)
	if err != nil {
		return err
	}

	chunk = g.placeBlocks(chunk)

	err = g.W.SetChunk(chunk, true)
	return err
}

// First stage
func (g *Generator) getBiome(position vec.Vec2I) Biome {
	return Biome{}
}

// Second stage: Returns a chunk, at this stage it only sets the block id to one of 2, 0 for air and 1 for land
// Needs to be tweaked
func (g *Generator) generateLandscape(position vec.Vec2I, biome Biome) *Chunk {
	wHeight := g.W.WorldHeight
	lSize := g.W.LayerSize
	rough := int(10 - biome.Properties.Roughness)

	noiseGen := g.NoiseGenerators["landscape"]

	cWorldPos := position.Clone()
	cWorldPos.MultiplyScalar(float64(lSize))

	layers := make([]*Layer, wHeight)

	// The actual generation
	for x := 0; x < lSize; x++ {
		for y := 0; y < lSize; y++ {
			// layer
			for z := 0; z < wHeight; z++ {
				// world positions
				wx := cWorldPos.X + x
				wy := cWorldPos.Y + y

				l := layers[z]
				if len(l.Blocks) == 0 {
					l.Blocks = make([]*Block, lSize*lSize)
				}
				index := g.W.CoordsToIndex(vec.Vec3I{x, y, 0})
				noise := noiseGen.Noise3(float64(wx/rough), float64(wy/rough), float64(z/rough))
				noise -= float64((z / wHeight) * 5) //TODO: tweak this
				noise *= 100

				b := Block{
					LocalPosition: vec.Vec2I{x, y},
					Layer:         l,
					Id:            int(noise), // Actual id's are assigned later
				}
				l.Blocks[index] = &b
			}
		}
	}

	return &Chunk{
		Position: position,
		Layers:   layers,
		Biome:    biome,
	}
}

// Smooths the chunk edges
func (g *Generator) smoothEedges(chunk *Chunk) (*Chunk, ferr.FortiaError) {
	return &Chunk{}, nil
}

// Assigns proper blocks to everything, stone should be stone etc...
func (g *Generator) placeBlocks(chunk *Chunk) *Chunk {
	return &Chunk{}
}
