/* this is the session
*/

package gate

import (
	"net"
)

type Session struct {
	ip string
	port int

	uid string
	
}

func Create() (session *Session, err error) {
	session = nil
	err = nil

	//todo:

	return
}

func (ssn *Session) Start() {
	//todo:
}

func (ssn *Session) Close() {
	//todo：
}

func (ssn *Session) Update() {
	//todo:
}

func (ssn *Session) Send() {
	//todo:
}

func (ssn *Session) onrecv() {
	//todo:
}

func (ssn *Session) heartbeat() {
	//todo:
}
