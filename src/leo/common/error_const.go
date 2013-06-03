/* this is common
*/

package common

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
