/* this is the app
*/

package robot

import pblib "code.google.com/p/goprotobuf/proto"

import (
	"fmt"
	"net"
	"bytes"
	"encoding/binary"
	"time"
	"leo/base"
//	"leo/common"
	"leo/proto"
)

type App struct {
	Conn *net.TCPConn;
}

var (
	Root *App
)

func NewApp() (app *App, err error) {
	app = new(App)
	err = app.init()
	Root = app
	return
}

func (app *App) init() error {
	return nil
}

func (app *App) Sendpkt(pkt *base.Packet) {
	buffer, _ := pkt.Bytes()
	l := int32(len(buffer))
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, l)
	buffer = append(buf.Bytes(), buffer...)
	
	_, err := app.Conn.Write(buffer)
	if err != nil {
		fmt.Println("write err:", err.Error())
	}	
}

func (app *App) Recvpkt() *base.Packet {
	tmp := make([]byte, 512)
	_, err := app.Conn.Read(tmp)
	if err != nil {
		fmt.Println("read err:", err.Error())
		return nil
	}

	var l int32 = 0
	buf := bytes.NewBuffer(tmp[:4])
	binary.Read(buf, binary.BigEndian, &l)

	pkt, _ := base.NewPacketFromBytes(tmp[4:l+4])
	return pkt
}

func (app *App) Startup() {
	con, err := net.Dial("tcp", "127.0.0.1:10001")
	if err != nil {
		fmt.Println("net dial failed:", err.Error())
		return
	}

	defer con.Close()
	app.Conn = con.(*net.TCPConn)

	test_all()
}

func test_all() {
	test_login()
	time.Sleep(1e9)
}

func test_login() {
	fmt.Println("login...")
	pb := &proto.Login{ Name : pblib.String("test"), Pwd : pblib.String("test")}
	val, _ := pblib.Marshal(pb)
	pkt := base.NewPacket(proto.LOGIN, val)
	Root.Sendpkt(pkt)
	
	pkt = Root.Recvpkt()
	pbr := &proto.LoginResp{}
	
	err := pblib.Unmarshal(pkt.Args, pbr)
	if err != nil {
		fmt.Println("parse proto failed:", err)
		return
	}
	fmt.Println("login resp:", pbr.GetErrorCode(), pbr.GetErrorMsg())
}