##Alha todo list

###Todo list
####Alpha 0.1
1. fortia packages
    - log (done)
    - error (done)
    - db (working on)
        + auth (done)
        + game (partially done)
2. Rest Server (done)
3. Auth server (mostly done)
    - Register (done)
    - Login (done)
    - Session token handling (done)
    - Get avaiable worlds (done)
4. Web frontend (mostly done)
    - login (done)
    - lobby (partially done)
        + servers and world joining (done)
5. Doc game server API (partially done)
6. World ticker (done)
    - Gen simple world (done)
7. Game server(working on)
    - Serve the world 
8. Web frontend(working on)
    - join a world <- big one
    - View the world
        + engine
            * load blocktypes and textures
            * handle input
            * draw loop

######Sudden things to do
 - figure out which way to store block types with attributes in both go and javascript
     + simple json file?
    
First public release 0.2
####Alpha 0.2
alpha 0.2:
 - vector math lib 
     + change all stuff to use the lib
 - selection
 - Build buildings
 - Train units
 - somewhat more advanced world gen
 - Extended lobby
     + Admin panel
     + User global settings
     + News
     + Help
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

####Post alpha 0.2
 - Chat
     + different channels
 - Websockets as an alternative backend to rest
 - optional 3d renderer
     + If you have a good computer you can use the 3d-
     + renderer, or if you dont you can use the 2d one
 - Build system
     + Numbered builds
 - Switch easily between dev and production modes (add -dev flag)
 - Easy server deployment
    - Simple daemon that is on the server, downloads updates and such for whatever the server is running
    - tell the daemon to install a server via web admin panel
 - Password recovery through email
 - 2fa maybe?

 
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
