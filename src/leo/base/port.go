/* this is server port
 */

package base

import (
	//	"fmt"
	"errors"
	"strconv"
	"sync"
)

//safe for goroutine
type Port struct {
	running bool

	id     int
	server *RpcServer

	lock    sync.RWMutex
	clients map[int]*RpcClient
}

type PortData struct {
	ID   int
	IP   string
	Port int
}

type HandShakeService struct {
	owner *Port
}

func (service *HandShakeService) Open(cl *PortData, sv *PortData) error {
	if cl == nil {
		return errors.New("cl arg is nil")
	}
	if sv == nil {
		return errors.New("sv arg is nil")
	}

	if _, ok := service.owner.clients[cl.ID]; ok {
		sv.ID = service.owner.ID()
		sv.IP = service.owner.Server().IP()
		sv.Port = service.owner.Server().Port()
		return nil
	}

	err := service.owner.OpenConnect(cl.ID, cl.IP, cl.Port)
	if err != nil {
		return err
	}
	sv.ID = service.owner.ID()
	sv.IP = service.owner.Server().IP()
	sv.Port = service.owner.Server().Port()
	return nil
}

func (service *HandShakeService) Close(cl *PortData, sv *PortData) error {
	if cl == nil {
		return errors.New("cl arg is nil")
	}
	if sv == nil {
		return errors.New("sv arg is nil")
	}

	if _, ok := service.owner.clients[cl.ID]; !ok {
		sv.ID = service.owner.ID()
		sv.IP = service.owner.Server().IP()
		sv.Port = service.owner.Server().Port()
		return nil
	}

	err := service.owner.CloseConnect(cl.ID)
	if err != nil {
		return err
	}
	sv.ID = service.owner.ID()
	sv.IP = service.owner.Server().IP()
	sv.Port = service.owner.Server().Port()
	return nil
}

func NewPort(id int, ip string, port int) (p *Port, err error) {
	p = new(Port)
	err = p.init(id, ip, port)
	return
}

func (port *Port) init(id int, ip string, pt int) error {
	s, err := NewRpcServer(ip, pt)
	if err != nil {
		return err
	}
	port.server = s
	port.id = id

	port.clients = make(map[int]*RpcClient)

	hs := &HandShakeService{owner: port}
	port.RegisterService(hs)
	return nil
}

func (port *Port) Start() error {
	return port.server.Start()
}

func (port *Port) Close() error {
	for _, v := range port.clients {
		v.Close()
	}
	port.clients = make(map[int]*RpcClient)

	if port.server != nil {
		port.server.Close()
		port.server = nil
	}

	return nil
}

func (port *Port) ID() int {
	return port.id
}

func (port *Port) Server() *RpcServer {
	return port.server
}

func (port *Port) RegisterService(sv interface{}) error {
	return port.server.RegisterService(sv)
}

func (port *Port) OpenConnect(port_id int, ip string, pt int) error {
	_, ok := port.clients[port_id]
	if ok {
		return errors.New("duplicate port id" + strconv.Itoa(port_id))
	}
	cl, err := NewRpcClient(ip, pt)
	if err != nil {
		return err
	}
	err = cl.Start()
	if err != nil {
		return err
	}
	port.lock.Lock()
	port.clients[port_id] = cl
	port.lock.Unlock()

	pd := &PortData{port.id, port.server.IP(), port.server.Port()}
	ps := new(PortData)
	err = cl.Call("HandShakeService.Open", pd, ps)
	return err
}

func (port *Port) CloseConnect(port_id int) error {
	v, ok := port.clients[port_id]
	if ok {
		port.lock.Lock()
		delete(port.clients, port_id)
		port.lock.Unlock()

		pd := &PortData{port.id, port.server.IP(), port.server.Port()}
		ps := new(PortData)
		err := v.Call("HandShakeService.Close", pd, ps)
		if err != nil {
			LoggerIns.Error("close connect err:", err)
		}

		v.Close()
	}
	return nil
}

func (port *Port) GetConnect(port_id int) (*RpcClient, bool) {
	port.lock.RLock()
	cl, ok := port.clients[port_id]
	port.lock.RUnlock()
	return cl, ok
}

func (port *Port) Call(port_id int, method string, args interface{}, reply interface{}) error {
	cl, ok := port.GetConnect(port_id)
	if !ok {
		return errors.New("target port do not exists " + strconv.Itoa(port_id))
	}
	return cl.Call(method, args, reply)
}

func (port *Port) CallAsync(port_id int, method string, args interface{}, reply interface{}, cb RpcCallback) error {
	cl, ok := port.GetConnect(port_id)
	if !ok {
		return errors.New("target port do not exists " + strconv.Itoa(port_id))
	}
	cl.CallAsync(method, args, reply, cb)
	return nil
}

func (port *Port) SendAsync(port_id int, method string, args interface{}) error {
	cl, ok := port.GetConnect(port_id)
	if !ok {
		return errors.New("target port do not exists " + strconv.Itoa(port_id))
	}
	var v int = 0
	cl.CallAsync(method, args, &v, nil)
	return nil
}
