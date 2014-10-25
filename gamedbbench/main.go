package main

import (
	"fmt"
	"github.com/jonas747/fortia/db"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	logger *log.LogClient
	authDb *db.AuthDB
	gameDb *db.GameDB
	config *Config
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
		l.Error(ferr.Wrap(nErr, ""))
	}

	l.Info("(2/4) Log client init successful! Creating database connection pools...")

	gdb, nErr := db.NewDatabase(config.GameDb)
	if nErr != nil {
		l.Warn("Not connected to database..." + nErr.Error())
	}

	gameDb = &db.GameDB{gdb}

	l.Info("World ticker started sucessfully")

	chunkSize := 10
	size := 1000

	chunkedSetTest(size, chunkSize)
	chunkedGetTest(size, chunkSize)

	//individualGenTest(size)
	//individualGetTest(size)
}

/*
func individualGenTest(size int) {
	logger.Info("Starting set test")
	totalSize := size * size * size
	blocks := genChunk(size)

	now := time.Now()
	// Set it, this is a ton of set requests...
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			for z := 0; z < size; z++ {
				//Flat[x + WIDTH * (y + DEPTH * z)] = Original[x, y, z]
				index := x + size*(y+size*z)
				b := blocks[index]
				err := gameDb.SetBlock(x, y, z, b)
				if err != nil {
					logger.Error(err)
					return
				}
			}
		}

	}

	duration := time.Since(now)
	logger.Infof("Time taken setting %d blocks: %s, Blocks/s: %.2f", totalSize, duration.String(), float64(totalSize)/duration.Seconds())
}

func individualGetTest(size int) {
	logger.Info("Starting get test")
	now := time.Now()
	// Set it, this is a ton of set requests...
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			for z := 0; z < size; z++ {
				_, err := gameDb.GetBlock(x, y, z)
				if err != nil {
					logger.Error(err)
					return
				}
			}
		}

	}
	duration := time.Since(now)
	totalSize := size * size * size
	logger.Infof("Time taken Getting %d blocks: %s, Blocks/s: %.2f", totalSize, duration.String(), float64(totalSize)/duration.Seconds())
}
*/
func chunkedSetTest(size, chunkSize int) {
	logger.Info("Starting Set test")
	RSize := size / chunkSize

	chunks := make(map[string][]int, 0)

	for x := 0; x < RSize; x++ {
		for y := 0; y < RSize; y++ {
			chunks[fmt.Sprintf("%d:%d", x, y)] = genChunk(chunkSize)
		}
	}

	now := time.Now()
	for k, v := range chunks {
		splitKey := strings.Split(k, ":")
		x, _ := strconv.Atoi(splitKey[0])
		y, _ := strconv.Atoi(splitKey[1])
		gameDb.SetChunk(x, y, v)
	}
	duration := time.Since(now)
	totalSize := size * size * size
	logger.Infof("Time taken Setting %d blocks: %s, Blocks/s: %.2f", totalSize, duration.String(), float64(totalSize)/duration.Seconds())

}

func genChunk(size int) []int {
	rand.Seed(int64(time.Now().Nanosecond()))
	totalSize := size * size * size

	// Generate the block array
	blocks := make([]int, totalSize)
	for i := 0; i < totalSize; i++ {
		blocks[i] = rand.Int()
	}
	return blocks
}

func chunkedGetTest(size, chunkSize int) {
	logger.Info("Starting Get test")
	RSize := size / chunkSize
	now := time.Now()
	for x := 0; x < RSize; x++ {
		for y := 0; y < RSize; y++ {
			gameDb.GetChunk(x, y)
		}
	}
	duration := time.Since(now)
	totalSize := size * size * size
	logger.Infof("Time taken Getting %d blocks: %s, Blocks/s: %.2f", totalSize, duration.String(), float64(totalSize)/duration.Seconds())
}
