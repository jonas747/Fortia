##Alha todo list

###Todo list
####Alpha 0.1
1. fortia packages
    - log (done)
    - error (done)
    - db (working on)
        + auth (done)
        + game (working on)
2. Rest Server (done)
3. Auth server (mostly done)
    - Register (done)
    - Login (done)
    - Session token handling (done)
    - Get avaiable worlds
4. Web frontend (mostly done)
    - login (done)
    - lobby (working on)
        + Only front page and world joining
5. Doc game server API (working on)
6. World ticker
    - Gen simple world
7. Game server
    - Serve the world 
8. Web frontend
    - view available worlds
    - join a world
    - View the world

First public release 0.2
####Alpha 0.2
1. Log server
    - Store logs in mongodb
2. Web frontend
    - Select tile(s)
        + view units and building on that tile and interact with them
3. Game server
    - Interactions with units and buildings

#####Alpha 0.2 Units
 - Actions
     + Move
     + Attack
     + Stop
     + Hold position
 - Monitor
     + Health

#####Alpha 0.2 Buildings
 - Actions
     + Self destruct
     + Cancel queue item
 - Monitor
     + Queue
     + Health
 
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
