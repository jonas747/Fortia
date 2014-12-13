#DB layout

##Auth
Auth database is in redis, for the time being
General stuff, online servers, configurations, user accounts etc...
####Keys
 - u:{name}
    + Hash
    + Represents a user
        * name   - Name of the user
        * mail   - Email of the user
        * pw     - Bcrypt password hash
        * ---Optional:
        * worlds - The worlds this user is on, seperated with ":" 
        * role - 0 normal, 1 mod, 2 admin, 3 fish, 4 god
        * donor - donor leve
 - t:{token}
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
        * 
 - c:{xpos}:{ypos}
     + protobuf with chunk data
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
 - ue:{name}
     + Set of entity id's
     + User Entities 
     + List of entities this user owns
 - actionQueue:{tick}
     + set of Actions
     + Actions to be executed at specified tick
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

