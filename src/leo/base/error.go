/* this is the error code
*/
package base

import (
)

type LeoError struct {
	code string
	msg string
}

func NewLeoError(errcode, errmsg string) LeoError {
	return LeoError{code:errcode, msg:errmsg}
}

func (err LeoError) Error() string {
	return "[code]" + err.code + "[msg]" + err.msg
}

func (err LeoError) Code() string {
	return err.code
}

func (err LeoError) Msg() string {
	return err.msg
}

