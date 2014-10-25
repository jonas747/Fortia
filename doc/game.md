#Game api
api url = {world}.fortia.com/

####/register
Method: `POST` 

Register the specified user up for this world

####/enter
Method: `POST`

Enters the world, returning a bunch of info

####/update
Method: `GET`

Returns info such as moved units or new messages since last update

####/chunk
Method: `GET`

Returns all the blocks visible to you in that chunk
params: x y

####/visiblechunks
Method: `GET`

Returns a list of chunks visible to you, or wether specified chunk is visible to you
optional params: x y

####/blocktypes
Method: `GET`

Returns all the block types with attributes for this server

#Types of game objects

Block types and entity types are stored in a json file. 
If a block type needs a custom handler it can register a custom handler

###Blocks

id - The unique id for this blocktype
name - human readable name for this blocktype
flags - list of flags for this blocktype
properties - proprties for this blocktype
    
default flags:
    solid

    {
        "id": 1,
        "name": "stone",
        "flags": [
            "solid",
            "mineable",
        ],
        "proprties": {
            "miningspeed": 100,
        }
    }
    {
        "id": 2,
        "name": "water",
        "flags": [
            "fluid",
            "freezeable"
        ],
        "proprties": {
            "fluidspeed": 10,
            "freezetemp": 0,
            "freezetoblockid": 4
        }
    }
    {
        "id": 2,
        "name": "lava",
        "flags": [
            "fluid",
            "heatemitter",
        ],
        "proprties": {
            "fluidspeed": 10,
        }
    }

flags and properties for those flags:

global properties:
    + mass
        * 1

 - solid
 - mineable
     + mineSpeed int
 - convert
     + convertTemp int
     + convertDamage int
 - fluid
     + fluidSpeed int
     + fluidLevel int
     + fluidSwimSpeed int
 - heatEmitter
     + heHeat int
     + heRadius int

###ItemsconvertTemp

###Entity