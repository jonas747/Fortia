package ticker

import (
	"github.com/jonas747/fortia/db"
	//ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/log"
	"github.com/jonas747/fortia/vec"
	"github.com/jonas747/fortia/world"
	"time"
)

var (
	logger    *log.LogClient
	authDb    *db.AuthDB
	gameDb    *db.GameDB
	gameWorld *world.World
)

func Run(l *log.LogClient, adb *db.AuthDB, gdb *db.GameDB, gw *world.World, addr string) {
	logger.Info("Running world ticker")
	logger = l
	authDb = adb
	gameDb = gdb
	gameWorld = gw
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
