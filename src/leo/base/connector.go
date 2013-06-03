/* this is connect to master
*/

package base

import (
	"net"
	"strings"
	"strconv"
)

type Connector struct {
	//todo
}

//var ConnectorIns *Connector = nil

func NewConnector () (conn *Connector, err error) {
	conn = new(Connector)
	err = conn.init()
	//ConnectorIns = conn
	return
}

func (conn *Connector) init() error {
	return nil
}

func (conn *Connector) Start() error {
	return nil
}

func (conn *Connector) Close() error {
	return nil
}

//the caller need to invoke the sessionmgr.AddSession(ssn) and ssn.Start()
func (conn *Connector) Connect(ip string, port int) (ssn *Session, err error) {
	arr := []string{ip, strconv.Itoa(port)}
	val := strings.Join(arr, ":")
	addr, err := net.ResolveTCPAddr("tcp", val)
	if err != nil {
		LoggerIns.Error("connect err:", ip, port, err)
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