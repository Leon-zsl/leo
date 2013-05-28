/* this is packet
*/

package base

import (
	"bytes"
	"encoding/binary"
)
type Packet struct {
	op int32
	args []byte
}

func NewPacket(opcode int, argements []byte) *Packet {
	return &Packet{op:int32(opcode), args:argements}
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
	pkt.op = op
	pkt.args = b[4:]

	return
}

func (pkt *Packet) Op() int32 {
	return pkt.op
}

func (pkt *Packet) Args() []byte {
	return pkt.args
}

func (pkt *Packet) Bytes() (b []byte, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, pkt.op)
	if err != nil {
		b = nil
		return
	}

	b = append(buf.Bytes(), pkt.args...)
	return
}