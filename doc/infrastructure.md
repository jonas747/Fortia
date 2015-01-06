## Layout
 - logserver
     + Compiles all logs from all servers
 - fserver
     + fserver instances fetch configurations from a database
     + Sending commands to these servers is done through a redis database (update, configuration etc..)
         * Removes the single point failure that is a single master server, can have redis slaves that gets promoted if/when one dies
     + Serivces:
         * Auth
             - Handles non-game specific stuff
             - Servers a public rest API
         * Game
             - Handles game specific stuff
             - serves a public rest API
         * World ticker
             - Keeps the world going
         * Scheduler
             - Tells al tickers when and what to tick
 - Fortia pkgs
    + log
        * Provides a logging client and server
    + rest
        * Simple rest server and errors
    + vec
        * vector math
    + world
        * Functions for maipulating the world
    + rdb
        * authserver/db and world/db implementations for redis
    + common   
        * Common functions used throughout fortia

How this works:
Login
Client post /login -> load balancer -> auth server -> redis auth server

Current parts of the fortia backend thats scaleable:

 - API
     + The fortia api is stateless so servers can be added and removed without any hassle

To be made scaleable

 - DB
     + Data is sharded between multiple databases, to spread out the load towards mutliple servers, the only penalty for adding/removing servers is the rebalancing
 - Logging

Maybe in the future

 - World scheduler
 - Master server
     + The master server is supposed to handle all the other servers, So with too many servers this might need to be made scaleable across multiple servers

##Sharding the data between redis instances

World data is sharded between several instances to spread the load and to increase the possible world size.

It is sharded by hashing the key and then take the modulo of whatever the number of active redis shards is

The key is also prepended a database version number, so the final key format is:
v{version}:{key}

The database version number increases everytime a redis server is added or removed.
This is to make the rebalancing simpler and also makes the world viewable while rebalancing

This has some drawbacks:

 - When adding a redis server everything will have to be rebalanced between all servers
 - Can't use more advanced redis features like transactions

It has however some advantaes too:

 - If a redis instance(and its backup/slave) dies, only the world data on that server is not available (the world will however pause as soon as a server dies but will still be viewable from clients, in read only mode)

##Ticking:

Ticking the world is done by many ticker servers, distributing the ticking 

Ticking is done in stages:

1. non-world editing, stuff that does not change the world(individual entities may change though)
2. world editing stage, here all the actions that are changing the world processed, (mined blocks, buildings placed etc)
3. entity moving stage, all the entities are moved, calculated paths for etc...


intruduces some problems:

###Problems with distributed tickers:
Because chunks are stored in the protobuf format, updating just one field to say add a entity to a block, requires updating the whole chunk. Which has to be done in 3 steps, 1 - retrieve and decode chunk. 2 - apply changes to chunk. 3 - Store the chunk
Because of this a chunk can only be edited by 1 worker at the same time.

Solutions: 
    
- Instead of directly applying changes to the world, queue the changes up and apply them all at the end
- Use some sort of locking

Either one will be testes to see which works out the best on large and small scale when time comes