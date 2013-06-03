/* this is the session
*/

package base

import (
	"fmt"
	"bytes"
	"encoding/binary"
	"time"
	"strings"
	"strconv"
	"net"
	"uuid"
	"runtime/debug"
)

type SessionHandler interface {
	HandleSessionStart(ssn *Session)
	HandleSessionMsg(ssn *Session, pkt* Packet)
	HandleSessionClose(ssn *Session)
	HandleSessionError(ssn *Session, err error)
}

type Session struct {
	closed bool

	addr string
	sid string

	conn *net.TCPConn

//	recvq *RingBuffer
	sendq *RingBuffer

	handlers []SessionHandler
}

func NewSession(conn *net.TCPConn) (session *Session, err error) {
	session = new(Session)
	err = session.init(conn)
	return
}

func (ssn *Session) init(conn *net.TCPConn) error {
	ssn.addr = conn.RemoteAddr().String()
	ssn.sid = uuid.New()

	ssn.handlers = make([]SessionHandler, 0)

	ssn.sendq = NewRingBuffer()
//	ssn.recvq = NewRingBuffer()

	conn.SetReadBuffer(2048)
	//	conn.SetNoDelay(true)
	//	conn.SetKeepAlive(true)
	//	conn.SetLinger(0)

	ssn.conn = conn
	
	return nil
}

func (ssn *Session) Start() error {
	ssn.closed = false

	//no handler, close the session
	if len(ssn.handlers) == 0 {
		ssn.Close()
		return nil
	}

	for _, l := range(ssn.handlers) {
		l.HandleSessionStart(ssn)
	}

	//go ssn.onmsg()
	go ssn.onsend()
	go ssn.onrecv()

	return nil
}

func (ssn *Session) Closed() bool {
	return ssn.closed
}

func (ssn *Session) Close() error {
	for _, v := range(ssn.handlers) {
		v.HandleSessionClose(ssn)
	}

	if ssn.conn != nil {
		ssn.conn.Close()
		ssn.conn = nil
	}
	ssn.closed = true

	return nil
}

func (ssn *Session) Conn() *net.TCPConn {
	return ssn.conn
}

func (ssn *Session) RegisterHandler(h SessionHandler) {
	for _, v := range(ssn.handlers) {
		if v == h {
			LoggerIns.Warn("duplicate session handler")
			return
		}
	}
	ssn.handlers = append(ssn.handlers, h)
}

func (ssn *Session) UnRegisterHandler(h SessionHandler) {
	for i, v := range(ssn.handlers) {
		if v == h {
			ssn.handlers = append(ssn.handlers[:i],
				ssn.handlers[i+1:]...)
			break
		}
	}
}

func (ssn *Session) Send(pk *Packet) {
	if pk == nil {
		return
	}
	ssn.sendq.Push(pk)
}

func (ssn *Session) IP() string {
	return strings.Split(ssn.addr, ":")[0]
}

func (ssn *Session) Port() int {
	p, _ := strconv.Atoi(strings.Split(ssn.addr, ":")[1])
	return p
}

func (ssn *Session) Addr() string {
	return ssn.addr
}

func (ssn *Session) SID() string {
	return ssn.sid
}

func (ssn *Session) handle_send_err(err error) {
	if len(ssn.handlers) > 0 {
		for _, v := range(ssn.handlers) {
			v.HandleSessionError(ssn, err)
		}
	} else {
		LoggerIns.Error("write session failed: " + err.Error())
	}
	ssn.Close()
}

func (ssn *Session) handle_recv_err(err error) {
	if len(ssn.handlers) > 0 {
		for _, v := range(ssn.handlers) {
			v.HandleSessionError(ssn, err)
		}
	} else {
		LoggerIns.Error("read session failed: " + err.Error())
	}
	ssn.Close()
}

// func (ssn *Session) onmsg() {
// 	for {
// 		for ssn.recvq.Empty() {
// 			time.Sleep(1e6)
// 			continue
// 		}

// 		pk := ssn.recvq.Pop()
// 		for _, h := range(ssn.handlers) {
// 			h.HandleSessionMsg(ssn, pk)
// 		}
// 	}
// }

func (ssn *Session) onrecv() {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				ssn.handle_recv_err(err)
			} else if LoggerIns != nil {
				LoggerIns.Critical(r, string(debug.Stack()))
			} else {
				fmt.Println("recv exception caught", r, string(debug.Stack()))
			}
		}
	}()

	var tmpbuf [2048]byte
	recvbuf := make([]byte, 0)
	for {
		if ssn.closed {
			break
		}

		l, err := ssn.conn.Read(tmpbuf[0:])
		if err != nil {
			ssn.handle_recv_err(err)
			continue
		}

		if l == 0 {
			continue
		}
		
		recvbuf = append(recvbuf, tmpbuf[:l]...)
		for {
			if len(recvbuf) == 0 {
				break
			}

			var ln int32 = 0
			buf := bytes.NewBuffer(recvbuf[:4])
			err := binary.Read(buf, binary.BigEndian, &ln)
			if err != nil {
				ssn.handle_recv_err(err)
				break
			}
			
			if len(recvbuf) < int(ln + 4) {
				break
			}

			pk, err := NewPacketFromBytes(recvbuf[4:ln+4])
			if err != nil {
				ssn.handle_recv_err(err)
				break
			}

			for _, h := range(ssn.handlers) {
				h.HandleSessionMsg(ssn, pk)
			}

 			recvbuf = recvbuf[ln+4:]

// 			ssn.recvq.Push(pk)
// 			recvbuf = recvbuf[ln+4:]
		}
	}
}

func (ssn *Session) onsend() {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				ssn.handle_send_err(err)
			} else if LoggerIns != nil {
				LoggerIns.Critical(r, string(debug.Stack()))
			} else {
				fmt.Println("send exception caught", r, string(debug.Stack()))
			}
		}
	}()

	for {
		if ssn.closed {
			break
		}

		for ssn.sendq.Empty() {
			time.Sleep(1e6 * 30) //wait 30ms
			continue
		}

		sendbuf := make([]byte, 0)
		for {
			pkt := ssn.sendq.Pop()
			if pkt == nil {
				break
			}

			buffer, err := pkt.Bytes()
			if err != nil {
				ssn.handle_send_err(err)
				break
			}

			var l int32 = 0
			l = int32(len(buffer))
			if l == 0 {
				continue
			}

			buf := new(bytes.Buffer)
			err = binary.Write(buf, binary.BigEndian, l)
			if err != nil && !ssn.closed {
				ssn.handle_send_err(err)
				break
			}

			buffer = append(buf.Bytes(), buffer...)
			sendbuf = append(sendbuf, buffer...)
			//sendbuf = append(sendbuf, buf.Bytes()..., buffer...)
		}

		if ssn.closed {
			break
		}

		ln := 0
		for ln < len(sendbuf) {
			l, err := ssn.conn.Write(sendbuf[ln:])
			if err != nil && !ssn.closed {
				ssn.handle_send_err(err)
				break
			}
			ln += l
		}
	}
}

