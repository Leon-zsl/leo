/* this is the log implemention by log4go
*/

package logger

import (
	"log4go"
	)

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