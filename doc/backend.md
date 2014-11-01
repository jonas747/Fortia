####Parts

 - Master server
    + Handles logging for all game servers
    + Manages all the server instances, keeps them up to date etc..
 - Auth server
    + Handles non-game specific stuff
 - Game server
    + Handles game specific stuff
 - World ticker
    + Keeps the world going
 - Web frontend
    + bootstrap
 - Fortia pkgs
    + log
        * Provides a logging client and server
    + rest
        * Simple rest server
    + resterrors
        * common errors returned by the rest servers
    + vec
        * vector math
    + world
        * Functions for maipulating the world
    + db
        * Abstracts db interactions

How this works:
Login
Client post /login -> load balancer -> auth server -> redis auth server
                                      auth server
                                      ...........
