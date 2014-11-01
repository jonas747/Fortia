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


Ticking:

Tasks are in a unordered set
the key of that set is the tick they are executed at

Tasks are distributed among the available tickservers so that they are as little idle as possible during each tick waiting for locks.

example:

Tick:100
actionread
actionread
actionwrite
actionwrite
actionread
actionwrite

aviable servers 2
The scheduler gives the following actions to the tickservers

s1:
actionwrite
actionread
actionwrite
s2:
actionread
actionwrite
actionread

Although not perfect its the way the sceduler works