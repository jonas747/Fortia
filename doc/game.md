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

####/layers
Method: `GET`
 - x: list (seperated by ,) of x coordinates
 - y: list (seperated by ,) of y coordinates
 - z: list (seperated by ,) of z coordinates

#Types of game objects
Todo: write doc on thi
###ItemsconvertTemp

###Entity