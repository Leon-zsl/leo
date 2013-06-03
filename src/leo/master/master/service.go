/* this is master service
*/

package master

import (
	"leo/base"
)

type MasterService struct {
	Clock *base.Clock
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