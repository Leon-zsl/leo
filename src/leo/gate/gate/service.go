/* this is the main logic
*/

package gate

import (
	"fmt"
	"strconv"
	"path"
	"ini"

	"leo/base"
)

type GateService struct {
	master_port_id int
	account_port_id int

	Clock *base.Clock
	RouterMgr *RouterMgr
}

var (
	ServiceIns *GateService = nil
)

func NewGateService() (service *GateService, err error) {
	if ServiceIns != nil {
		service = ServiceIns
		err = nil
	} else {
		service = new(GateService)
		err = service.init()
		ServiceIns = service
	}
	return
}

func (service *GateService) init() error {
	service.Clock, _ = base.NewClock()

	mgr, err := NewRouterMgr()
	if err != nil {
		return err
	}
	Root.Acceptor.RegisterAcceptedSessionListener(mgr)
	service.RouterMgr = mgr

	return nil
}

func (service *GateService) Start() error {
	service.Clock.Start()
	service.connect_master()
	return nil
}

func (service *GateService) Close() error {
	service.disconnect_master()
	Root.Acceptor.UnRegisterAcceptedSessionListener(service.RouterMgr)
	service.RouterMgr.Close()
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

func (service *GateService) MasterServer() int {
	return service.master_port_id
}

func (service *GateService) AccountServer() int {
	return service.account_port_id
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

func (service *GateService) connect_account() error {
	fmt.Println("connect account")

	//parse config file
	confile := path.Join(CONF_PATH, CONF_FILE)
	conf, err := ini.LoadFile(confile)
	if err != nil {
		return err
	}

	id, _ := conf.Get("account", "id")
	ip, _ := conf.Get("account", "ip")
	pt, _ := conf.Get("account", "port")
	port_id, _ := strconv.Atoi(id)
	port, _ := strconv.Atoi(pt)
	Root.Port.OpenConnect(port_id, ip, port)

	service.account_port_id = port_id
	return nil
}

func (service *GateService) disconnect_account() error {
	fmt.Println("disconnect_account")
	Root.Port.CloseConnect(service.account_port_id)
	return nil
}