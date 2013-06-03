/* this is master service
*/

package master

type MasterService struct {
}

func NewMasterService() (service *MasterService, err error) {
	service = new(MasterService)
	err = service.init()
	return
}

func (service *MasterService) init() error {
	return nil
}

func (service *MasterService) Start() error {
	return nil
}

func (service *MasterService) Close() error {
	return nil
}

func (service *MasterService) Tick() error {
	//fmt.Println("world service tick")
	return nil
}

func (service *MasterService) Save() error {
	return nil
}