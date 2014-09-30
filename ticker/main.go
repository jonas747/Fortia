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
}

func run() {

}
