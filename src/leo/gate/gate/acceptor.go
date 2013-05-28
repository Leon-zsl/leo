/* this is connect manager
*/

package gate

import (
	"fmt"
	"strconv"
	"strings"
	"net"
	"sync"
	"runtime"
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
	addr, err := net.ResolveTCPAddr("tcp", val)
	if err != nil {
		return err
	}
	mgr.addr = val
	mgr.listen_count = count

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	mgr.listener = listener
	
	mgr.connListeners = make([]ConnectListener, 0, 8)
	return nil
}

func (mgr *Acceptor) Start() {
	mgr.running = true
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

func (mgr *Acceptor) Update() {
	//todo:
}

func (mgr *Acceptor) handle_accept() {
	defer func() {
		if r := recover(); r != nil {
			if Root != nil && Root.Logger != nil {
				fmt.Println("accept exception caught!")
				Root.Logger.Critical(r)
			} else {
				fmt.Println("handle_accept exception")
			}
		}
	}()

	for {
		conn, err := mgr.listener.AcceptTCP()
		if err != nil {
			if !mgr.running {
				runtime.Goexit()
			}

			if Root != nil && Root.Logger != nil {
				fmt.Println("accept error caught!")
				Root.Logger.Error(err)
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