# Fortia

## Layout

Parts:
    - Master server
        + Handles logging for all game servers
        + Servers a web page
    - Auth server
        + Handles non-game specific stuff
    - Game server
        + Handles game specific stuff
    - Web frontend
        + bootstrap
    - Fortia Lib
        + mostly database stuff

Joining a world:
    1. Click join world button
    2. Send REST request to the auth server
    3. Auth server checks session token etc
    4. Auth server makes another rest request to the specified game server
    5. Game server adds an entry for the user in the game db

Databases:
    - Auth
        + For anything that dosent belong in a world database
        + password hashes
        + usernames
        + emails
        + etc... 
    - Game (per world)
        + For game speicic stuff
        + The game world itself
        
Users have to create new names each time they join a world

## Pre-alpha state
Game server config is done via json files for pre-alpha stage
Game:
    - 3d world viewed in 2d (dwarf fortress style)
    - Simple structures
        + barracks
            * Traines population into warriors
    - Simple ui
        + Create buildings
        + Population
            * warrios
            * farmers
            * etc..
        + Resources
    
Account mgmnt:
    - Simple login (username password)
    - Simple registration (email user passx2)
        + bcrypt
    - Select world

Master server:
    - View all game servers
    - Logs

Game server
    - 1 game server per world 

Todo list:
    
    1. Fortia pkg
    2. Auth server (reg, login)
    3. Doc game server API
    4. Web frontend
    5. Game server 
    6. Master server

    - Fortia pkg
        + Logging to master server
        + Auth db
        + Game db
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
