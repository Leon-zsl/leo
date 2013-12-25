/* this is the error code
 */
package base

import (
	"strconv"
)

type LeoError struct {
	code int
	msg  string
}

func NewLeoError(errcode int, errmsg string) LeoError {
	return LeoError{code: errcode, msg: errmsg}
}

func (err LeoError) Error() string {
	return "[code]" + strconv.Itoa(err.code) + "[msg]" + err.msg
}

func (err LeoError) Code() int {
	return err.code
}

func (err LeoError) Msg() string {
	return err.msg
}
