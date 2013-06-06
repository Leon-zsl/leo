/* this is account service
*/

package account

import (
	"fmt"
	"path"
	"ini"
	"strconv"

	"leo/base"
)

type AccountService struct {
	master_port_id int
	db_port_id int
	gate_port_id int
	Clock *base.Clock
}

var (
	AccountServiceIns *AccountService = nil
)

func NewAccountService() (service *AccountService, err error) {
	if AccountServiceIns != nil {
		service = AccountServiceIns
		err = nil
	} else {
		service = new(AccountService)
		err = service.init()
	}
	return
}

func (service *AccountService) init() error {
	service.Clock, _ = base.NewClock()
	return nil
}

func (service *AccountService) Start() error {
	service.Clock.Start()
	err := service.connect_master()
	if err != nil {
		return err
	}
	err = service.connect_db()
	if err != nil {
		return err
	}
	err = service.get_gate()
	if err != nil {
		return err
	}
	return nil
}

func (service *AccountService) Close() error {
	service.disconnect_db()
	service.disconnect_master()
	service.Clock.Close()
	return nil
}

func (service *AccountService) Tick() error {
	service.Clock.Tick()
	//fmt.Println("account service tick")
	return nil
}

func (service *AccountService) Save() error {
	return nil
}

func (service *AccountService) MasterServer() int {
	return service.master_port_id
}

func (service *AccountService) DBServer() int {
	return service.db_port_id
}

func (service *AccountService) GateServer() int {
	return service.gate_port_id
}

func (service *AccountService) connect_master() error {
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

func (service *AccountService) disconnect_master() error {
	fmt.Println("disconnect_master")
	Root.Port.CloseConnect(service.master_port_id)
	return nil
}

func (service *AccountService) connect_db() error {
	fmt.Println("connect db")

	//parse config file
	confile := path.Join(CONF_PATH, CONF_FILE)
	conf, err := ini.LoadFile(confile)
	if err != nil {
		return err
	}

	id, _ := conf.Get("db", "id")
	ip, _ := conf.Get("db", "ip")
	pt, _ := conf.Get("db", "port")
	port_id, _ := strconv.Atoi(id)
	port, _ := strconv.Atoi(pt)
	Root.Port.OpenConnect(port_id, ip, port)

	service.db_port_id = port_id
	return nil
}

func (service *AccountService) disconnect_db() error {
	fmt.Println("disconnect db")
	Root.Port.CloseConnect(service.db_port_id)
	return nil
}

func (service *AccountService) get_gate() error {
	//parse config file
	confile := path.Join(CONF_PATH, CONF_FILE)
	conf, err := ini.LoadFile(confile)
	if err != nil {
		return err
	}

	id, _ := conf.Get("gate", "id")
	port_id, _ := strconv.Atoi(id)
	service.gate_port_id = port_id
	return nil
}