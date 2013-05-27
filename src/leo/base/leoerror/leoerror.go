/* this is the error code
*/
package leoerror

import (
	"strconv"
)
const (
	ErrNo = iota
	ErrSys
	ErrStartFailed
	ErrRuntimeExcept
)

var errMap = map[int] string {
ErrNo : "no error",
ErrSys : "sys error",
ErrStartFailed : "start failed",
ErrRuntimeExcept : "runtime exception",
}

type LeoError struct {
	code int
	msg string
}

func CreateLeoError(errcode int, errmsg string) LeoError {
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

