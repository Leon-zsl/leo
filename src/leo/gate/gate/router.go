/* this is router */

package gate

import (
	//"fmt"
	"runtime/debug"

	"leo/base"
	"leo/proto"
	"leo/common"
)

type Router struct {
	mgr *RouterMgr
	ssn *base.Session
	world_port_id int
}

func NewRouter(mgr *RouterMgr, ssn *base.Session) (router *Router, err error) {
	router = new(Router)
	err = router.init(mgr, ssn)
	return
}

func (router *Router) init(mgr *RouterMgr, ssn *base.Session) error {
	router.ssn = ssn
	router.mgr = mgr
	return nil
}

func (router *Router) WorldServer() int {
	return router.world_port_id
}

func (router *Router) SetWorldServer(port int) {
	router.world_port_id = port
}

func (router *Router) Session() *base.Session {
	return router.ssn
}

func (router *Router) HandleSessionStart(ssn *base.Session) {
	base.LoggerIns.Info("session start:", ssn.SID(), ssn.Addr())
}

func (router *Router) HandleSessionMsg(ssn *base.Session, pkt *base.Packet) {
	target := proto.RouteMap[int(pkt.Op())]
	req := &common.RpcClientRequest{ssn.SID(), pkt}
	switch target {
	case "master":
		Root.Port.SendAsync(ServiceIns.MasterServer(),"ClientReqService.Request", req)
	case "account":
		Root.Port.SendAsync(ServiceIns.AccountServer(),"ClientReqService.Request", req)
	case "world":
		Root.Port.SendAsync(router.world_port_id, "ClientReqService.Request", pkt)
	default:
		base.LoggerIns.Error("target server is nil", pkt.Op())
	}
}

func (router *Router) HandleSessionError(ssn *base.Session, err error) {
	base.LoggerIns.Error("session error:", ssn.SID(), ssn.Addr(), err)
	debug.PrintStack()
}

func (router *Router) HandleSessionClose(ssn *base.Session) {
	base.LoggerIns.Info("session close:", ssn.SID(), ssn.Addr())
	router.mgr.DelRouter(ssn.SID())
}