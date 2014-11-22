package ticker

import (
	"github.com/jonas747/fortia/authserver"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/world"
	"time"
)

var (
	logger    *log.LogClient
	authDb    authserver.AuthDB
	gameDb    world.GameDB
	gameWorld *world.World
	handlers  map[string]Handler
)

type Handler func(*world.Action) ferr.FortiaError

func Run(l *log.LogClient, adb authserver.AuthDB, gdb world.GameDB) ferr.FortiaError {
	logger.Info("Running world ticker")
	logger = l
	authDb = adb
	gameDb = gdb

	gameWorld = &world.World{
		Logger: logger,
		Db:     gameDb,
	}
	err := gameWorld.LoadSettingsFromDb()
	if err != nil {
		return err
	}
	registerAllHandlers()
	run()
	return nil
}

func registerHandler(name string, handler Handler) {
	handlers[name] = handler
}

func registerAllHandlers() {
	//....
}

func run() {
	ticker := time.NewTicker(time.Duration(10) * time.Second)
	for {
		<-ticker.C
		logger.Debug("Ticking now...")
		startedTick := time.Now()

		nTick, err := gameDb.IncrTick()
		if err != nil {
			logger.Error(err)
			continue
		}

		for {
			action, err := gameDb.PopAction(nTick)
			if err != nil {
				if err.GetCode() == world.ErrCodeNotFound {
					taken := time.Since(startedTick)
					logger.Infof("Took %s to tick", taken.String())
					continue
				} else {
					logger.Error(err)
				}
			}

			err = processAction(action)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

func processAction(action *world.Action) ferr.FortiaError {
	handlerFunc, ok := handlers[action.Handler]
	if !ok {
		return ferr.New("Unknown handler in action \"", action.Name, "\", Handler: \"", action.Handler, "\"")
	}
	return handlerFunc(action)
}
