package logger

import (
	"fmt"
	"log"
	"os"
)

type LOGGER_LEVEL string

const (
	DEBUG_LEVEL LOGGER_LEVEL = "debug"
	INFO_LEVEL  LOGGER_LEVEL = "info"
)

type logger struct {
	debug *log.Logger
	info  *log.Logger
	lvl   LOGGER_LEVEL
}

var lg *logger

func SetupLogger(prefix string, level LOGGER_LEVEL) {
	flags := log.Lmicroseconds | log.Ldate
	debug := log.New(os.Stdout, fmt.Sprintf("[%s] [DEBUG] ", prefix), flags)
	info := log.New(os.Stdout, fmt.Sprintf("[%s] [INFO ] ", prefix), flags)
	lg = &logger{
		debug: debug,
		info:  info,
		lvl:   level,
	}
}

func Debugf(format string, v ...interface{}) {
	if lg.lvl == INFO_LEVEL {
		return
	}
	output(lg.debug, fmt.Sprintf(format, v...))
}

func Debug(v ...interface{}) {
	if lg.lvl == INFO_LEVEL {
		return
	}
	output(lg.debug, fmt.Sprintln(v...))
}

func Infof(format string, v ...interface{}) {
	output(lg.info, fmt.Sprintf(format, v...))

}

func Info(v ...interface{}) {
	output(lg.info, fmt.Sprintln(v...))
}

func output(lg *log.Logger, s string) {
	lg.Output(3, s)
}
