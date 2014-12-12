#Forta Action system

In fortia when you want to move a unit from say [0,0,0] to [0,0,10]
Then the gameserver creates an action that has the time index of the time it takes to travel 1 tile + the currend time index, the action besically says "move 1 tile closer to [0,0,10]".

Every tick, actions are designated onto tickers by the scheduler. Because of this modular system many "ticker" servers can be added at any time to improve the speed of ticking the world.