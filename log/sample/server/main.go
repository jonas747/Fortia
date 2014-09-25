package main

import (
	"fmt"
	"github.com/jonas747/fortia/log"
)

func main() {
	initServer()
}

func initServer() {
	logServer, err := log.NewServer(":8080")
	if err != nil {
		fmt.Println("Error launching logsever:", err)
		return
	}
	for {
		msg := <-logServer.Out
		fmt.Print(msg.StringDetailed(true))
	}
}
