/* this is the gate entry
 */
package gate

import (
	"fmt"
	"errors"
	"strconv"
	"path"
	"ini"
	"runtime"
	"runtime/debug"

	"leo/base"
)

type Gate struct {
	running bool
	cfgFile ini.File

	Logger *base.Logger
	Acceptor *Acceptor
	Connector *Connector
	SessionMgr *SessionMgr

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
	err := gate.parseConfig()
	if err != nil {
		gate.close()
		return err
	}

	//init logger
	file, ok := gate.cfgFile.Get("logger", "config_file")
	if !ok {
		gate.close()
		return errors.New("can not find logger/config_file in gate config file")
	}
	ty, ok := gate.cfgFile.Get("logger", "log_type")
	if !ok {
		gate.close()
		return errors.New("can not find logger/log_type in gate config file")
	}
	err = gate.createLogger(ty, path.Join(CONF_PATH, file))
	if err != nil {
		gate.close()
		return err
	}

	//init acceptor
	port, ok := gate.cfgFile.Get("acceptor", "port")
	if !ok {
		gate.close()
		return errors.New("can not find acceptor/port in gate config file")
	}
	val, err := strconv.Atoi(port)
	if err != nil {
		gate.close()
		return err
	}
	ip, ok := gate.cfgFile.Get("acceptor", "ip")
	if !ok {
		gate.close()
		return errors.New("can not find acceptor/ip in gate config file")
	}
	count, ok := gate.cfgFile.Get("acceptor", "count")
	if !ok {
		gate.close()
		return errors.New("can not find acceptor/count in gate config file")
	}
	cval, err := strconv.Atoi(count)
	if err != nil {
		gate.close()
		return err
	}
	ac, err := NewAcceptor(ip, val, cval)
	if err != nil {
		gate.close()
		return err
	}
	gate.Acceptor = ac

	//init connector
	conn, err := NewConnector()
	if err != nil {
		gate.close()
		return err
	}
	gate.Connector = conn

	//init session mgr
	sm, err := NewSessionMgr()
	if err != nil {
		gate.close()
		return err
	}
	gate.SessionMgr = sm
	
	//init service
	sv, err := NewService()
	if err != nil {
		gate.close()
		return err
	}
	gate.Service = sv

	return nil
}

func (gate *Gate) Start() {
	gate.running = true

	gate.Acceptor.Start()
	gate.Connector.Start()
	gate.SessionMgr.Start()
	gate.Service.Start()
}

func (gate *Gate) Run() {
	defer func() {
		if r := recover(); r != nil {
			if gate.Logger != nil {
				gate.Logger.Critical(r, string(debug.Stack()))
			} else {
				fmt.Println("run time exception:", r, string(debug.Stack()))
			}
		}

		gate.save()
		gate.close()
	}()

	for {
		if !gate.running {
			break
		}

		gate.SessionMgr.Update()
		gate.Service.Update()
	}
}

func (gate *Gate) Shutdown() {
	gate.running = false
}

func (gate *Gate) close() {
	if gate.Service != nil {
		gate.Service.Close()
		gate.Service = nil
	}
	if gate.SessionMgr != nil {
		gate.SessionMgr.Close()
		gate.SessionMgr = nil
	}
	if gate.Acceptor != nil {
		gate.Acceptor.Close()
		gate.Acceptor = nil
	}
	if gate.Connector != nil {
		gate.Connector.Close()
		gate.Connector = nil
	}
	if gate.Logger != nil {
		gate.Logger.Close()
		gate.Logger = nil
	}

	gate.cfgFile = nil
	Root = nil
}

func (gate *Gate) save() {
	//todo:
}

func (gate *Gate) parseConfig() error {
	confile := path.Join(CONF_PATH, CONF_FILE)
	conf, err := ini.LoadFile(confile)
	
	if err == nil {
		gate.cfgFile = conf
	}
	return err
}

func (gate *Gate) createLogger(ty, confile string) error {
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
		gate.Logger = lg
	}
	return err
}