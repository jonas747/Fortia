package world

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/cheggaaa/pb"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/messages"
	"github.com/jonas747/fortia/simplex"
	"github.com/jonas747/fortia/vec"
	"math/rand"
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
	World           *World // The world
	Settings        *messages.WorldSettings
	NoiseGenerators map[string]*simplex.Noise
	Size            int
}

// Returns a new generator
func NewGenerator(world *World, seed int64) *Generator {
	// map the blocktypes to biomes for ease of access
	// Init the noise generators
	ngens := make(map[string]*simplex.Noise)
	rng := rand.New(rand.NewSource(seed))

	ngens["landscape"] = simplex.NewNoise(rng)
	ngens["caves"] = simplex.NewNoise(rng)

	generator := &Generator{
		World:           world,
		NoiseGenerators: ngens,
	}

	return generator
}

// Generates a world
func (g *Generator) GenerateWorld() ferr.FortiaError {
	// Start by generating the base chunks
	g.World.Logger.Info("Generating landscape")
	p := pb.StartNew(g.Size * g.Size)
	for x := 0; x < g.Size; x++ {
		for y := 0; y < g.Size; y++ {
			pos := vec.Vec2I{x, y}
			chunk, err := g.generateBaseChunk(pos)
			if err != nil {
				return err
			}
			// Save the chunk
			err = g.World.SetChunk(chunk)
			if err != nil {
				return err
			}
			p.Increment()
		}
	}
	p.Finish()

	// Caves
	g.World.Logger.Info("Generating caves")
	g.GenStage(g.generateCaves)

	// Smooth between chunks
	// g.World.Logger.Info("Smoothing chunk borders")
	// g.GenStage(g.smoothChunk)

	// Grow trees
	// More advanced block placement
	// Flag blocks
	g.World.Logger.Info("Placing blocks")
	g.GenStage(g.basePlaceBlocks)

	g.World.Logger.Info("Flagging hidden blocks and layers")
	g.GenStage(g.flagHidden)

	return nil
}

func (g *Generator) GenStage(f func(*Chunk) ferr.FortiaError) {
	p := pb.StartNew(g.Size * g.Size)
	defer p.Finish()
	for x := 0; x < g.Size; x++ {
		for y := 0; y < g.Size; y++ {
			chunk, err := g.World.GetChunk(vec.Vec2I{x, y})
			if err != nil {
				g.World.Logger.Error(err)
				continue
			}
			err = f(chunk)
			if err != nil {
				g.World.Logger.Error(err)
				continue
			}
			err = g.World.SetChunk(chunk)
			if err != nil {
				g.World.Logger.Error(err)
				continue
			}

			p.Increment()
		}
	}
}

func (g *Generator) flagHidden(chunk *Chunk) ferr.FortiaError {
	chunk.FlagHidden(map[vec.Vec2I]*Chunk{})
	return nil
}

// 0.2
func (g *Generator) ExpandWorld() {

}

func (g *Generator) smoothChunk(chunk *Chunk) ferr.FortiaError {
	return nil
}

// Generates a chunk, saves chunk information and layers to db.
func (g *Generator) generateBaseChunk(position vec.Vec2I) (*Chunk, ferr.FortiaError) {
	biome, potency, err := g.getBiome(position)

	chunk := g.generateLandscape(position, biome)
	if err != nil {
		return nil, err
	}

	chunk.RawChunk.Potency = proto.Int(potency)
	return chunk, err
}

// First stage
// Returns the biome and the potency
// the dominating biome is the one with the most potency? or the ones with most neighbours?(going with potency for now, experimenting later)
func (g *Generator) getBiome(position vec.Vec2I) (*messages.Biome, int, ferr.FortiaError) {
	highestPotency := 0
	highestBiome := &messages.Biome{}
	// Get souroding chunks
	for x := -1; x < 1; x++ {
		for y := -1; y < 1; y++ {
			chunk, err := g.World.GetChunk(vec.Vec2I{x + position.X, y + position.Y})
			if err != nil {
				if err.GetCode() != 404 {
					return &messages.Biome{}, 0, err
				}
			}

			// skip if the chunk is not found
			if chunk == nil {
				continue
			}
			if int(chunk.RawChunk.GetPotency()) > highestPotency {
				highestPotency = int(chunk.RawChunk.GetPotency())
				highestBiome = g.World.GetBiomeFromId(int(chunk.RawChunk.GetBiome()))
			}
		}
	}

	// if there were no souronding chunks assign a random biome based on biome probabilities
	if highestPotency == 0 {
		// Get a random biome
		biomes := make([]int, 0)
		for k, v := range g.World.Settings.Biomes.Biomes {
			for i := 0; i < int(v.GetProbability()); i++ {
				biomes = append(biomes, k)
			}
		}
		bNum := biomes[rand.Intn(len(biomes))]
		highestBiome = g.World.Settings.Biomes.Biomes[bNum]
		highestPotency = 5 // Add one because we subtract one later
	}
	return highestBiome, highestPotency - 1, nil
}

