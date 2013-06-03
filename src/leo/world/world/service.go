/* this is world service
*/

package world

import (
	"fmt"
	"path"
	"ini"
	"strconv"
)

type WorldService struct {
	master_port_id int
}

func NewWorldService() (service *WorldService, err error) {
	service = new(WorldService)
	err = service.init()
	return
}

func (service *WorldService) init() error {
	return nil
}

func (service *WorldService) Start() error {
	return nil
}

func (service *WorldService) Close() error {
	return nil
}

func (service *WorldService) Tick() error {
	//fmt.Println("world service tick")
	return nil
}

func (service *WorldService) Save() error {
	return nil
}

func (service *WorldService) connect_master() error {
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

func (service *WorldService) disconnect_master() error {
	fmt.Println("disconnect_master")
	Root.Port.CloseConnect(service.master_port_id)
	return nil
}