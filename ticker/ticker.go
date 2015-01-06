package ticker

import (
	"github.com/jonas747/fortia/db"
	"github.com/jonas747/fortia/errorcodes"
	"github.com/jonas747/fortia/errors"
	"github.com/jonas747/fortia/game"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/messages"
	"time"
)

var (
	logger    *log.LogClient
	authDb    db.AuthDB
	gameDb    db.GameDB
	gameWorld *game.World
	handlers  map[string]Handler
)

type Handler func(*messages.Action) errors.FortiaError

func Run(l *log.LogClient, adb db.AuthDB, gdb db.GameDB) errors.FortiaError {
	logger = l
	logger.Info("Running world ticker")
	authDb = adb
	gameDb = gdb

	gameWorld = &game.World{
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
	ticker := time.NewTicker(time.Duration(1) * time.Second)
	for {
		<-ticker.C
		logger.Debug("Ticking now...")
		startedTick := time.Now()

		nTick, err := gameDb.IncrTick()
		if err != nil {
			logger.Error(err)
			continue
		}
		// 0 is an alias for immediate, things such as placing units gets put with tick number 0 to process
		// the action as soon as possible
		num, err := processTick(0)
		if err != nil {
			logger.Error(err)
		}
		num2, err := processTick(nTick)
		if err != nil {
			logger.Error(err)
		}

		taken := time.Since(startedTick)
		logger.Debugf("Finnished tick [%d] took %s, Processed: %d action(s)", nTick, taken.String(), num+num2)
	}
}

// Returns the number of actions processed and errors if any
func processTick(curTick int) (int, errors.FortiaError) {
	stage := 1
	numActions := 0
	for {
		action, err := gameDb.PopAction(curTick, stage)
		if err != nil {
			if err.GetCode() == errorcodes.ErrorCode_RedisKeyNotFound {
				stage++
				if stage > 3 {
					return numActions, nil
				}
				continue
			} else {
				return numActions, err
			}
		}

		err = processAction(action)
		numActions++
		if err != nil {
			return numActions, err
		}
	}
	return numActions, nil
}

func processAction(action *messages.Action) errors.FortiaError {
	handlerFunc, ok := handlers[action.GetHandler()]
	if !ok {
		return errors.New(errorcodes.ErrorCode_NoHandlerForAction, "No handler for action %s, action.Handler: %s", action.GetName(), action.GetHandler())
	}
	return handlerFunc(action)
}
