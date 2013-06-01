/* this is rpc server
*/

package base

import (
	"fmt"
	"net/rpc"
	"runtime/debug"
)

type RpcServer struct {
	running bool
	acceptor *Acceptor
	ssns []*Session
}

func NewRpcServer(ip string, port int) (server *RpcServer, err error) {
	server = new(RpcServer)
	err = server.init(ip, port)
	return
}

func (server *RpcServer) init(ip string, port int) error {
	ac, err := NewAcceptor(ip, port, 1)
	if err != nil {
		return err
	}
	ac.RegisterAcceptedSessionListener(server)
	server.acceptor = ac

	server.ssns = make([]*Session, 0)
	return nil
}

func (server *RpcServer) Start() {
	server.acceptor.Start()
	server.running = true
}

func (server *RpcServer) Close() {
	server.running = false

	for _, ssn := range server.ssns {
		ssn.Close()
	}
	server.ssns = make([]*Session, 0)

	if server.acceptor != nil {
		server.acceptor.UnRegisterAcceptedSessionListener(server)
		server.acceptor.Close()
		server.acceptor = nil
	}
}

func (server *RpcServer) Register(sv interface{}) error {
	return rpc.Register(sv)
}

func (server *RpcServer) HandleAcceptedSession(ssn *Session) {
	if ssn == nil {
		return
	}
	server.ssns = append(server.ssns, ssn)

	go server.serve_rpc(ssn)
}

func (server *RpcServer) serve_rpc(ssn *Session) {
	defer func() {
		if !server.running {
			return
		}
		if r := recover(); r != nil {
			if LoggerIns != nil {
				LoggerIns.Critical(r, string(debug.Stack()))
			} else {
				fmt.Println(r, string(debug.Stack()))
			}
		}
	}()

	rpc.ServeConn(ssn.Conn())
}