/*this is main service
*/

package db

import (
	"path"
	"ini"
	"strconv"
	"time"
	"fmt"
	"leo/base"
)

func test_add() {
	rcd, _ := base.NewRecord()

	rcd.SetValue("uid", 1)
	rcd.SetValue("account", "test")
	rcd.SetValue("pwd", "test")
	err := Root.Driver.Add("main_user", 1, "uid", rcd)
	if err != nil {
		fmt.Println("add err:", err)
	} else {
		fmt.Println("add ok")
	}
}

func test_set() {
	rcd, _ := base.NewRecord()

	rcd.SetValue("uid", 1)
	rcd.SetValue("account", "test")
	rcd.SetValue("pwd", "test")
	err := Root.Driver.Set("main_user", 1, "uid", rcd)
	if err != nil {
		fmt.Println("set err:", err)
	} else {
		fmt.Println("set ok")
	}
}

func test_get() {
	rcd, err := Root.Driver.Get("main_user", 1, "uid")
	if err != nil {
		fmt.Println("get err:", err)
	} else {
		fmt.Println("get ok")
		fmt.Println(rcd.Values())
	}
}

func test_del() {
	err := Root.Driver.Del("main_user", 1, "uid")
	if err != nil {
		fmt.Println("del err:", err)
	} else {
		fmt.Println("del ok")
	}
}

type DBService struct {
	master_port_id int
}

func NewDBService() (service *DBService, err error) {
	service = new(DBService)
	err = service.init()
	return
}

func (service *DBService) init() error {
	return nil
}

func (service *DBService) Start() error {
	return nil
}

func (service *DBService) Close() error {
	return nil
}

func (service *DBService) Save() error {
	return nil
}

func (service *DBService) Tick() error {
	//fmt.Println("db service tick")
	test_add()
	test_get()
	test_set()
	test_del()
	time.Sleep(1e9)
	return nil
}

func (service *DBService) connect_master() error {
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

func (service *DBService) disconnect_master() error {
	fmt.Println("disconnect_master")
	Root.Port.CloseConnect(service.master_port_id)
	return nil
}