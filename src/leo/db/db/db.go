/* this is db module
 */

package db

import (
	"errors"
	"fmt"
	"path"
	//	"ini"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"leo/base"
)

import ini "github.com/vaughan0/go-ini"

type DB struct {
	running bool

	Driver  *base.Driver
	Port    *base.Port
	Service base.Service
}

var (
	Root *DB = nil
)

func NewDB() (db *DB, err error) {
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)
	fmt.Println("number if cpu: ", cpu)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("create db failed!", r, string(debug.Stack()))
			db = nil
		}
	}()

	if Root != nil {
		db = Root
		return
	}

	Root = new(DB)
	err = Root.init()
	if err != nil {
		fmt.Println("init db failed", err)
		debug.PrintStack()
		Root = nil
		return
	}

	db = Root
	return
}

func (db *DB) init() error {
	//parse config file
	confile := path.Join(CONF_PATH, CONF_FILE)
	conf, err := ini.LoadFile(confile)
	if err != nil {
		db.close()
		return err
	}

	//init logger
	file, ok := conf.Get("logger", "config_file")
	if !ok {
		db.close()
		return errors.New("can not find logger/config_file in db config file")
	}
	ty, ok := conf.Get("logger", "log_type")
	if !ok {
		db.close()
		return errors.New("can not find logger/log_type in db config file")
	}
	v := base.LOG_TYPE_SYS
	switch ty {
	case "sys":
		v = base.LOG_TYPE_SYS
	case "log4go":
		v = base.LOG_TYPE_LOG4GO
	default:
		fmt.Println("invalid log type", ty)
	}
	_, err = base.NewLogger(v, path.Join(CONF_PATH, file))
	if err != nil {
		db.close()
		return err
	}

	//init db driver
	db_addr, ok := conf.Get("db", "host")
	if !ok {
		db.close()
		return errors.New("can not find db/host in db config file")
	}
	db_name, ok := conf.Get("db", "database")
	if !ok {
		db.close()
		return errors.New("can not find db/database in db config file")
	}
	db_account, ok := conf.Get("db", "username")
	if !ok {
		db.close()
		return errors.New("can not find db/username in db config file")
	}
	db_pwd, ok := conf.Get("db", "password")
	if !ok {
		db.close()
		return errors.New("can not find db/password in db config file")
	}
	cache, ok := conf.Get("db", "cache")
	if !ok {
		db.close()
		return errors.New("can not find db/cache in db config file")
	}
	db_cache := false
	if cache == "true" {
		db_cache = true
	}
	dr, err := base.NewDriver(db_addr, db_name, db_account, db_pwd, db_cache)
	if err != nil {
		db.close()
		return err
	}
	db.Driver = dr

	//init port
	cid, ok := conf.Get("port_server", "id")
	if !ok {
		db.close()
		return errors.New("can not find port_server/id in db config file")
	}
	id, err := strconv.Atoi(cid)
	if err != nil {
		db.close()
		return err
	}
	ip, ok := conf.Get("port_server", "ip")
	if !ok {
		db.close()
		return errors.New("can not find port_server/ip in db config file")
	}
	pt, ok := conf.Get("port_server", "port")
	if !ok {
		db.close()
		return errors.New("can not find port_server/port in db config file")
	}
	port, err := strconv.Atoi(pt)
	if err != nil {
		db.close()
		return err
	}
	p, err := base.NewPort(id, ip, port)
	if err != nil {
		db.close()
		return err
	}
	db.Port = p

	//init rpc service
	err = BuildRpcService(p)
	if err != nil {
		db.close()
		return err
	}

	//init service
	sv, err := NewDBService()
	if err != nil {
		db.close()
		return err
	}
	db.Service = sv

	return nil
}

func (db *DB) Run() {
	defer func() {
		if r := recover(); r != nil {
			if base.LoggerIns != nil {
				base.LoggerIns.Critical(r, string(debug.Stack()))
			} else {
				fmt.Println("run time exception:", r, string(debug.Stack()))
			}
		}
		db.close()
	}()

	db.start()
	c := time.Tick(60 * time.Millisecond)
	for _ = range c {
		db.Service.Tick()
		if !db.running {
			break
		}
	}
	db.save()
}

func (db *DB) Shutdown() {
	db.running = false
}

func (db *DB) save() {
	db.Service.Save()
}

func (db *DB) start() {
	db.running = true

	db.Driver.Start()
	db.Port.Start()
	db.Service.Start()
}

func (db *DB) close() {
	db.running = false

	if db.Service != nil {
		db.Service.Close()
		db = nil
	}
	if db.Port != nil {
		db.Port.Close()
		db.Port = nil
	}
	if db.Driver != nil {
		db.Driver.Close()
		db = nil
	}
	if base.LoggerIns != nil {
		base.LoggerIns.Close()
		base.LoggerIns = nil
	}

	Root = nil
}
