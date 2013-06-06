/* this is handler interface */

package account

import (
	"leo/common"
)

type ClientReqHandler interface {
	Handle(replay *common.RpcClientRequest)
}

type RegisterHandler struct {
}

func (h *RegisterHandler) Handle(reply *common.RpcClientRequest) {
	//todo
}

type LoginHandler struct {
}

func (h *LoginHandler) Handle(reply *common.RpcClientRequest) {
	//todo:
}