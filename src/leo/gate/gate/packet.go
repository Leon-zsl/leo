/* this is packet
*/

package gate

type Packet struct {
	Op int
	Args []byte
}

func (pk *Packet) Write(op int, args []byte) {
	//todo:
}
