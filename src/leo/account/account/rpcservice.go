/* this is client req service */

package account

import pblib "code.google.com/p/goprotobuf/proto"

import (
	"leo/base"
	"leo/common"
	"leo/proto"
)

type ClientReqService struct {
	handler_map map[int] ClientReqHandler
}

func NewClientReqService() (sv *ClientReqService, err error) {
	sv = new(ClientReqService)
	err = sv.init()
	return
}

func (sv *ClientReqService) init() error {
	sv.handler_map = make(map[int] ClientReqHandler)
	sv.Register(proto.REGISTER, new(RegisterHandler))
	sv.Register(proto.LOGIN, new(LoginHandler))
	return nil
}

func (sv *ClientReqService) Request(reply *common.RpcClientRequest, v interface{}) error {
	op := int(reply.Pkt.Op())
	h, ok := sv.handler_map[op]
	if !ok {
		pb := &proto.Error{ ErrorCode : pblib.Int32(proto.EC_UNKNOWN_OP), 
			ErrorMsg : pblib.String("")}
		val, err := pblib.Marshal(pb)
		if err != nil {
			base.LoggerIns.Error("pb marshal err:", err)
			return err
		}
		pkt := base.NewPacket(proto.ERROR, val)
		resp := &common.RpcSendTo{reply.Sid, pkt}
		Root.Port.SendAsync(AccountServiceIns.GateServer(), "RpcService.SendTo", resp)
	} else {
		go h.Handle(reply)
	}
	return nil
}

func (sv *ClientReqService) Register(op int, h ClientReqHandler) {
	sv.handler_map[op] = h
}

