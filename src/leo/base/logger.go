/* this is the log interface
*/
package logger

type Logger interface {
	Init()
	Close()

	Debug()
	Info()
	Warning()
	Error()
	Critical()
}