## Pre alha todo list

####Todo list

1. fortia packages
    - log (done)
    - error (done)
    - db (working on)
        + auth (done)
        + game (working on)
2. Rest Server
3. Auth server 
    - Register
    - Login
    - Session token handling
4. Doc game server API
5. Web frontend
6. Game server 
7. Master server

#### Detailed

 - Fortia pkgs
    - error
    - log
    - db

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
