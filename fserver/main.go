package main

import (
	"flag"
	"github.com/jonas747/fortia/authserver"
	"github.com/jonas747/fortia/common"
	"github.com/jonas747/fortia/db"
	"github.com/jonas747/fortia/gameserver"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/world"
)

var (
	fUpdateWorld = flag.Bool("-u", false, "Updates the world with the wgen.json and blocktypes.json")
	fCreateWorld = flag.Bool("-c", false, "Creates a world with settings from world.json, wgen.json and blocktypes.json")
)

var (
	gdb    *db.GameDB
	adb    *db.AuthDB
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
	if config.ServerAuth || config.ServeGame || *fCreateWorld {
		if config.ServeGame || *fCreateWorld {
			gdbRaw, err := db.NewDatabase(config.Game)
			panicErr(err)
			gdb = &db.GameDB{gdbRaw}
		}
		adbRaw, err := db.NewDatabase(config.Auth)
		panicErr(err)
		adb = &db.AuthDB{adbRaw}
	}
	// Gen world initialises a new world
	if *fCreateWorld {
		createWorld()
	}
	if config.ServeGame {
		go gameserver.Run(logger, gdb, adb, ":8081")
	}
	if config.ServerAuth {
		go authserver.Run(logger, adb, ":8080")
	}
	select {}
}

func createWorld() {
	// Load all the settings from files
	var winfo WorldConfig
	err := common.LoadJsonFile("world.json", &winfo)
	panicErr(err)

	logger.Info("Creating world", winfo.Name)

	btypes, err := world.BlockTypesFromFile("blocks.json")
	panicErr(err)

	biomes, err := world.BiomesFromFile("wgen.json")
	panicErr(err)

	// Set the world info to the auth db
	err = adb.SetWorldInfo(winfo.Name, map[string]interface{}{"size": winfo.LayerSize, "name": winfo.Name})
	panicErr(err)

	w := &world.World{
		Logger:      logger,
		Db:          gdb,
		LayerSize:   winfo.LayerSize,
		WorldHeight: winfo.WorldHeight,
		BlockTypes:  btypes,
		Biomes:      biomes,
	}
	// Save the settings to the game db
	err = w.SaveSettingsToDb()
	panicErr(err)
}
