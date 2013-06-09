/* this is router */

package gate

import pblib "code.google.com/p/goprotobuf/proto"
import (
	"fmt"
	"io"
	"errors"
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
	target, ok := proto.RouteMap[int(pkt.Op)]
	if !ok {
		router.send_err_resp(ssn, pkt.Op, proto.EC_UNKNOWN_OP, "")
		return
	}

	req := &common.RpcClientRequest{ssn.SID(), pkt}
	var err error = nil
	switch target {
	case "master":
		err = Root.Port.SendAsync(ServiceIns.MasterServer(),"ClientReqService.Request", req)
	case "account":
		err = Root.Port.SendAsync(ServiceIns.AccountServer(),"ClientReqService.Request", req)
	case "world":
		err = Root.Port.SendAsync(router.world_port_id, "ClientReqService.Request", req)
	default:
		err = errors.New("unknown target service")
	}
	
	if err != nil {
		base.LoggerIns.Error("router err", pkt.Op, err.Error())
		router.send_err_resp(ssn, pkt.Op, proto.EC_UNKNOWN_OP, err.Error())
	}
}

func (router *Router) HandleSessionError(ssn *base.Session, err error) {
	//means the remote close the sock
	if err == io.EOF {
		fmt.Println("the remote close the sock")
		return
	}
	base.LoggerIns.Error("session error:", ssn.SID(), ssn.Addr(), err)
	debug.PrintStack()
}

func (router *Router) HandleSessionClose(ssn *base.Session) {
	base.LoggerIns.Info("session close:", ssn.SID(), ssn.Addr())
	router.mgr.DelRouter(ssn.SID())
}

func (router *Router) send_err_resp(ssn *base.Session, op int32, code int32, msg string) {
	pb := &proto.Error{ Op : pblib.Int32(op), ErrorCode : pblib.Int32(code), ErrorMsg : pblib.String(msg)}
	val, _ := pblib.Marshal(pb)
	pkt := base.NewPacket(proto.ERROR, val)
	ssn.Send(pkt)
}