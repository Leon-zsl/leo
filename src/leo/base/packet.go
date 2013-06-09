/* this is packet
*/

package base

import (
	"bytes"
	"encoding/binary"
)

//safe for goroutine
type Packet struct {
	Op int32
	Args []byte
}

func NewPacket(opcode int, argements []byte) *Packet {
	return &Packet{Op:int32(opcode), Args:argements}
}

func NewPacketFromBytes(b []byte) (pkt *Packet, err error) {
	if b == nil {
		return
	}

	var op int32 = 0
	buf := bytes.NewBuffer(b[:4])
	err = binary.Read(buf, binary.BigEndian, &op)
	if err != nil {
		pkt = nil
		return
	}
	
	pkt = new(Packet)
	pkt.Op = op
	pkt.Args = b[4:]

	return
}

// func (pkt *Packet) Op() int32 {
// 	return pkt.op
// }

// func (pkt *Packet) Args() []byte {
// 	return pkt.args
// }

func (pkt *Packet) Bytes() (b []byte, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, pkt.Op)
	if err != nil {
		b = nil
		return
	}

	b = append(buf.Bytes(), pkt.Args...)
	return
}