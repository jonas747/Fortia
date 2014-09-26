#DB layout

##Auth
Auth database is in redis
General stuff, online servers, configurations, user accounts etc...
####Keys
 - user:{name}
    + Hash
    + Represents a user
        * name   - Nameof the user
        * mail   - Email of the user
        * pw     - Bcrypt password hash
        * worlds - The worlds this user is on
 - token:{user}:{token}
     + Session token for the specified user


##Game
Game database is entirely in redis

Options for storing the world:
Benchmark to see which is faster

#####Single blocks
 - b:{xpos}:{ypos}:{zpos}
     + unsorted set
     + or list
         * 0 in the list = lowest entity on the tile (for example floor)
     + values in list/set:
         * entity id's
         
#####Chunks
 - c:{xpos}:{ypos}
     + list
         * index = local x, y, z
         * value
             - entity id's seperated by ":"

 - e:{eid}
     + hash
     + An entity(live entity, building, item, etc..)
         * type - Entity type
         * pos  - Position of entity
         * owner- Owner of this entity
         * other fields depending on which entity type
 - u:{name}
     + hash
     + World specific info about this user
    
##Log
Log database is in mongodb? maybe