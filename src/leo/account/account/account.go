/* this is stage module
 */

package account

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	//	"ini"

	"leo/base"
)

import ini "github.com/vaughan0/go-ini"

type Account struct {
	running bool

	Port    *base.Port
	Service base.Service
}

var (
	Root *Account = nil
)

func NewAccount() (account *Account, err error) {
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)
	fmt.Println("number if cpu: ", cpu)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("create account failed!", r, string(debug.Stack()))
			account = nil
		}
	}()

	if Root != nil {
		account = Root
		return
	}

	Root = new(Account)
	err = Root.init()
	if err != nil {
		fmt.Println("init account failed", err)
		debug.PrintStack()
		Root = nil
		return
	}

	account = Root
	return
}

func (account *Account) init() error {
	//parse config file
	confile := path.Join(CONF_PATH, CONF_FILE)
	conf, err := ini.LoadFile(confile)
	if err != nil {
		account.close()
		return err
	}

	//init logger
	file, ok := conf.Get("logger", "config_file")
	if !ok {
		account.close()
		return errors.New("can not find logger/config_file in account config file")
	}
	ty, ok := conf.Get("logger", "log_type")
	if !ok {
		account.close()
		return errors.New("can not find logger/log_type in account config file")
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
		account.close()
		return err
	}

	//init port
	cid, ok := conf.Get("port_server", "id")
	if !ok {
		account.close()
		return errors.New("can not find port_server/id in account config file")
	}
	id, err := strconv.Atoi(cid)
	if err != nil {
		account.close()
		return err
	}
	ip, ok := conf.Get("port_server", "ip")
	if !ok {
		account.close()
		return errors.New("can not find port_server/ip in account config file")
	}
	pt, ok := conf.Get("port_server", "port")
	if !ok {
		account.close()
		return errors.New("can not find port_server/port in account config file")
	}
	port, err := strconv.Atoi(pt)
	if err != nil {
		account.close()
		return err
	}
	p, err := base.NewPort(id, ip, port)
	if err != nil {
		account.close()
		return err
	}
	account.Port = p

	//rpc service
	err = BuildRpcService(p)
	if err != nil {
		account.close()
		return err
	}

	//init service
	sv, err := NewAccountService()
	if err != nil {
		account.close()
		return err
	}
	account.Service = sv

	return nil
}

func (account *Account) Run() {
	defer func() {
		if r := recover(); r != nil {
			if base.LoggerIns != nil {
				base.LoggerIns.Critical(r, string(debug.Stack()))
			} else {
				fmt.Println("runtine exception:", r, string(debug.Stack()))
			}
		}

		account.close()
	}()

	account.start()
	c := time.Tick(60 * time.Millisecond)
	for _ = range c {
		account.Service.Tick()
		if !account.running {
			break
		}
	}
	account.save()
}

func (account *Account) Shutdown() {
	account.running = false
}

func (account *Account) close() {
	account.running = false

	if account.Port != nil {
		account.Port.Close()
		account.Port = nil
	}
	if account.Service != nil {
		account.Service.Close()
		account.Service = nil
	}

	Root = nil
}

func (account *Account) start() {
	account.running = true
	account.Port.Start()
	account.Service.Start()
}

func (account *Account) save() {
	account.Service.Save()
}
