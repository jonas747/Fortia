## Stages
Worldgen is a process done in multiple stages:

1. Landscape, caves and whatnot
2. Smoothing between biomes
3. Block placement
4. Vegetation
5. Flagging

### Landscape, caves, Rivers etc
This is where all the landscape is actually generated, no actual blocks is placed yet, but each block has a id with the "life" variable

### Smoothing between biomes
It can be a bit rough in the transition between biomes, this stages smoothes things out a little bit

### Block placement
This stage actually places the blocks based on the block's life variable assigned in the first stage

### Flagging
Flag surface blocks for vegetation

### Vegetation
Here trees, plants etc are placed.

### Flagging
All The blocks have certain flags, such as wether that block is fully hidden or not, if its connected to the ground, wetherer it is covered and so on.

## Block classes

1. Soil
2. Deep soil
3. Metal
4. Rock
5. Deep rock
6. Wood
7. Vegetation
8. Gem
9. Deep gem
10. Water
11. Magma

## Layers

### Surface
 - soil
 - wood
 - vegetation
 - Water

The surface layer, depending on biome full, full of trees happy things and whatnot

### Dirt
 - soil
 - Water

Mostly dirt on this layer

### ... 
 - deep soil
 - rock

### ...
 - deep soil
 - rock
 - metals
 - gems

### ...
 - deep rock
 - metals
 - deep metals
 - gems

### Hell 
 - deep rock
 - deep metals
 - deep gems
 - magma

This layer is full of lava, dangerous creatures and valuable recources
