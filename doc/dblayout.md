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
        * ---Optional:
        * worlds - The worlds this user is on, seperated with ":" 
        * role - 0 normal, 1 mod, 2 admin, 3 fish, 4 god
        * donor - donor leve
 - token:{user}:{token}
     + Session token for the specified user

 - worlds
    + unordered set
    + Contains a list of all world names
 
 - world:name
    + hash
    + info about a specified world
        * started   - when the world started
        * players   - number of players 
        * size      - size in chunks

##Game
Game database is entirely in redis

- worldInfo
    + hash
    + general information about this server(name...)
        * name
        * size
        * players

went with chunks before:....
 - c:{xpos}:{ypos}
     + list
         * index = local x, y, z
         * value
             - entity id's seperated by ":" with first one being block id
going with layers
- l:{xpos}:{ypos}:{layer(or zpos if you want)}
     + json of block info
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
    
Entity specific properties:
 - Units
     + CurrentAction
         * An action such as MOVE, ATTACK, etc..
     + Health
 - Buldings
     + Health

 - Unit Action fields
     + Move
         * MovePos
         * WhenOff
             - Unix timestamp of when this unit is supposed to be off this tile(and onto another)
##Login
Log database is in mongodb? maybe

