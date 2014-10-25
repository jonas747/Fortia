package main

import (
	//"fmt"
	"github.com/jonas747/fortia/db"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/vec"
	"github.com/jonas747/fortia/world"
	//"math/rand"
	// "strconv"
	// "strings"
	"time"
)

var (
	logger     *log.LogClient
	authDb     *db.AuthDB
	gameDb     *db.GameDB
	blockTypes []*BlockType
	config     *Config
	gameWorld  *world.World
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	c, err := loadConfig("config.json")
	panicErr(err)
	config = c

	l, nErr := log.NewLogClient(config.LogServer, -1, "authAPI")
	logger = l
	if nErr != nil {
		l.Warn(ferr.Wrap(nErr, "Error connecting to logserver, client wont send log messages"))
	}

	l.Info("Log client init successful! Creating database connection pools...")

	gdb, nErr := db.NewDatabase(config.GameDb)
	if nErr != nil {
		l.Fatal(ferr.Wrap(nErr, ""))
		return
	}

	gameDb = &db.GameDB{gdb}

	l.Info("Game db init sucessfull, Reading blcoktypes now")

	blockTypes, err = loadBlockTypes("blocks.json")
	if err != nil {
		l.Fatal(err)
		return
	}

	l.Info("Blocktypes loaded, initialising world")
	gameWorld = &world.World{
		Logger:      logger,
		Db:          gameDb,
		LayerSize:   50,
		LayerHeight: 100,
	}

	l.Info("World ticker started sucessfully")
	generate(2, 2)
	run()
}

func generate(numX, numY int) {
	num := 0
	for x := 0; x < numX; x++ {
		for y := 0; y < numY; y++ {
			for z := 0; z < 100; z++ {
				num++
				layer := gameWorld.GenLayer(vec.Vec3I{x, y, z})
				err := gameWorld.SetLayer(layer)
				if err != nil {
					logger.Error(err)
				}
				logger.Info("Generated layer ", num, " out of ", numX*numY*100)
			}
		}
	}
}

func run() {
	ticker := time.NewTicker(time.Duration(10) * time.Second)
	for {
		logger.Debug("Ticking now...")
		startedTick := time.Now()
		// keys, err := gameDb.GetChunkList()
		// if err != nil {
		// 	logger.Error(err)
		// 	continue
		// }
		// finChan := make(chan bool)
		// for _, v := range keys {
		// 	split := strings.Split(v, ":")
		// 	x, _ := strconv.Atoi(split[1])
		// 	y, _ := strconv.Atoi(split[1])

		// 	chunk, err := gameWorld.GetChunk(vec.Vec3I{X: x, Y: y})
		// 	if err != nil {
		// 		logger.Error(err)
		// 		continue
		// 	}
		// 	go tickChunk(chunk, finChan)
		// 	<-finChan
		// }

		taken := time.Since(startedTick)
		logger.Infof("Took %s to tick", taken.String())
		<-ticker.C
	}
}

func tickLayer(layer *world.Layer, done chan bool) {
	done <- true
}
