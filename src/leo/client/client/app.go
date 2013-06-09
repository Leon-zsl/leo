/* this is the app
*/

package client

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

func (app *App) sendpkt(pkt *base.Packet) {
	fmt.Println("send pkt...")

	buffer, _ := pkt.Bytes()
	l := int32(len(buffer))
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, l)
	buffer = append(buf.Bytes(), buffer...)
	
	_, err := app.Conn.Write(buffer)
	if err != nil {
		fmt.Println("write err:", err.Error())
	}	

	fmt.Println("send pkt end...")
}

func (app *App) recvpkt() *base.Packet {
	fmt.Println("recv pkt...")
	tmp := make([]byte, 512)
	_, err := app.Conn.Read(tmp)
	if err != nil {
		fmt.Println("read err:", err.Error())
		return nil
	}

	var l int32 = 0
	buf := bytes.NewBuffer(tmp[:4])
	binary.Read(buf, binary.BigEndian, &l)

	fmt.Println("recv len ", l)
	pkt, _ := base.NewPacketFromBytes(tmp[4:l+4])

	fmt.Println("recv pkt end...")
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

	for {
		pb := &proto.Login{ Name : pblib.String("test"), Pwd : pblib.String("test")}
		val, _ := pblib.Marshal(pb)
		pkt := base.NewPacket(proto.LOGIN, val)
		app.sendpkt(pkt)
		
		pkt = app.recvpkt()
		pbr := &proto.LoginResp{}
	
		fmt.Println("recv: ", pkt.Op)
		err := pblib.Unmarshal(pkt.Args, pbr)
		if err != nil {
			fmt.Println("parse proto failed:", err)
			break
		}
		fmt.Println("login resp:", pbr.GetErrorCode(), pbr.GetErrorMsg())

		time.Sleep(1e9)
		break
	}

}