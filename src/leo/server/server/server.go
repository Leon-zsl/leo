/* this is server 
*/

package server

import (
	"fmt"
	"net"
// 	"bytes"
// 	"encoding/binary"
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
	addr, err := net.ResolveTCPAddr("tcp", ":10001")
	if err != nil {
		fmt.Println("resolve addr err:", err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("listen err:", err)
	}
	
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("accept err:", err)
		}
		go app.handle_conn(conn)
	}
}

func (app *App) handle_conn(conn *net.TCPConn) {
	var buf [2048]byte
	for {
		rn, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println("read err:", err)
		}
		fmt.Println("read len:", rn)

		wn, err := conn.Write(buf[:rn])
		if err != nil {
			fmt.Println("write err:", err)
		}
		fmt.Println("wr len:", wn)
	}

	conn.Close()
}