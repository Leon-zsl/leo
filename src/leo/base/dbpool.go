/* this is db pool */

package base

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DBPool struct {
	running bool

	addr    string
	name    string
	account string
	pwd     string

	size  int
	conns chan *sql.DB
}

func NewDBPool(size int, addr, name, account, pwd string) (pool *DBPool, err error) {
	pool = new(DBPool)
	err = pool.init(size, addr, name, account, pwd)
	return
}

func (pool *DBPool) init(size int, addr, name, account, pwd string) error {
	pool.size = size
	pool.addr = addr
	pool.name = name
	pool.account = account
	pool.pwd = pwd
	pool.conns = make(chan *sql.DB, size)
	return nil
}

func (pool *DBPool) Start() error {
	pool.running = true
	go pool.product()
	return nil
}

func (pool *DBPool) Close() error {
	pool.running = false
	l := len(pool.conns)
	for i := 0; i < l; i++ {
		conn := <-pool.conns
		conn.Close()
	}
	close(pool.conns)
	return nil
}

func (pool *DBPool) Get() *sql.DB {
	return <-pool.conns
}

func (pool *DBPool) product() {
	go func() {
		if r := recover(); r != nil {
			if LoggerIns != nil {
				LoggerIns.Critical(r)
			} else {
				fmt.Println(r)
			}
		}
	}()

	for pool.running {
		conn, err := sql.Open("mysql", pool.account+":"+pool.pwd+
			"@tcp"+"("+pool.addr+")"+
			"/"+pool.name+
			"?"+"charset="+"utf8")
		if err != nil {
			LoggerIns.Error("create db conn failed", err)
			break
		}

		if !pool.running {
			conn.Close()
			break
		}
		pool.conns <- conn
	}
}
