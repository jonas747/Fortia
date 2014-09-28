## Pre alha todo list

####Todo list

1. fortia packages
    - log (done)
    - error (done)
    - db (working on)
        + auth (done)
        + game (working on)
2. Rest Server (done)
3. Auth server 
    - Register
    - Login
    - Session token handling
4. Web frontend
5. Doc game server API
6. World ticker
    - Gen simple world
7. Game server 
8. Master server

#### Detailed
 - World ticker 
     + Simple package
         * Before game server is complete
             - only generate a world
         * After
             - Move units and such

 - Fortia pkgs
    - error
        + Simple error package (mostly copied from dropbox with slight changes)
    - log
        + Simple log client-server package
    - rest
        + Simple rest server
    - db
        + Various database functions for game and auth

 - Master server
    + Handle logging from game servers
        * Store them in leveldb
    + Simple web interface
        * Start stop world servers
        * View logs
            - Combined
            - Split between server
 - Game server
     + Handle ticks
    + Handle incoming requests
 - Auth server
    + Registration
    + Login
    + Recover password
    + Join worlds
 - Web frontend
    + Home
    + Game
