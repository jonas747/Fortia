#Game api
api url = {world}.fortia.com/

####/info
Method: `GET` 

Returns info about this world
such as blocktpes, size etc..

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

#Types of game objects
Todo: write doc on thi
###ItemsconvertTemp

###Entity