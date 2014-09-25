package log

import (
	"fmt"
	"time"
)

const (
	col_nocolor = "\x1b[0m"
	col_black   = "\x1b[30m"
	col_red     = "\x1b[31m"
	col_green   = "\x1b[32m"
	col_yellow  = "\x1b[33m"
	col_blue    = "\x1b[34m"
	col_magenta = "\x1b[35m"
	col_cyan    = "\x1b[36m"
	col_white   = "\x1b[37m"
)

var (
	col_bold_black   = col_black[:4] + ";1m"
	col_bold_red     = col_red[:4] + ";1m"
	col_bold_green   = col_green[:4] + ";1m"
	col_bold_yellow  = col_yellow[:4] + ";1m"
	col_bold_blue    = col_blue[:4] + ";1m"
	col_bold_magenta = col_magenta[:4] + ";1m"
	col_bold_cyan    = col_cyan[:4] + ";1m"
	col_bold_white   = col_white[:4] + ";1m"
)

/*
Log levels:
-1: debug
 0: Info
 1: Warning
 2: Error
 3: Fatal
*/

var LogLvlStr = map[int]string{
	-1: "Debug",
	0:  "Info",
	1:  "Warning",
	2:  "Error",
	3:  "Fatal",
}

var LogColors = map[int]string{
	-1: col_cyan,
	0:  col_white,
	1:  col_yellow,
	2:  col_magenta,
	3:  col_red,
}

// Represents a log message
type LogMsg struct {
	Lvl       int
	Msg       string
	Data      map[string]interface{}
	Host      string
	Timestamp time.Time
}

func NewLogMsg(lvl int, msg string, data map[string]interface{}) LogMsg {
	now := time.Now()
	return LogMsg{
		Lvl:       lvl,
		Msg:       msg,
		Data:      data,
		Timestamp: now,
	}
}

func (l *LogMsg) String() string {
	return l.StringDetailed(false)
}

func (l *LogMsg) StringDetailed(host bool) string {
	timeStr := l.Timestamp.Format(time.Stamp)
	stackTrace, ok := l.Data["stacktrace"]
	str := ""
	if host {
		str = fmt.Sprintf("{%s}", l.Host)
	}
	if ok {
		str += fmt.Sprintf("%s[%s] %s: %s\n%s\x1b[0m", LogColors[l.Lvl], timeStr, LogLvlStr[l.Lvl], l.Msg, stackTrace)
	} else {
		str += fmt.Sprintf("%s[%s] %s: %s\x1b[0m\n", LogColors[l.Lvl], timeStr, LogLvlStr[l.Lvl], l.Msg)
	}
	return str
}
