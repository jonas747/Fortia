package log

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/jonas747/fortia/errors"
	"github.com/jonas747/fortia/messages"
	"net"
)

type LogClient struct {
	Conn     net.Conn
	PrintLvl int
	Host     string
}

func NewLogClient(addr string, Plvl int, host string) (*LogClient, errors.FortiaError) {
	// Dial the address
	conn, err := net.Dial("tcp", addr)
	lc := &LogClient{
		PrintLvl: Plvl,
		Conn:     conn,
		Host:     host,
	}
	if err != nil {
		return lc, errors.Wrap(err, messages.ErrorCode_NetDialErr, "", nil)
	}
	return lc, nil
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

// For log.Logger output
func (l *LogClient) Write(p []byte) (n int, err error) {
	// TODO: Check for prefix
	ferr := errors.New(messages.ErrorCode_UnKnownErr, string(p), nil)
	l.Error(ferr)
	return len(p), nil
}

// lvl -1 (not recorded or sent to master(maybe sent to master))
func (l *LogClient) Debug(msgs ...interface{}) {
	str := fmt.Sprint(msgs...)
	l.Log(NewLogMsg(-1, str, make(map[string]interface{})))
}

func (l *LogClient) Debugf(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	l.Log(NewLogMsg(-1, str, make(map[string]interface{})))
}

// lvl 0
func (l *LogClient) Info(msg ...interface{}) {
	l.Infoa(fmt.Sprint(msg...), make(map[string]interface{}))
}

func (l *LogClient) Infof(format string, args ...interface{}) {
	l.Infoa(fmt.Sprintf(format, args...), make(map[string]interface{}))
}

func (l *LogClient) Infoa(msg string, data map[string]interface{}) {
	l.Log(NewLogMsg(0, msg, data))
}

// lvl 1
func (l *LogClient) Warn(msg ...interface{}) {
	l.Warna(fmt.Sprint(msg...), make(map[string]interface{}))
}

func (l *LogClient) Warnf(format string, args ...interface{}) {
	l.Warna(fmt.Sprintf(format, args...), make(map[string]interface{}))
}

func (l *LogClient) Warna(msg string, data map[string]interface{}) {
	l.Log(NewLogMsg(1, msg, data))
}

// lvl 2
func (l *LogClient) Error(err errors.FortiaError) {
	msg := LogMsgFromError(2, err)
	l.Log(msg)
}

// lvl 3
func (l *LogClient) Fatal(err errors.FortiaError) {
	msg := LogMsgFromError(3, err)
	l.Log(msg)
}
