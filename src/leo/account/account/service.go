/* this is account service
*/

package account

import (
	"fmt"
	"path"
	"ini"
	"strconv"
)

type AccountService struct {
	master_port_id int
}

func NewAccountService() (service *AccountService, err error) {
	service = new(AccountService)
	err = service.init()
	return
}

func (service *AccountService) init() error {
	return nil
}

func (service *AccountService) Start() error {
	service.connect_master()
	return nil
}

func (service *AccountService) Close() error {
	service.disconnect_master()
	return nil
}

func (service *AccountService) Tick() error {
	//fmt.Println("account service tick")
	return nil
}

func (service *AccountService) Save() error {
	return nil
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