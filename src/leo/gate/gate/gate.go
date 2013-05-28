/* this is the gate entry
 */
package gate

import (
	"fmt"
	"errors"
	"strconv"
	"path"
	"ini"
	"log4go"

	"leo/base"
)

type Gate struct {
	running bool

	cfgFile ini.File

	Logger log4go.Logger

	SessionMgr *SessionMgr
	Acceptor *Acceptor
}

var (
	Root *Gate = nil
)

func NewGate() (gt *Gate, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("create gate failed!")
			gt = nil
			err = base.NewLeoError(base.LeoErrStartFailed, 
				"gate start failed")
		}
	}()

	err = nil
	gt = nil

	if Root != nil {
		gt = Root
		return
	}

	Root = new(Gate)
	err = Root.init()
	if err != nil {
		fmt.Println("init gate failed")
		Root = nil
		return
	}

	gt = Root
	return
}

func (gate *Gate) init() error {
	//parse config file
	gate.parseConfig()

	//init logger
	file, ok := gate.cfgFile.Get("logger", "config_file")
	if !ok {
		gate.close()
		return errors.New("can not find logger/config_file in gate config file")
	}
	name, ok := gate.cfgFile.Get("logger", "log_name")
	if !ok {
		gate.close()
		return errors.New("can not find logger/log_name in gate config file")
	}
	err := gate.createLogger(name, path.Join(CONF_PATH, file))
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
		return errors.New("listen port is invalid")
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
		return errors.New("listen count is invalid")
	}
	ac, err := NewAcceptor(ip, val, cval)
	if err != nil {
		gate.close()
		return err
	}
	gate.Acceptor = ac

	//init session mgr
	sm, err := NewSessionMgr()
	if err != nil {
		gate.close()
		return err
	}
	gate.SessionMgr = sm

	//init connector
	//todo:
	
	gate.running = true
	return nil
}

func (gate *Gate) Start() {
	gate.Logger.Info("gate start now")
	gate.Acceptor.Start()
	gate.SessionMgr.Start()
}

func (gate *Gate) Run() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("runtime exception catched!")
			gate.Logger.Critical(r)
		}

		gate.save()
		gate.close()
	}()

//	for gate.running {
		gate.Acceptor.Update()
		gate.SessionMgr.Update()
//	}
}

func (gate *Gate) Shutdown() {
	gate.running = false
}

func (gate *Gate) close() {
	gate.Logger.Info("gate close now")

	if gate.SessionMgr != nil {
		gate.SessionMgr.Close()
		gate.SessionMgr = nil
	}
	if gate.Acceptor != nil {
		gate.Acceptor.Close()
		gate.Acceptor = nil
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

func (gate *Gate) createLogger(logname, confile string) error {
	log := make(log4go.Logger)
	log.LoadConfiguration(confile)
	
	filter, ok := log[logname]
	if ok {
		if IS_DEBUG {
			filter.Level = log4go.DEBUG
		} else {
			filter.Level = log4go.ERROR
		}
	} else {
		return errors.New("can not find log name " + logname)
	}

	gate.Logger = log
	return nil
}