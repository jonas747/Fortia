package log

import (
	//"encoding/json"
	"fmt"
	ferr "github.com/jonas747/fortia/error"
	"strings"
)

/*
Log levels:
-1: debug
 0: Info
 1: Warning
 2: Error
 3: Fatal
*/

type LogMsg struct {
	Lvl  int
	Msg  string
	Data map[string]interface{}
}

func NewLogMsg(lvl int, msg string, data map[string]interface{}) LogMsg {
	return LogMsg{
		Lvl:  lvl,
		Msg:  msg,
		Data: data,
	}
}

type LogClient struct {
	// TODO
	// connection to server
	// etc
	PrintLvl int
}

func NewLogClient(addr string, Plvl int) (*LogClient, error) {
	return &LogClient{
		PrintLvl: Plvl,
	}, nil
}

// Base log function
func (l *LogClient) Log(msg LogMsg) {
	// Print if were higher than printlvl
	if l.PrintLvl <= msg.Lvl {
		fmt.Println(msg.Msg)
	}

	// Send to logserver
	if msg.Lvl < 0 {
		return
	}
	//TODO
}

// lvl -1 (not recorded or sent to master(maybe sent to master))
func (l *LogClient) Debug(msgs ...string) {
	str := strings.Join(msgs, "")
	l.Log(NewLogMsg(-1, str, make(map[string]interface{})))
}

// lvl 0
func (l *LogClient) Info(msg string, data map[string]interface{}) {
	l.Log(NewLogMsg(0, msg, data))
}

// lvl 1
func (l *LogClient) Warn(msg string, data map[string]interface{}) {
	l.Log(NewLogMsg(1, msg, data))
}

// lvl 2
func (l *LogClient) Error(err ferr.FortiaError) {
	msg := err.GetMessage()
	data := err.GetData()

	data["stacktrace"] = err.GetStack()
	data["context"] = err.GetContext()

	l.Log(NewLogMsg(2, msg, data))
}

// lvl 3
func (l *LogClient) Fatal(err ferr.FortiaError) {
	msg := err.GetMessage()
	data := err.GetData()

	data["stacktrace"] = err.GetStack()
	data["context"] = err.GetContext()

	l.Log(NewLogMsg(2, msg, data))
}
