/*this is main service
 */

package db

import (
//	"path"
// 	"ini"
// 	"strconv"
// 	"fmt"
//	"leo/base"
)

type DBService struct {
	//master_port_id int
	//Clock *base.Clock
}

func NewDBService() (service *DBService, err error) {
	service = new(DBService)
	err = service.init()
	return
}

func (service *DBService) init() error {
	//service.Clock, _ = base.NewClock()
	return nil
}

func (service *DBService) Start() error {
	//service.Clock.Start()
	//service.connect_master()
	return nil
}

func (service *DBService) Close() error {
	//service.disconnect_master()
	//service.Clock.Close()
	return nil
}

func (service *DBService) Save() error {
	return nil
}

func (service *DBService) Tick() error {
	//service.Clock.Tick()
	//fmt.Println("db service tick")
	return nil
}
