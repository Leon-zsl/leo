/* this is the main logic
*/

package gate

import (
	"fmt"
	"runtime/debug"

	"leo/base"
)

type EchoHandler struct {
	ssn *base.Session
}

func NewEchoHandler(ssn *base.Session) *EchoHandler {
	fmt.Println("new echo handler")
	h := new(EchoHandler)
	h.ssn = ssn
	ssn.RegisterHandler(h)
	return h
}

func (h *EchoHandler) HandleSessionStart(ssn *base.Session) {
	fmt.Println("session start:", ssn.Addr())
}

func (h *EchoHandler) HandleSessionMsg(ssn *base.Session, pkt *base.Packet) {
	fmt.Println(pkt.Op(), string(pkt.Args()))
	ssn.Send(pkt)
}

func (h *EchoHandler) HandleSessionError(ssn *base.Session, err error) {
	fmt.Println("session error:", ssn.Addr(), ",", err.Error())
	debug.PrintStack()
}

func (h *EchoHandler) HandleSessionClose(ssn *base.Session) {
	fmt.Println("session close:", ssn.Addr())
}

type GateService struct {
}

func NewService() (service *GateService, err error) {
	service = new(GateService)
	err = service.init()
	return
}

func (srv *GateService) init() error {
	Root.Acceptor.RegisterAcceptedSessionListener(srv)
	return nil
}

func (srv *GateService) Start() {
	//do nothing
}

func (srv *GateService) Close() {
	//do nothing
}

func (srv *GateService) Tick() {
	//do nothing
}

func (srv *GateService) HandleAcceptedSession(ssn *base.Session) {
	fmt.Println("accept session: ", ssn.Addr())
	NewEchoHandler(ssn)
}