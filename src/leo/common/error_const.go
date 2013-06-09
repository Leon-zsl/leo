/* this is common
*/

package common

const (
	LeoErrNo = iota
	LeoErrSys
	LeoErrStartFailed
	LeoErrRuntimeExcept

	//rpc
	LeoErrRpcNoArg
	LeoErrRpcIllegalArg
	LeoErrRpcIllegalResp
)

var ErrMap = map[int] string {
LeoErrNo : "no error",
LeoErrSys : "sys error",
LeoErrStartFailed : "start failed",
LeoErrRuntimeExcept : "runtime exception",

LeoErrRpcNoArg : "rpc invoker no arg",
LeoErrRpcIllegalArg : "rpc invoker illeagal arg",
LeoErrRpcIllegalResp : "rpc invoker illeagal arg",
}