// Second stage: Returns a chunk, at this stage it only sets the block id to one of 2, 0 for air and 1 for land
// Needs to be tweaked
func (g *Generator) generateLandscape(position vec.Vec2I, biome *messages.Biome) *Chunk {
	chunkHeight := int(g.World.Settings.GetChunkHeight())
	chunkWidth := int(g.World.Settings.GetChunkWidth())
	noiseGen := g.NoiseGenerators["landscape"]

	//rough := biome.Properties.Roughness
	cWorldPos := position.Clone()
	cWorldPos.MultiplyScalar(float64(chunkWidth))

	blocks := make([]*messages.Block, chunkWidth*chunkWidth*chunkHeight)

	c := &Chunk{
		World: g.World,
		RawChunk: &messages.Chunk{
			Blocks: blocks,
			Biome:  proto.Int32(biome.GetId()),
			X:      proto.Int(position.X),
			Y:      proto.Int(position.Y),
		},
	}

	// The actual generation
	for x := 0; x < chunkWidth; x++ {
		for y := 0; y < chunkWidth; y++ {
			for z := 0; z < chunkHeight; z++ {
				// world position
				wx := cWorldPos.X + x
				wy := cWorldPos.Y + y

				noise := noiseGen.Noise3(float64(wx)/float64(50), float64(wy)/float64(50), float64(z)/float64(50))
				life := 7 - 10*(float64(z)/float64(chunkHeight)) // Conrolls the elevation
				//g.World.Logger.Info(z, life, noise)
				life += noise / 4 // This controlls the sensitivity
				life *= 100

				raw := &messages.Block{
					Kind: proto.Int32(int32(life)),
				}
				index := g.World.CoordsToIndex(vec.Vec3I{x, y, z})
				blocks[index] = raw
			}
		}
	}
	return c
}

// Generates caves underground
func (g *Generator) generateCaves(chunk *Chunk) ferr.FortiaError {
	noiseGen := g.NoiseGenerators["caves"]

	chunkWidth := int(g.World.Settings.GetChunkWidth())
	chunkHeight := int(g.World.Settings.GetChunkHeight())

	cWorldPos := vec.Vec2I{int(chunk.RawChunk.GetX()), int(chunk.RawChunk.GetY())}
	cWorldPos.MultiplyScalar(float64(chunkWidth))

	for x := 0; x < chunkWidth; x++ {
		for y := 0; y < chunkWidth; y++ {
			for z := 0; z < chunkHeight; z++ {
				wx := float64(cWorldPos.X + x)
				wy := float64(cWorldPos.Y + y)
				life := 6 - 10*(float64(z)/float64(chunkHeight))
				life += (noiseGen.Noise3(wx/float64(40), wy/float64(40), float64(z)/float64(40))) * 2
				if life > 0 {
					index := g.World.CoordsToIndex(vec.Vec3I{x, y, z})
					currentLife := chunk.RawChunk.Blocks[index].GetKind()
					currentLife -= int32(life * 100)
					chunk.RawChunk.Blocks[index].Kind = proto.Int32(currentLife)
					//l.Blocks[index].Id = 0
				}
			}
		}
	}

	return nil
}

// Smooths the chunk edges
func (g *Generator) smoothEedges(chunk *Chunk) (*Chunk, ferr.FortiaError) {
	return chunk, nil
}

// Assigns proper blocks to everything, stone should be stone etc...
// TODO: More advanced block placement
func (g *Generator) basePlaceBlocks(chunk *Chunk) ferr.FortiaError {
	for _, block := range chunk.RawChunk.Blocks {
		life := block.GetKind()
		if life > 50 {
			block.Kind = proto.Int32(1) // rock
		} else if life <= 50 && life > 0 {
			block.Kind = proto.Int32(2) // grass
		} else {
			block.Kind = proto.Int32(0) // Air
		}
	}
	return nil
}

// TODO
func (g *Generator) GrowTree() {

}
