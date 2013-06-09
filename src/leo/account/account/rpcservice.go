/* this is client req service */

package account

import pblib "code.google.com/p/goprotobuf/proto"

import (
	"leo/base"
	"leo/common"
	"leo/proto"
)

type ClientReqService struct {
	owner *ClientReqDispatcher
}

type ClientReqDispatcher struct {
	handlermap map[int] ClientReqHandler
	service *ClientReqService
}

func NewClientReqDispatcher(port *base.Port) (sv *ClientReqDispatcher, err error) {
	sv = new(ClientReqDispatcher)
	err = sv.init(port)
	return
}

func (d *ClientReqDispatcher) init(port *base.Port) error {
	d.handlermap = make(map[int] ClientReqHandler)
	d.Register(proto.REGISTER, new(RegisterHandler))
	d.Register(proto.LOGIN, new(LoginHandler))

	d.service = new(ClientReqService)
	d.service.owner = d
	port.RegisterService(d.service)
	return nil
}

func (sv *ClientReqService) Request(reply *common.RpcClientRequest, v *int) error {
	*v = 0
	op := int(reply.Pkt.Op)
	h, ok := sv.owner.handlermap[op]
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
		return nil
	} else {
// 		ch := make(chan error)
// 		go h.Handle(reply, ch)
// 		return <-ch
		h.Handle(reply)
		//go h.Handle(reply)
		return nil
	}
}

func (d *ClientReqDispatcher) Register(op int, h ClientReqHandler) {
	d.handlermap[op] = h
}