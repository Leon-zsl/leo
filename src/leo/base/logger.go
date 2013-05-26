/* this is the log interface
*/
package logger

import (
	"log4go"
	)

type Logger interface {
	Init()
	Close()

	Debug()
	Info()
	Warning()
	Error()
	Critical()
}

type Logger4Go struct {
	Level int
	
}

func CreateLogger() *Logger {
	return nil
}

func (lg *Logger) Init() {
	//todo
}

func (lg *Logger)Close() {
	//todo
}

func (lg *Logger) Debug() {
	//todo:
}

func (lg *Logger) Info() {
	//todo
}

func (lg *Logger) Warning() {
	//todo
}

func (lg *Logger) Error() {
	//todo
}

func (lg *Logger) Critical() {
	//todo
}