/*this is main service
*/

package db

import (
	
)

type Service struct {
}

func NewService() (service *Service, err error) {
	service = new(Service)
	err = service.init()
	return
}

func (service *Service) init() error {
	return nil
}

func (service *Service) Start() {
	
}

func (service *Service) Close() {
}

func (service *Service) Tick() {
}