/* this is connect manager
*/

package gate

import (
	"fmt"
	"strconv"
	"strings"
	"net"
	"sync"
//	"syscall"
	"runtime/debug"
)

type ConnectListener interface {
	Handle(conn *net.TCPConn)
}

type Acceptor struct {
	running bool
	addr string
	listen_count int
	listener *net.TCPListener
	connListeners []ConnectListener
	lock sync.Mutex
}

func NewAcceptor(ip string, port int, count int) (mgr *Acceptor, err error) {
	mgr = new(Acceptor)
	err = mgr.init(ip, port, count)
	if err != nil {
		mgr = nil
	}
	return
}

func (mgr *Acceptor) init(ip string, port int, count int) error {
	arr := []string{ip, strconv.Itoa(port)}
	val := strings.Join(arr, ":")
	mgr.addr = val
	if count <= 0 {
		count = 1
	}
	mgr.listen_count = count

	mgr.connListeners = make([]ConnectListener, 0, 8)
	return nil
}

func (mgr *Acceptor) Start() {
	mgr.running = true

	addr, err := net.ResolveTCPAddr("tcp", mgr.addr)
	if err != nil {
		Root.Logger.Critical("acceptor start failed:", err)
		return
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		Root.Logger.Critical("acceptor start failed:", err)
		return
	}

// 	file, err := listener.File()
// 	if err != nil {
// 		Root.Logger.Critical("get tcp listener file failed", err)
// 		return
// 	}
// 	err = syscall.SetsockoptInt(int(file.Fd()),
// 		syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
// 	if err != nil {
// 		Root.Logger.Critical("set tcp listener reuseaddr failed", err)
// 		return
// 	}
// 	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
// 	if err != nil {
// 		Root.Logger.Critical("get tcp listener file failed", err)
// 		return
// 	}
// 	err = syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
// 	if err != nil {
// 		Root.Logger.Critical("set tcp listener reuseaddr failed", err)
// 		return
// 	}
	mgr.listener = listener
	
	for i := 0; i < mgr.listen_count; i++ {
		go mgr.handle_accept()
	}
}

func (mgr *Acceptor) Close() {
	mgr.running = false
	if mgr.listener != nil {
		mgr.listener.Close()
	}
}

func (mgr *Acceptor) RegisterConnListener(l ConnectListener) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	for _, v := range(mgr.connListeners) {
		if v == l {
			return
		}
	}
	mgr.connListeners = append(mgr.connListeners, l)
}

func (mgr *Acceptor) UnRegisterConnListener(l ConnectListener) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	for i, v := range(mgr.connListeners) {
		if v == l {
			mgr.connListeners = append(mgr.connListeners[:i],
				mgr.connListeners[i+1:]...)
			break
		}
	}
}

func (mgr *Acceptor) handle_accept() {
	defer func() {
		if r := recover(); r != nil {
			if Root != nil && Root.Logger != nil {
				Root.Logger.Critical(r, string(debug.Stack()))
			} else {
				fmt.Println("handle_accept except", r, string(debug.Stack()))
			}
		}
	}()

	for {
		if !mgr.running {
			break
		}

		conn, err := mgr.listener.AcceptTCP()
		if err != nil {
			if Root != nil && Root.Logger != nil {
				Root.Logger.Error(err)
				debug.PrintStack()
			} else {
				fmt.Println("accept tcp error:", err.Error())
			}
			continue
		}
		mgr.handle_conn(conn)
	}
}

func (mgr *Acceptor) handle_conn(conn *net.TCPConn) {
//	conn.SetNoDelay(true)
//	conn.SetKeepAlive(true)
//	conn.SetLinger(0)

	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	for _, v := range mgr.connListeners {
		v.Handle(conn)
	}
}