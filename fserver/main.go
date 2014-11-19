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
	"runtime"
)

var (
	fUpdateWorld = flag.Bool("u", false, "Updates the world with the wgen.json and blocktypes.json, then exits")
	fCreateWorld = flag.Bool("c", false, "Creates a world with settings from world.json, wgen.json and blocktypes.json")
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
	runtime.GOMAXPROCS(runtime.NumCPU())

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
		createWorld()
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

func createWorld() {
	// Load the world info json
	var biomesInfo world.BiomesInfo
	err := common.LoadJsonFile("biomes.json", &biomesInfo)
	panicErr(err)

	var btypes []world.BlockType
	err = common.LoadJsonFile("blocks.json", &btypes)
	panicErr(err)

	winfo := &world.WorldInfo{
		BlockTypes: btypes,
		Biomes:     biomesInfo,
	}

	err = common.LoadJsonFile("world.json", winfo)

	world := &world.World{
		Logger:      logger,
		Db:          gdb,
		GeneralInfo: winfo,
	}
	err = world.SaveSettingsToDb()
	panicErr(err)

	generator := world.NewGenerator()
	generator.Size = winfo.Size

	err = generator.GenerateWorld()
	panicErr(err)
}
