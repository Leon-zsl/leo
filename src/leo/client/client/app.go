/* this is the app
*/

package client

import (
	"fmt"
	"net"
	"bytes"
	"encoding/binary"
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

func (app *App) Startup() {
	con, err := net.Dial("tcp", "127.0.0.1:10001")
	if err != nil {
		fmt.Println("net dial failed:", err.Error())
		return
	}

	defer con.Close()
	app.Conn = con.(*net.TCPConn)

	for {
		snd_data := []byte("this is client")

		var cmd int32 = 16
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.BigEndian, cmd)
		snd_cmd := buf.Bytes()

		ln := int32(len(snd_cmd) + len(snd_data))
		buf = new(bytes.Buffer)
		binary.Write(buf, binary.BigEndian, ln)
		snd_ln := buf.Bytes()

		fmt.Println("send: ", ln, snd_ln, snd_cmd, snd_data)
		
		val := append(append(snd_ln, snd_cmd...), snd_data...)
		_, err := con.Write(val)
		if err != nil {
			fmt.Println("write err:", err.Error())
			continue
		}

		tmp := make([]byte, 512)
		n, err := con.Read(tmp)
		if err != nil {
			fmt.Println("read err:", err.Error())
			continue
		}

		var lr int32 = 0
		buf = bytes.NewBuffer(tmp[:4])
		binary.Read(buf, binary.BigEndian, &lr)
		
		var lc int32 = 0
		buf = bytes.NewBuffer(tmp[4:8])
		binary.Read(buf, binary.BigEndian, &lc)

		rcv_data := string(tmp[8:n])
		fmt.Println("rcv: ", lr, lc, rcv_data)
	}
}