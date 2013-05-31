/* this is the gate entry
 */
package gate

import (
	"fmt"
	"errors"
	"strconv"
	"path"
	"ini"
	"time"
	"runtime"
	"runtime/debug"

	"leo/base"
)

type Gate struct {
	running bool

	Acceptor *base.Acceptor
	Connector *base.Connector

	Service *GateService
}

var (
	Root *Gate = nil
)

func NewGate() (gt *Gate, err error) {
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)
	fmt.Println("number if cpu: ", cpu)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("create gate failed!", r, string(debug.Stack()))
			gt = nil
		}
	}()

	if Root != nil {
		gt = Root
		return
	}

	Root = new(Gate)
	err = Root.init()
	if err != nil {
		fmt.Println("init gate failed", err)
		debug.PrintStack()
		Root = nil
		return
	}

	gt = Root
	return
}

func (gate *Gate) init() error {
	//parse config file
	confile := path.Join(CONF_PATH, CONF_FILE)
	conf, err := ini.LoadFile(confile)
	if err != nil {
		gate.close()
		return err
	}

	//init logger
	file, ok := conf.Get("logger", "config_file")
	if !ok {
		gate.close()
		return errors.New("can not find logger/config_file in gate config file")
	}
	ty, ok := conf.Get("logger", "log_type")
	if !ok {
		gate.close()
		return errors.New("can not find logger/log_type in gate config file")
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
		gate.close()
		return err
	}

	//init acceptor
	port, ok := conf.Get("acceptor", "port")
	if !ok {
		gate.close()
		return errors.New("can not find acceptor/port in gate config file")
	}
	val, err := strconv.Atoi(port)
	if err != nil {
		gate.close()
		return err
	}
	ip, ok := conf.Get("acceptor", "ip")
	if !ok {
		gate.close()
		return errors.New("can not find acceptor/ip in gate config file")
	}
	count, ok := conf.Get("acceptor", "count")
	if !ok {
		gate.close()
		return errors.New("can not find acceptor/count in gate config file")
	}
	cval, err := strconv.Atoi(count)
	if err != nil {
		gate.close()
		return err
	}
	ac, err := base.NewAcceptor(ip, val, cval)
	if err != nil {
		gate.close()
		return err
	}
	gate.Acceptor = ac

	//init connector
	co, err := base.NewConnector()
	if err != nil {
		gate.close()
		return err
	}
	gate.Connector = co

	//init service
	sv, err := NewService()
	if err != nil {
		gate.close()
		return err
	}
	gate.Service = sv

	return nil
}

func (gate *Gate) Run() {
	defer func() {
		if r := recover(); r != nil {
			if base.LoggerIns != nil {
				base.LoggerIns.Critical(r, string(debug.Stack()))
			} else {
				fmt.Println("run time exception:", r, string(debug.Stack()))
			}
		}
		//gate.save()
		gate.close()
	}()

	gate.start()

	c := time.Tick(60 * time.Millisecond)
	for _ = range c {
		gate.Service.Tick()
		if !gate.running {
			break
		}
	}

// 	for {
// 		if !gate.running {
// 			break
// 		}
// 		gate.Service.Update()

// 		time.Sleep(1e6 * 30)
// 		fmt.Printf("%v\n", time.Now())
// 	}
}

func (gate *Gate) Shutdown() {
	gate.running = false
}

func (gate *Gate) start() {
	gate.running = true

	gate.Acceptor.Start()
	gate.Connector.Start()

	gate.Service.Start()
}


func (gate *Gate) close() {
	gate.running = false

	if gate.Service != nil {
		gate.Service.Close()
		gate.Service = nil
	}

	if gate.Acceptor != nil {
		gate.Acceptor.Close()
		gate.Acceptor = nil
	}
	if gate.Connector != nil {
		gate.Connector.Close()
		gate.Connector = nil
	}
	if base.LoggerIns != nil {
		base.LoggerIns.Close()
		base.LoggerIns = nil
	}

	Root = nil
}

func (gate *Gate) save() {
	//todo:
}