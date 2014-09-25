## Pre alha todo list

####Todo list

1. fortia packages
    - log (done)
    - error (done)
    - db (working on)
2. Auth server 
    - Register
    - Login
    - Session token handling
3. Doc game server API
4. Web frontend
5. Game server 
6. Master server

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
