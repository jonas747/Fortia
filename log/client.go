package log

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	ferr "github.com/jonas747/fortia/error"
	"net"
	"strings"
)

type LogClient struct {
	Conn     net.Conn
	PrintLvl int
	Host     string
}

func NewLogClient(addr string, Plvl int, host string) (*LogClient, error) {
	// Dial the address
	conn, err := net.Dial("tcp", addr)
	fmt.Println(conn.LocalAddr())
	fmt.Println(conn.RemoteAddr())
	return &LogClient{
		PrintLvl: Plvl,
		Conn:     conn,
		Host:     host,
	}, err
}

// Base log function
func (l *LogClient) Log(msg LogMsg) {
	// Print if were higher than printlvl
	if l.PrintLvl <= msg.Lvl {
		fmt.Print(msg.String())
	}

	if msg.Lvl < 0 || l.Conn == nil {
		return
	}

	msg.Host = l.Host

	// Send to logserver
	serialized, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	// First 8 bytes are the length of the log message
	var buffer bytes.Buffer
	length := int64(len(serialized))
	binary.Write(&buffer, binary.LittleEndian, length)
	buffer.Write(serialized)
	raw := buffer.Bytes()
	_, err = l.Conn.Write(raw)
	if err != nil {
		fmt.Println("Error sending log: ", err)
	}
}

// lvl -1 (not recorded or sent to master(maybe sent to master))
func (l *LogClient) Debug(msgs ...string) {
	str := strings.Join(msgs, "")
	l.Log(NewLogMsg(-1, str, make(map[string]interface{})))
}

// lvl 0
func (l *LogClient) Info(msg string) {
	l.Infoa(msg, make(map[string]interface{}))
}

func (l *LogClient) Infoa(msg string, data map[string]interface{}) {
	l.Log(NewLogMsg(0, msg, data))
}

// lvl 1
func (l *LogClient) Warn(msg string) {
	l.Warna(msg, make(map[string]interface{}))
}

func (l *LogClient) Warna(msg string, data map[string]interface{}) {
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

	l.Log(NewLogMsg(3, msg, data))
}
