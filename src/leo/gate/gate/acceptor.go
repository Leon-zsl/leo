/* this is connect manager
*/

package gate

// import (
// 	"net"
// )

type Acceptor struct {
	running bool
	port int
}

func NewAcceptor(port int) (mgr *Acceptor, err error) {
	mgr = nil
	err = nil
	//todo:
	return
}

func (mgr *Acceptor) Start(port int) {
	//todo:
}

func (mgr *Acceptor) Close() {
	//todo:
}

func (mgr *Acceptor) Run() {
	//todo:
}