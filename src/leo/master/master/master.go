/* this is master module
*/

package master

import (
	"fmt"
	"path"
	"time"
	"strconv"
	"errors"
	"runtime"
	"runtime/debug"

	"ini"

	"leo/base"
)
type Master struct {
	running bool

	rpcserver *base.RpcServer

	Service *Service
}

var (
	Root *Master = nil
)

func NewMaster() (master *Master, err error) {
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)
	fmt.Println("number if cpu: ", cpu)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("create master failed!", r, string(debug.Stack()))
			master = nil
		}
	}()

	if Root != nil {
		master = Root
		return
	}

	Root = new(Master)
	err = Root.init()
	if err != nil {
		fmt.Println("init master failed", err)
		debug.PrintStack()
		Root = nil
		return
	}

	master = Root
	return
}

func (master *Master) init() error {
	//parse config file
	confile := path.Join(CONF_PATH, CONF_FILE)
	conf, err := ini.LoadFile(confile)
	if err != nil {
		master.close()
		return err
	}

	//init logger
	file, ok := conf.Get("logger", "config_file")
	if !ok {
		master.close()
		return errors.New("can not find logger/config_file in master config file")
	}
	ty, ok := conf.Get("logger", "log_type")
	if !ok {
		master.close()
		return errors.New("can not find logger/log_type in master config file")
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
		master.close()
		return err
	}

	//init rpc server
	ip, ok := conf.Get("rpc_server", "ip")
	if !ok {
		master.close()
		return errors.New("can not find rpc_server/ip in master config file")
	}
	pt, ok := conf.Get("rpc_server", "port")
	if !ok {
		master.close()
		return errors.New("can not find rpc_server/port in master config file")
	}
	port, err := strconv.Atoi(pt)
	if err != nil {
		master.close()
		return err
	}
	rpcs, err := base.NewRpcServer(ip, port)
	if err != nil {
		master.close()
		return err
	}
	master.rpcserver = rpcs

	//init service
	sv, err := NewService()
	if err != nil {
		master.close()
		return err
	}
	master.Service = sv

	return nil
}

func (master *Master) Run() {
	defer func() {
		if r := recover(); r != nil {
			if base.LoggerIns != nil {
				base.LoggerIns.Critical(r, string(debug.Stack()))
			} else {
				fmt.Println("runtine exception:", r, string(debug.Stack()))
			}
		}
		
		master.close()
	}()

	master.start()

	c := time.Tick(60 * time.Millisecond)
	for _ = range c {
		master.Service.Tick()
		if !master.running {
			break
		}
	}
}

func (master *Master) Shutdown() {
	master.running = false
}

func (master *Master) close() {
	master.running = false

	if master.rpcserver != nil {
		master.rpcserver.Close()
		master.rpcserver = nil
	}
	if master.Service != nil {
		master.Service.Close()
		master.Service = nil
	}

	Root = nil
}

func (master *Master) start() {
	master.running = true

	master.rpcserver.Start()

	master.Service.Start()
}

func (master *Master) save() {
	//todo:
}