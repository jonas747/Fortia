# Fortia

Fortia is a strategy game that takes place in a fantasy world. You start off with a small town with a small population, but if you play the game right you can end up with an empire with a population in the houndreds of thousands. 

You can be an aggressive player, fighting everyone you can manage, plundering and threatening your way up the ladder of world dominance.

You can be a team player forming alliances with other player and use diplomacy to climb the ladder.

You can be a strategist carefully planning your every move like setting up the battlefield between you and your enemy with tunnels and traps before the battle starts.

You can even trade yourself up the ladder by trading resources for protection or putting out hits on other players towns.

How you play the game is up to you, there is endless possibilities.

A world can last years without a winner. Depending on what time you have you can set tasks to be automated or you can manually do them for the absolute optimal speed to get those tasks done, although when you get a really high population micro management is impossible and you will have to rely on higher level of automation.

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

####Alpha

Alpha 0.1 is mostly the base work, the only thing that will work is view the world. There are no units or building just blocks

Alpha 0.2 Intruduces 2 units and 1 building and interaction between units from different players

Alpha 0.3 Dosent really add a lot of stuff, mostly backend stuff like the master server, server managegement, build system and load balancers is gonna be worked on here, improving the server infrastructure.
The only game related thing 0.3 is gonna add is making the world destructible by making units be able to destroy blocks(how long they use on each block depends on which block)

after this the fun work starts! (physics(water preasure, gravity, gas etc) is gonna be a hard one though...)