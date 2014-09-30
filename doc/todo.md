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
6. World ticker (working on)
    - Gen simple world
7. Game server
    - Serve the world 
8. Web frontend
    - view available worlds
    - join a world
    - View the world

First public release 0.2
####Alpha 0.2
alpha 0.2:
 - vector math lib 
     + change all stuff to use the lib
 - selection
 - Build buildings
 - Train units

######Avaiable buildings for alpha 0.2
 - Barracks
    + Train farmers/builders(default units) into warriors

######Avaiable units:
 - farmers
     + Default unit
 - warrior
     + 3 times tronger than farmers
     + takes up twice as much food

######Alpha 0.2 Units
 - Actions
     + Move
     + Attack
     + Stop
     + Hold position
 - Monitor
     + Health

######Alpha 0.2 Buildings
 - Actions
     + Self destruct
     + Cancel queue item
 - Monitor
     + Queue
     + Health
 
#### notes
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
