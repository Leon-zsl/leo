/* this is rpc handlers
*/

package gate

import (
//	"fmt"
	"leo/common"
)

type RpcService struct {
	//do nothing
}

func NewRpcService() (sv *RpcService, err error) {
	sv = new(RpcService)
	err = nil
	return
}

func (rs *RpcService) SendTo(reply *common.RpcSendTo, v *int) error {
	rt, ok := ServiceIns.RouterMgr.Router(reply.Sid)
	if ok {
		rt.Session().Send(reply.Pkt)
	}
	*v = 0
	return nil
}

//it will hurt the response speed seriously
func (rs *RpcService) SendToAll(reply *common.RpcSendToAll, v *int) error {
	for _, v := range(ServiceIns.RouterMgr.Routers()) {
		v.Session().Send(reply.Pkt)
	}
	*v = 0
	return nil
}

func (rs *RpcService) Broadcast(reply *common.RpcBroadcast, v *int) error {
	for _, v := range(reply.Sids) {
		rt, ok := ServiceIns.RouterMgr.Router(v)
		if ok {
			rt.Session().Send(reply.Pkt)
		}
	}
	*v = 0
	return nil
}

func (rs *RpcService) MoveWorld(reply *common.RpcMoveWorld, v *int) error {
	rt, ok := ServiceIns.RouterMgr.Router(reply.Sid)
	if ok {
		rt.SetWorldServer(reply.PortID)
	}
	*v = 0
	return nil
}