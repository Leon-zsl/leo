/* this is connect manager
*/

package base

import (
	"fmt"
	"strconv"
	"strings"
	"net"
	"runtime/debug"
)

type AcceptedSessionListener interface {
	HandleAcceptedSession(ssn *Session)
}

//goroutine safe if register before start
type Acceptor struct {
	running bool
	addr string
	listen_count int
	listener *net.TCPListener

	ssnListeners []AcceptedSessionListener
}

//var AcceptorIns *Acceptor = nil

func NewAcceptor(ip string, port int, count int) (mgr *Acceptor, err error) {
	mgr = new(Acceptor)
	err = mgr.init(ip, port, count)
	//AcceptorIns = mgr
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

	mgr.ssnListeners = make([]AcceptedSessionListener, 0)
	return nil
}

func (mgr *Acceptor) Start() error {
	mgr.running = true

	addr, err := net.ResolveTCPAddr("tcp", mgr.addr)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

// 	file, err := listener.File()
// 	if err != nil {
// 		LoggerIns.Critical("get tcp listener file failed", err)
// 		return
// 	}
// 	err = syscall.SetsockoptInt(int(file.Fd()),
// 		syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
// 	if err != nil {
// 		LoggerIns.Critical("set tcp listener reuseaddr failed", err)
// 		return
// 	}
// 	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
// 	if err != nil {
// 		LoggerIns.Critical("get tcp listener file failed", err)
// 		return
// 	}
// 	err = syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
// 	if err != nil {
// 		LoggerIns.Critical("set tcp listener reuseaddr failed", err)
// 		return
// 	}
	mgr.listener = listener
	
	for i := 0; i < mgr.listen_count; i++ {
		go mgr.handle_accept()
	}
	return nil
}

func (mgr *Acceptor) Close() error {
	mgr.running = false
	if mgr.listener != nil {
		mgr.listener.Close()
	}
	return nil
}

func (mgr *Acceptor) IP() string {
	arr := strings.Split(mgr.addr, ":")
	if len(arr) > 0 {
		return arr[0]
	} else {
		return ""
	}
}

func (mgr *Acceptor) Port() int {
	arr := strings.Split(mgr.addr, ":")
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

func (mgr *Acceptor) RegisterAcceptedSessionListener(l AcceptedSessionListener) {
// 	mgr.connlock.Lock()
// 	defer mgr.connlock.Unlock()

	for _, v := range(mgr.ssnListeners) {
		if v == l {
			return
		}
	}
	mgr.ssnListeners = append(mgr.ssnListeners, l)
}

func (mgr *Acceptor) UnRegisterAcceptedSessionListener(l AcceptedSessionListener) {
// 	mgr.connlock.Lock()
// 	defer mgr.connlock.Unlock()

	for i, v := range(mgr.ssnListeners) {
		if v == l {
			mgr.ssnListeners = append(mgr.ssnListeners[:i],
				mgr.ssnListeners[i+1:]...)
			break
		}
	}
}

func (mgr *Acceptor) handle_accept() {
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
		if !mgr.running {
			break
		}

		conn, err := mgr.listener.AcceptTCP()
		if err != nil {
			if LoggerIns != nil {
				LoggerIns.Error(err)
				debug.PrintStack()
			} else {
				fmt.Println("accept tcp error:", err.Error())
			}
			continue
		}

		// 	mgr.connlock.Lock()
		// 	defer mgr.connlock.Unlock()

		if len(mgr.ssnListeners) == 0 {
			conn.Close()
			continue
		}

		ssn, err := NewSession(conn)
		if err != nil {
			LoggerIns.Error("create new session failed: " + conn.RemoteAddr().String())
			return
		}

		for _, v := range mgr.ssnListeners {
			v.HandleAcceptedSession(ssn)
		}
		ssn.Start()
	}
}