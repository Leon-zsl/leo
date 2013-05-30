/* this is connect to master
*/

package gate

import (
	"net"
	"strings"
	"strconv"
)

type Connector struct {
	//todo
}

func NewConnector () (conn *Connector, err error) {
	conn = new(Connector)
	err = conn.init()
	return
}

func (conn *Connector) init() error {
	return nil
}

func (conn *Connector) Start() {
	//do nothing
}

func (conn *Connector) Close() {
	//do nothing
}

//the caller need to invoke the sessionmgr.AddSession(ssn) and ssn.Start()
func (conn *Connector) Connect(ip string, port int) (ssn *Session, err error) {
	arr := []string{ip, strconv.Itoa(port)}
	val := strings.Join(arr, ":")
	addr, err := net.ResolveTCPAddr("tcp", val)
	if err != nil {
		return
	}

	tcpconn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return
	}

	ssn, err = NewSession(tcpconn)
	if err != nil {
		ssn = nil
		return
	}

	return
}