# Scheduler

The fortia scheduler controlls the world ticking of a single world.

Once the scheduler wants the world to move forward by one tick, it sends a signal to each ticker with the stage a stage, once out of work each ticker replies that they finnished with that stage, and once all the tickers are done with that stage the ticker signals them to start the next stage untill there are no more stages, then the ticker waits untill the next tick is suppossed to happen.