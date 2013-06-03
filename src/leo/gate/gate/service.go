/* this is the main logic
*/

package gate

import (
	"fmt"
	"strconv"
	"path"
	"ini"
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
	Clock *base.Clock
}

func NewGateService() (service *GateService, err error) {
	service = new(GateService)
	err = service.init()
	return
}

func (service *GateService) init() error {
	service.Clock, _ = base.NewClock()
	Root.Acceptor.RegisterAcceptedSessionListener(service)
	return nil
}

func (service *GateService) Start() error {
	service.Clock.Start()
	service.connect_master()
	return nil
}

func (service *GateService) Close() error {
	service.disconnect_master()
	service.Clock.Close()
	return nil
}

func (service *GateService) Tick() error {
	service.Clock.Tick()
	//fmt.Println("gate service tick")
	return nil
}

func (service *GateService) Save() error {
	return nil
}

func (service *GateService) HandleAcceptedSession(ssn *base.Session) {
	fmt.Println("accept session: ", ssn.Addr())
	NewEchoHandler(ssn)
}

func (service *GateService) connect_master() error {
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
	Root.Port.OpenConnect(port_id, ip, port)

	service.master_port_id = port_id
	return nil
}

func (service *GateService) disconnect_master() error {
	fmt.Println("disconnect_master")
	Root.Port.CloseConnect(service.master_port_id)
	return nil
}