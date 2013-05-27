/* this is connect manager
*/

package gate

import (
	"net"
)

type Accept struct {
	Running bool
}

func Create() (mgr *Accept, err error) {
	mgr = nil
	err = nil
	//todo:
	return
}

func (mgr *Accept) Start() {
	//todo:
}

func (mgr *Accept) Close() {
	//todo:
}

func (mgr *Accept) Run() {
	//todo:
}