/* this is rpc server
*/

package base

import (
	"fmt"
	"strings"
	"strconv"
	"net"
	"net/rpc"
	"runtime/debug"
)

type RpcServer struct {
	running bool
	addr string
	listener *net.TCPListener
	conns []*net.TCPConn
}

func NewRpcServer(ip string, port int) (server *RpcServer, err error) {
	server = new(RpcServer)
	err = server.init(ip, port)
	return
}

func (server *RpcServer) init(ip string, port int) error {
	arr := []string{ip, strconv.Itoa(port)}
	val := strings.Join(arr, ":")
	server.addr = val
	server.conns = make([]*net.TCPConn, 0)
	return nil
}

func (server *RpcServer) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", server.addr)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	server.listener = l
	server.running = true

	go server.handle_accept()

	return nil
}

func (server *RpcServer) Close() error {
	server.running = false

	for _, conn := range server.conns {
		conn.Close()
	}
	server.conns = make([]*net.TCPConn, 0)

	if server.listener != nil {
		server.listener.Close()
		server.listener = nil
	}

	return nil
}

func (server *RpcServer) IP() string {
	arr := strings.Split(server.addr, ":")
	if len(arr) > 0 {
		return arr[0]
	} else {
		return ""
	}
}

func (server *RpcServer) Port() int {
	arr := strings.Split(server.addr, ":")
	if len(arr) >1 {
		v, err := strconv.Atoi(arr[1])
		if err != nil {
			return 0
		} else {
			return v
		}
	} else {
		return 0
	}
}

func (server *RpcServer) RegisterService(sv interface{}) error {
	return rpc.Register(sv)
}

func (server *RpcServer) handle_accept() {
	defer func() {
		if r := recover(); r != nil {
			if LoggerIns != nil {
				LoggerIns.Critical(r, string(debug.Stack()))
			} else {
				fmt.Println("handle_accept except", r, string(debug.Stack()))
			}
		}
	}()

	for {
		if !server.running {
			break
		}

		conn, err := server.listener.AcceptTCP()
		if err != nil {
			if LoggerIns != nil {
				LoggerIns.Error(err)
				debug.PrintStack()
			} else {
				fmt.Println("accept tcp error:", err.Error())
			}
			continue
		}
		server.conns = append(server.conns, conn)

		go server.serve_rpc(conn)
	}
}

func (server *RpcServer) serve_rpc(conn *net.TCPConn) {
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

	rpc.ServeConn(conn)
}