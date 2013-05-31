/* this is db driver
*/

package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"leo/base"
)

type Driver struct {
	running bool

	addr string
	name string
	account string
	pwd string

	db *sql.DB
}

func NewDriver(addr, name, account, pwd string) (driver *Driver, err error) {
	driver = new(Driver)
	err = driver.init(addr, name, account, pwd)
	return
}

func (driver *Driver) init(addr, name, account, pwd string) error {
	driver.addr = addr
	driver.name = name
	driver.account = account
	driver.pwd = pwd
	return nil
}

func (driver *Driver) Start() {
	db, err := sql.Open("mysql", driver.account + ":" + driver.pwd + 
		"@tcp" + "(" + driver.addr + ")" + 
		"/" + driver.name + 
		"?" + "charset=" + "utf8")

	if err != nil {
		base.LoggerIns.Critical(err)
		return
	}

	driver.db = db
	driver.running = true
}

func (driver *Driver) Close() {
	driver.running = false
	if driver.db != nil {
		driver.db.Close()
		driver.db = nil
	}
}

func (driver *Driver) Exec(sql string) {
}

func (driver *Driver) Query(sql string) {
}
