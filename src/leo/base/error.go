/* this is the error code
*/
package base

import (
	"strconv"
)

const (
	LeoErrNo = iota
	LeoErrSys
	LeoErrStartFailed
	LeoErrRuntimeExcept
)

var errMap = map[int] string {
LeoErrNo : "no error",
LeoErrSys : "sys error",
LeoErrStartFailed : "start failed",
LeoErrRuntimeExcept : "runtime exception",
}

type LeoError struct {
	code int
	msg string
}

func NewLeoError(errcode int, errmsg string) LeoError {
	return LeoError{code:errcode, msg:errmsg}
}

func (err LeoError) Error() string {
	val, ok := errMap[err.code]
	if !ok {
		val = strconv.Itoa(err.code)
	}
	return "[code]" + val + "[msg]" + err.msg
}

func (err LeoError) Code() int {
	return err.code
}

func (err LeoError) Msg() string {
	return err.msg
}

