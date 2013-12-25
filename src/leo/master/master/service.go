/* this is master service
 */

package master

import (
	"fmt"
	"path"
	"strconv"
	//	"ini"
	"leo/base"
)

import ini "github.com/vaughan0/go-ini"

type MasterService struct {
	db_port_id   int
	gate_port_id int
	Clock        *base.Clock
}

func NewMasterService() (service *MasterService, err error) {
	service = new(MasterService)
	err = service.init()
	return
}

func (service *MasterService) init() error {
	service.Clock, _ = base.NewClock()
	return nil
}

func (service *MasterService) Start() error {
	service.Clock.Start()
	return nil
}

func (service *MasterService) Close() error {
	service.Clock.Close()
	return nil
}

func (service *MasterService) Tick() error {
	service.Clock.Tick()
	//fmt.Println("world service tick")
	return nil
}

func (service *MasterService) Save() error {
	return nil
}

func (service *MasterService) DBServer() int {
	return service.db_port_id
}

func (service *MasterService) GateServer() int {
	return service.gate_port_id
}

func (service *MasterService) connect_db() error {
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

func (service *MasterService) disconnect_db() error {
	fmt.Println("disconnect db")
	Root.Port.CloseConnect(service.db_port_id)
	return nil
}

func (service *MasterService) get_gate() error {
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
