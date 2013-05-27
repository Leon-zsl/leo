/* this is the gate entry
 */
package gate

import (
	"fmt"
)

type Gate struct {
	Running    bool
	SessionMgr *SessionMgr
}

func Create() (gate *Gate, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("create gate failed!")
			gate = nil
			err = -1
		}
	}()

	gt := new(Gate)
	e = gt.Init()
	if e != nil {
		fmt.Println("init gate failed")
		gate = nil
		err = e
	} else {
		gate = gt
		err = e
	}

	return
}

func (gate *Gate) Init() error {
	return nil
}

func (gate *Gate) Run() {
	//todo:
}

func (gate *Gate) Close() error {
	Running = false
	gate.SessionMgr.Close()
	return nil
}
