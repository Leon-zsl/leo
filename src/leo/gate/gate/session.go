/* this is the session
*/

package gate

import (
	"fmt"
	"bytes"
	"encoding/binary"
	"time"
	"strings"
	"strconv"
	"net"
	"uuid"

	"leo/base"
)

type ErrorListener interface {
	Handle(err error)
}

type Session struct {
	closed bool

	addr string
	sid string

	recvq *base.RingBuffer
	sendq *base.RingBuffer

	recvbuf []byte
	tempbuf []byte
	conn *net.TCPConn

	senderrlistener ErrorListener
	recverrlistener ErrorListener
}

func NewSession(conn *net.TCPConn) (session *Session, err error) {
	session = new(Session)
	err = session.init(conn)
	return
}

func (ssn *Session) init(conn *net.TCPConn) error {
	ssn.closed = false
	ssn.addr = conn.RemoteAddr().String()
	ssn.sid = uuid.New()

	ssn.recvq = base.NewRingBuffer()
	ssn.sendq = base.NewRingBuffer()

	ssn.recvbuf = make([]byte, 0)
	ssn.tempbuf = make([]byte, 2048)
	conn.SetReadBuffer(2048)
	
	go ssn.onsend()
	go ssn.onrecv()

	return nil
}

func (ssn *Session) Closed() bool {
	return ssn.closed
}

func (ssn *Session) Close() {
	if ssn.conn != nil {
		ssn.conn.Close()
		ssn.conn = nil
	}
	ssn.closed = true
}

func (ssn *Session) Update() {
	//todo:
}

func (ssn *Session) Recv() *base.Packet {
	return ssn.recvq.Pop()
}

func (ssn *Session) Send(pk *base.Packet) {
	if pk == nil {
		return
	}
	ssn.sendq.Push(pk)
}

func (ssn *Session) onrecv() {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				ssn.handle_recv_err(err)
			} else if Root != nil && Root.Logger != nil {
				Root.Logger.Critical(r)
			} else {
				fmt.Println("recv exception caught", ssn.addr)
			}
		}
	}()

	for {
		if ssn.closed {
			break
		}

		l, err := ssn.conn.Read(ssn.tempbuf)
		if err != nil {
			ssn.handle_recv_err(err)
			continue
		}

		if l == 0 {
			continue
		}

		ssn.recvbuf = append(ssn.recvbuf, ssn.tempbuf[:l]...)
		for {
			var ln int32 = 0
			buf := bytes.NewBuffer(ssn.recvbuf[:4])
			err := binary.Read(buf, binary.BigEndian, &ln)
			if err != nil {
				ssn.handle_recv_err(err)
				break
			}
			
			if len(ssn.recvbuf) < int(ln) {
				break
			}

			pk, err := base.NewPacketFromBytes(ssn.recvbuf[4:ln])
			if err != nil {
				ssn.handle_recv_err(err)
				break
			}

			ssn.recvq.Push(pk)
			ssn.recvbuf = ssn.recvbuf[ln:]
		}
	}
}

func (ssn *Session) onsend() {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				ssn.handle_send_err(err)
			} else if Root != nil && Root.Logger != nil {
				Root.Logger.Critical(r)
			} else {
				fmt.Println("send exception caught", ssn.addr)
			}
		}
	}()

	for {
		if ssn.closed {
			break
		}

		for ssn.sendq.Empty() {
			time.Sleep(1e6 * 30) //wait 30ms
		}

		sendbuf := make([]byte, 0)
		for {
			if ssn.closed {
				break
			}

			pkt := ssn.sendq.Pop()
			if pkt == nil {
				break
			}

			buffer, err := pkt.Bytes()
			if err != nil {
				ssn.handle_send_err(err)
				continue
			}

			var l int32 = 0
			l = int32(len(buffer))
			if l == 0 {
				continue
			}
			l += 4

			buf := new(bytes.Buffer)
			err = binary.Write(buf, binary.BigEndian, l)
			if err != nil && !ssn.closed {
				ssn.handle_send_err(err)
				continue
			}
			buffer = append(buf.Bytes(), buffer...)
			sendbuf = append(sendbuf, buffer...)
			//sendbuf = append(sendbuf, buf.Bytes()..., buffer...)
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

func (ssn *Session) SetSendErrListener(l ErrorListener) {
	ssn.senderrlistener = l
}

func (ssn *Session) SetRecvErrListener(l ErrorListener) {
	ssn.recverrlistener = l
}

func (ssn *Session) handle_send_err(err error) {
	if ssn.senderrlistener != nil {
		ssn.senderrlistener.Handle(err)
	} else {
		Root.Logger.Error("write session failed: " + err.Error())
	}
	ssn.Close()
}

func (ssn *Session) handle_recv_err(err error) {
	if ssn.recverrlistener != nil {
		ssn.recverrlistener.Handle(err)
	} else {
		Root.Logger.Error("write session failed: " + err.Error())
	}
	ssn.Close()
}