package main

import (
	"flag"
	"github.com/jonas747/fortia/authserver"
	"github.com/jonas747/fortia/common"
	"github.com/jonas747/fortia/gameserver"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/rdb"
	//"github.com/jonas747/fortia/ticker"
	"github.com/jonas747/fortia/world"
)

var (
	fUpdateWorld = flag.Bool("u", false, "Updates the world with the wgen.json and blocktypes.json, then exits")
	fCreateWorld = flag.Bool("c", false, "Creates a world with settings from world.json, wgen.json and blocktypes.json, then exits")
)

var (
	gdb    *rdb.GameDB
	adb    *rdb.AuthDB
	logger *log.LogClient
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	var config Config
	err := common.LoadJsonFile("config.json", &config)
	panicErr(err)

	logger, err = log.NewLogClient("", -1, "Fluffy")
	if err != nil {
		if logger != nil {
			logger.Warn("Error connecting log client to server", err.Error())
		} else {
			panicErr(err)
		}
	}

	// Connect to databases and master server
	// TODO: Master server
	if config.ServeAuth || config.ServeGame || *fCreateWorld {
		if config.ServeGame || *fCreateWorld {
			gdbRaw, err := rdb.NewDatabase(config.Game)
			panicErr(err)
			gdb = &rdb.GameDB{gdbRaw}
		}
		adbRaw, err := rdb.NewDatabase(config.Auth)
		panicErr(err)
		adb = &rdb.AuthDB{adbRaw}
	}
	// Gen world initialises a new world
	if *fCreateWorld {
		createWorld(6)
		logger.Info("Generated world, exiting...")
		return
	}
	if config.ServeGame {
		go gameserver.Run(logger, gdb, adb, ":8081")
	}
	if config.ServeAuth {
		go authserver.Run(logger, adb, ":8080")
	}
	if config.RunTicker {
		//go ticker.Run(logger, adb, gdb)
	}
	select {}
}

func createWorld(size int) {
	world := &world.World{
		Logger: logger,
		Db:     gdb,
	}
	err := world.LoadSettingsFromDb()
	panicErr(err)
}
