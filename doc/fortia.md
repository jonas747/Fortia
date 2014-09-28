# Fortia

####Parts

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

####Joining a world

1. Click join world button
2. Send REST request to the auth server
3. Auth server checks session token etc
4. Auth server makes another rest request to the specified game server
5. Game server adds an entry for the user in the game db

####Databases

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
Units walk through building and can fly&walk thorugh ground
(pathfinding first thing alpha)
####Game
 - 3d world viewed in 2d (dwarf fortress style)
 - Simple structures
 - Simple ui

More detailed:

 - Units
     + Farmers
         * Generates food
     + Warriors (Trained from farmers)
         * 5 times better than farmers at fighting
         * Consumes food
 - Buildings
     + Barracks
         * Trains Farmers into warriors
         * Cost: none
 - World gen
     + Simple
     + block types:
         * Grass
         * Rock
 - Resource gathering
     + Food (farmers by being around generates food)
 - World altering
     + none
 - Pathfinding
     + Try to implement some cheap pathfinding, last priority
 - Fighting
     + Simple
         * Warriors select the closest target withing a set radius
 - UI
     + View food
     + Click on a tile to open a menu
         * Move the units on that tile
         * Build a building on that tile

####Account management
 - Simple login (username password)
 - Simple registration (email user passx2)
    + bcrypt
 - Select world

####Master server
 - View all game servers
 - Logs

####Fortia-auth
 - Auth rest server

####Fortia-game
 - Game rest server
 
####World ticker
 - Keeps a world spinning