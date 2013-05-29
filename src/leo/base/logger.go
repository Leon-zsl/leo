/* this is logger
*/

package base

import (
	"os"
	"runtime"
	"log"
	"log4go"
)

type Logger struct {
	lgtype int
	logger4go log4go.Logger
	loggersys *log.Logger
}

const (
	LOG_TYPE_SYS = iota
	LOG_TYPE_LOG4GO
)


func NewLogger(ty int, confile string) (lg *Logger, err error) {
	lg = new(Logger)
	err = lg.init(ty, confile)
	return
}

func (lg *Logger) init(ty int, confile string) error {
	if ty == LOG_TYPE_SYS {
		lg.loggersys = log.New(os.Stderr, "", log.Lshortfile)
		return nil
	} else {
		log := make(log4go.Logger)
		log.LoadConfiguration(confile)
		lg.logger4go = log
		return nil
	}
}

func (lg *Logger) Close() {
	if lg.logger4go != nil {
		lg.logger4go.Close()
	}
}

func (lg *Logger) Critical(v ...interface{}) {
	_, f, l, _ := runtime.Caller(1)
	if lg.lgtype == LOG_TYPE_LOG4GO {
		lg.logger4go.Critical("[", f, "]", "[", l, "]", v)
	} else {
		_, f, l, _ := runtime.Caller(1)
		lg.loggersys.Println("[", f, "]", "[", l, "]", v)
	}
}

func (lg *Logger) Error(v ...interface{}) {
		_, f, l, _ := runtime.Caller(1)
	if lg.lgtype == LOG_TYPE_LOG4GO {
		lg.logger4go.Error("[", f, "]", "[", l, "]", v)
	} else {
		lg.loggersys.Println("[", f, "]", "[", l, "]", v)
	}
}

func (lg *Logger) Warn(v ...interface{}) {
	_, f, l, _ := runtime.Caller(1)
	if lg.lgtype == LOG_TYPE_LOG4GO {
		lg.logger4go.Warn("[", f, "]", "[", l, "]", v)
	} else {
		_, f, l, _ := runtime.Caller(1)
		lg.loggersys.Println("[", f, "]", "[", l, "]", v)
	}
}

func (lg *Logger) Info(v ...interface{}) {
	_, f, l, _ := runtime.Caller(1)
	if lg.lgtype == LOG_TYPE_LOG4GO {
		lg.logger4go.Info("[", f, "]", "[", l, "]", v)
	} else {
		_, f, l, _ := runtime.Caller(1)
		lg.loggersys.Println("[", f, "]", "[", l, "]", v)
	}
}

func (lg *Logger) Debug(v ...interface{}) {
	_, f, l, _ := runtime.Caller(1)
	if lg.lgtype == LOG_TYPE_LOG4GO {
		lg.logger4go.Debug("[", f, "]", "[", l, "]", v)
	} else {
		_, f, l, _ := runtime.Caller(1)
		lg.loggersys.Println("[", f, "]", "[", l, "]", v)
	}
}