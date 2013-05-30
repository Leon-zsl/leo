/* this is the main logic
*/

package gate

import (
	"fmt"
	"runtime/debug"
	"leo/base"
)

type EchoHandler struct {
	ssn *Session
}

func NewEchoHandler(ssn *Session) *EchoHandler {
	fmt.Println("new echo handler")
	h := new(EchoHandler)
	h.ssn = ssn
	ssn.RegisterHandler(h)
	return h
}

func (h *EchoHandler) HandleSessionMsg(ssn *Session, pkt *base.Packet) {
	fmt.Println(pkt.Op(), string(pkt.Args()))
	ssn.Send(pkt)
}

func (h *EchoHandler) HandleSessionError(ssn *Session, err error) {
	fmt.Println("session error:", ssn.Addr(), ",", err.Error())
	debug.PrintStack()
}

func (h *EchoHandler) HandleSessionClose(ssn *Session) {
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
	//donothing
	return nil
}

func (srv *GateService) Start() {
	Root.SessionMgr.RegisterNewSessionListener(srv)
}

func (srv *GateService) Close() {
	//donothing
}

func (srv *GateService) Update() {
	//donothing
}

func (srv *GateService) HandleNewSession(ssn *Session) {
	fmt.Println("new session: ", ssn.Addr())
	NewEchoHandler(ssn)
}