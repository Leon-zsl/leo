/*this is main service
*/

package db

import (
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

type Service struct {
}

func NewDBService() (service *Service, err error) {
	service = new(Service)
	err = service.init()
	return
}

func (service *Service) init() error {
	return nil
}

func (service *Service) Start() error {
	return nil
}

func (service *Service) Close() error {
	return nil
}

func (service *Service) Save() error {
	return nil
}

func (service *Service) Tick() error {
	test_add()
	test_get()
	test_set()
	test_del()
	time.Sleep(1e9)
	return nil
}