## Layout
 - Master server
    + Handles logging for all fservers
    + Manages all the server instances, keeps them up to date etc..
 - fserver
     + fserver instances are controlled by the master server
     + the master server can tell an fserver instance to:
         * Update
         * start service
         * stop service
         * Change settings for service
     + Auth server
        + Handles non-game specific stuff
     + Game server
        + Handles game specific stuff
     + World ticker
        + Keeps the world going
 - Fortia pkgs
    + log
        * Provides a logging client and server
    + rest
        * Simple rest server and errors
    + vec
        * vector math
    + world
        * Functions for maipulating the world
    + db
        * Abstracts db interactions
    + common   
        * Common functions used throughout fortia

How this works:
Login
Client post /login -> load balancer -> auth server -> redis auth server


##Ticking:

Ticking the world is done by many ticker servers, distributing the ticking 

Ticking is done in stages:

1. non-world editing, stuff that does not change the world
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