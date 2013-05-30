/* this is db module
*/

package db

import (
	"fmt"
	"errors"
	"path"
	"ini"
	"runtime"
	"runtime/debug"

	"leo/base"
)

type DB struct {
	running bool
	cfgFile ini.File
	
	Logger *base.Logger
	Driver *Driver
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
	err := db.parse_config()
	if err != nil {
		db.close()
		return err
	}

	//init logger
	file, ok := db.cfgFile.Get("logger", "config_file")
	if !ok {
		db.close()
		return errors.New("can not find logger/config_file in db config file")
	}
	ty, ok := db.cfgFile.Get("logger", "log_type")
	if !ok {
		db.close()
		return errors.New("can not find logger/log_type in db config file")
	}
	err = db.create_logger(ty, path.Join(CONF_PATH, file))
	if err != nil {
		db.close()
		return err
	}

	db_addr, ok := db.cfgFile.Get("db", "host")
	if !ok {
		db.close()
		return errors.New("can not find db/host in db config file")
	}
	db_name, ok := db.cfgFile.Get("db", "database")
	if !ok {
		db.close()
		return errors.New("can not find db/database in db config file")
	}
	db_account, ok := db.cfgFile.Get("db", "username")
	if !ok {
		db.close()
		return errors.New("can not find db/username in db config file")
	}
	db_pwd, ok := db.cfgFile.Get("db", "password")
	if !ok {
		db.close()
		return errors.New("can not find db/password in db config file")
	}
	dr, err := NewDriver(db_addr, db_name, db_account, db_pwd)
	if err != nil {
		db.close()
		return err
	}
	db.Driver = dr

	//todo:

	return nil
}

func (db *DB) Start() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("start exception:", r, string(debug.Stack()))
			db.close()
		}
	}()

	db.running = true

	db.Driver.Start()
	//todo:
}

func (db *DB) Run() {
	defer func() {
		if r := recover(); r != nil {
			if db.Logger != nil {
				db.Logger.Critical(r, string(debug.Stack()))
			} else {
				fmt.Println("run time exception:", r, string(debug.Stack()))
			}
		}

		db.save()
		db.close()
	}()

	for {
		if !db.running {
			break
		}

		db.Driver.Update()
		//todo:
	}
}

func (db *DB) Shutdown() {
	db.running = false
}

func (db *DB) save() {
	//todo:
}

func (db *DB) close() {
	//todo:
	db.running = false

	if db.Driver != nil {
		db.Driver.Close()
		db = nil
	}
	if db.Logger != nil {
		db.Logger.Close()
		db.Logger = nil
	}
	db.cfgFile = nil
	Root = nil
}

func (db *DB) parse_config() error {
	confile := path.Join(CONF_PATH, CONF_FILE)
	conf, err := ini.LoadFile(confile)
	
	if err == nil {
		db.cfgFile = conf
	}
	return err
}

func (db *DB) create_logger(ty, confile string) error {
	v := base.LOG_TYPE_SYS
	switch ty {
	case "sys":
		v = base.LOG_TYPE_SYS
	case "log4go":
		v = base.LOG_TYPE_LOG4GO
	default:
		fmt.Println("invalid log type", ty)
	}
	lg, err := base.NewLogger(v, confile)
	if err == nil {
		db.Logger = lg
	}
	return err
}
