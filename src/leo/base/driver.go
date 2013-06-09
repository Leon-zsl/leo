/* this is db driver
*/

package base

import (
	"strconv"
	"errors"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//goroutine safe
type Driver struct {
	running bool

	addr string
	name string
	account string
	pwd string

	db *sql.DB

	usecache bool
	cache *Cache
}

func NewDriver(addr, name, account, pwd string, usecache bool) (driver *Driver, err error) {
	driver = new(Driver)
	err = driver.init(addr, name, account, pwd, usecache)
	return
}

func (driver *Driver) init(addr, name, account, pwd string, usecache bool) error {
	driver.addr = addr
	driver.name = name
	driver.account = account
	driver.pwd = pwd
	driver.usecache = usecache

	driver.cache, _ = NewCache()
	return nil
}

func (driver *Driver) Start() error {
	db, err := sql.Open("mysql", driver.account + ":" + driver.pwd + 
		"@tcp" + "(" + driver.addr + ")" + 
		"/" + driver.name + 
		"?" + "charset=" + "utf8")

	if err != nil {
		return err
	}

	db.SetMaxIdleConns(32)
	driver.db = db
	driver.running = true
	driver.cache.Start()
	return nil
}

func (driver *Driver) Close() error {
	driver.running = false
	if driver.cache != nil {
		driver.cache.Close()
		driver.cache = nil
	}
	if driver.db != nil {
		driver.db.Close()
		driver.db = nil
	}

	return nil
}

func (driver *Driver) DB() *sql.DB {
	return driver.db
}

func (driver *Driver) Get(table string, key int, keyname string) (*Record, error){
	if table == "" {
		return nil, errors.New("table is invalid")
	}

	if driver.usecache {
		val := driver.cache.Get(table, key)
		if val != nil {
			return val, nil
		}
	}
	
	sql := "SELECT * FROM " + table + " WHERE " + keyname + "=?"
	rows, err := driver.db.Query(sql, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rcd, err := NewRecord()
	if err != nil {
		return nil, err
	}
	
	for rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		}

		names, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		val_addrs := make([]interface{}, len(names))
		for i, _ := range(names) {
			var v interface{} = nil
			val_addrs[i] = &v
		}
		err = rows.Scan(val_addrs...)
		if err != nil {
			return nil, err
		}

		for i, v := range(val_addrs) {
			rcd.SetValue(names[i], *(v.(*interface{})))
		}

		//wo only need the 1st line
		break
	}

	if driver.usecache {
		driver.cache.Add(table, key, rcd)
	}

	return rcd, nil
}

func (driver *Driver) Set(table string, key int, keyname string, record *Record) error {
	sql := "UPDATE " + table + " SET "
	idx := 0
	for _, name := range(record.Names) {
		comma := ","
		if idx == len(record.Names) - 1 {
			comma = ""
		}

		sql += name + "=?" + comma + " "
		idx++
	}
	
	sql += " WHERE " + keyname + "=" + strconv.Itoa(key)
	_, err := driver.db.Exec(sql, record.Values...)
	if err != nil {
		return err
	}

	if driver.usecache {
		driver.cache.Set(table, key, record)
	}

	return nil
}

func (driver *Driver) Add(table string, key int, keyname string, record *Record) error {
	sql := "INSERT INTO " + table + " VALUES " + "("
	idx := 0
	for _, _ = range(record.Names) {
		comma := ","
		if idx == len(record.Names) - 1 {
			comma = ""
		}

		sql += "?" + comma + " "
		idx++
	}
	sql += ")"

	_, err := driver.db.Exec(sql, record.Values...)
	if err != nil {
		return err
	}

	if driver.usecache {
		driver.cache.Add(table, key, record)
	}

	return nil
}

func (driver *Driver) Del(table string, key int, keyname string) error {
	sql := "DELETE FROM " + table + " WHERE " + keyname + "=?"
	_, err := driver.db.Exec(sql, key)
	if err != nil {
		return err
	}

	if driver.usecache {
		driver.cache.Del(table, key)
	}

	return nil
}

func (driver *Driver) Query(sql string, args []interface{}) ([]*Record, error) {
	rows, err := driver.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rcds := make([]*Record, 0)
	names, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		}

		rcd, _ := NewRecord()
		val_addrs := make([]interface{}, len(names))
		for i, _ := range(names) {
			var v interface{} = nil
			val_addrs[i] = &v
		}
		err = rows.Scan(val_addrs...)
		if err != nil {
			return nil, err
		}

		for i, v := range(val_addrs) {
			rcd.SetValue(names[i], *(v.(*interface{})))
		}

		rcds = append(rcds, rcd)
	}

	return rcds, nil
}

func (driver *Driver) Count(table string) int {
	rows, err := driver.db.Query("select * from " + table)
	if err != nil {
		return 0
	}

	count := 0
	for rows.Next() {
		count++
	}
	return count
}