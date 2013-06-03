/* this is the main logic
*/

package gate

import (
	"fmt"
	"strconv"
	"path"
	"ini"
	"time"
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
	master_port_id int
	connect bool
}

func NewGateService() (service *GateService, err error) {
	service = new(GateService)
	err = service.init()
	return
}

func (srv *GateService) init() error {
	Root.Acceptor.RegisterAcceptedSessionListener(srv)
	return nil
}

func (srv *GateService) Start() error {
	srv.connect_master()
	return nil
}

func (srv *GateService) Close() error {
	return nil
}

func (srv *GateService) Tick() error {
	if !srv.connect {
		time.Sleep(1e9)
		srv.connect_master()
	} else {
		time.Sleep(1e9)
		srv.disconnect_master()
	}
	return nil
}

func (srv *GateService) Save() error {
	return nil
}

func (srv *GateService) HandleAcceptedSession(ssn *base.Session) {
	fmt.Println("accept session: ", ssn.Addr())
	NewEchoHandler(ssn)
}

func (srv *GateService) connect_master() error {
	fmt.Println("connect master")

	//parse config file
	confile := path.Join(CONF_PATH, CONF_FILE)
	conf, err := ini.LoadFile(confile)
	if err != nil {
		return err
	}

	id, _ := conf.Get("master", "id")
	ip, _ := conf.Get("master", "ip")
	pt, _ := conf.Get("master", "port")
	port_id, _ := strconv.Atoi(id)
	port, _ := strconv.Atoi(pt)
	fmt.Println("connect to master:", port_id, ip, port)
	Root.Port.OpenConnect(port_id, ip, port)

	srv.master_port_id = port_id
	srv.connect = true
	return nil
}

func (srv *GateService) disconnect_master() error {
	fmt.Println("disconnect_master")
	Root.Port.CloseConnect(srv.master_port_id)
	srv.connect = false
	return nil
}