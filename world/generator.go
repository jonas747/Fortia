package world

import (
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/simplex"
	"github.com/jonas747/fortia/vec"
	"math/rand"
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
	biome, potency, err := g.getBiome(position)

	chunk := g.generateLandscape(position, biome)
	if err != nil {
		return err
	}

	chunk, err = g.smoothEedges(chunk)
	if err != nil {
		return err
	}
	chunk.Potency = potency

	chunk = g.placeBlocks(chunk)

	err = g.W.SetChunk(chunk, true)
	return err
}

// First stage
// Returns the biome and the potency
// the dominating biome is the one with the most potency? or the ones with most neighbours?(going with potency for now, experimenting later)
func (g *Generator) getBiome(position vec.Vec2I) (Biome, int, ferr.FortiaError) {
	highestPotency := 0
	highestBiome := Biome{}
	// Get souroding chunks
	for x := -1; x < 1; x++ {
		for y := -1; y < 1; y++ {
			chunk, err := g.W.GetChunk(x+position.X, y+position.Y, false)
			if err != nil {
				return Biome{}, 0, err
			}

			// skip if the chunk is not found
			if chunk == nil {
				continue
			}
			if chunk.Potency > highestPotency {
				highestPotency = chunk.Potency
				highestBiome = chunk.Biome
			}
		}
	}

	// if there were no souronding chunks assign a random biome based on biome probabilities
	if highestPotency == 0 {
		// Get a random biome
		biomes := make([]int, 0)
		for k, v := range g.Biomes.Biomes {
			for i := 0; i < int(v.Probability); i++ {
				biomes = append(biomes, k)
			}
		}
		bNum := biomes[rand.Intn(len(biomes))]
		highestBiome = g.Biomes.Biomes[bNum]
		highestPotency = 5 // Add one because we subtract one later
	}
	return highestBiome, highestPotency - 1, nil
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
				if l == nil {
					l = &Layer{
						World:    g.W,
						Position: vec.Vec3I{position.X, position.Y, z},
					}
					layers[z] = l
				}
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

// TODO
func (g *Generator) caves(chunk *Chunk) *Chunk {
	return chunk
}

// Smooths the chunk edges
func (g *Generator) smoothEedges(chunk *Chunk) (*Chunk, ferr.FortiaError) {
	return chunk, nil
}

// Assigns proper blocks to everything, stone should be stone etc...
// TODO: More advanced block placement
func (g *Generator) placeBlocks(chunk *Chunk) *Chunk {
	for x := 0; x < g.W.LayerSize; x++ {
		for y := 0; y < g.W.LayerSize; y++ {
			for z := 0; z < g.W.WorldHeight; z++ {
				l := chunk.Layers[z]
				index := g.W.CoordsToIndex(vec.Vec3I{x, y, 0})
				b := l.Blocks[index]
				if b.Id > 50 {
					b.Id = 1 // rock
				} else if b.Id <= 50 && b.Id > 0 {
					b.Id = 2 // grass
				} else {
					b.Id = 0 // Air
				}
			}
		}
	}
	return chunk
}
