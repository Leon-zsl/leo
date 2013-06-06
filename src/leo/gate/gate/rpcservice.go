/* this is rpc handlers
*/

package gate

import (
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

func (rs *RpcService) SendTo(reply *common.RpcSendTo, v interface{}) error {
	rt, ok := ServiceIns.RouterMgr.Router(reply.Sid)
	if ok {
		rt.Session().Send(reply.Pkt)
	}	
	return nil
}

//it will hurt the response speed seriously
func (rs *RpcService) SendToAll(reply *common.RpcSendToAll, v interface{}) error {
	for _, v := range(ServiceIns.RouterMgr.Routers()) {
		v.Session().Send(reply.Pkt)
	}
	return nil
}

func (rs *RpcService) Broadcast(reply *common.RpcBroadcast, v interface{}) error {
	for _, v := range(reply.Sids) {
		rt, ok := ServiceIns.RouterMgr.Router(v)
		if ok {
			rt.Session().Send(reply.Pkt)
		}
	}
	return nil
}

func (rs *RpcService) MoveWorld(reply *common.RpcMoveWorld, v interface{}) error {
	rt, ok := ServiceIns.RouterMgr.Router(reply.Sid)
	if ok {
		rt.SetWorldServer(reply.PortID)
	}
	return nil
}