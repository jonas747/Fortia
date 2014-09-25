package main

import (
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/log"
	"strconv"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		go client("localhost:8080", strconv.Itoa(i))
	}
	client("localhost:8080", "end")
	time.Sleep(time.Duration(1) * time.Second)
}

func client(addr string, host string) {
	logClient, err := log.NewLogClient(addr, -1, host)
	if err != nil {
		logClient.Error(ferr.New("Error Connecting to logserver"))
	}

	logClient.Debug("Hello debug")
	logClient.Info("Testing info")

	logClient.Warn("Woops warning")

	ferror := ferr.New("Someting went wrong")
	fatal := ferr.New("Something went terribly wrong")
	logClient.Error(ferror)
	logClient.Fatal(fatal)
	logClient.Info("Testing info")
}
