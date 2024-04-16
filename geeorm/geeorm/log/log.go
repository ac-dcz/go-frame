package log

import (
	"io"
	stlog "log"
	"os"
	"sync"
	"sync/atomic"
)

const (
	BLUE  = "\033[34m"
	RED   = "\033[31m"
	RESET = "\033[0m"
)

type LogType int

const (
	INFO LogType = iota
	ERROR
)

const (
	Disabled int32 = iota
	InfoLevel
	ErrorLevel
)

var (
	infoLog  = stlog.New(os.Stdout, BLUE+"[INFO ]"+RESET, stlog.LstdFlags|stlog.Lshortfile)
	errorLog = stlog.New(os.Stderr, RED+"[ERROR]"+RESET, stlog.LstdFlags|stlog.Lshortfile)
	mu       = sync.Mutex{}
	logmap   = make(map[LogType]*stlog.Logger)
	level    atomic.Int32
)

func init() {
	logmap[INFO] = infoLog
	logmap[ERROR] = errorLog
	level.Store(ErrorLevel)
}

func SetOutPut(typ LogType, w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	logger := logmap[typ]
	logger.SetOutput(w)
}

var (
	Info   = infoLog.Print
	Infof  = infoLog.Printf
	Error  = errorLog.Print
	Errorf = errorLog.Printf
)
