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

	"leo/base/leoerror"
)

type Gate struct {
	running bool

	cfgFile ini.File

	Logger log4go.Logger

	SessionMgr *SessionMgr
	Acceptor *Acceptor
}

var (
	Root *Gate
)

func NewGate() (gt *Gate, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("create gate failed!")
			gt = nil
			err = leoerror.CreateLeoError(leoerror.ErrStartFailed, 
				"gate start failed")
		}
	}()

	gt = new(Gate)
	err = gt.Startup()
	if err != nil {
		fmt.Println("init gate failed")
		gt = nil
		return
	}

	Root = gt
	return
}

func (gate *Gate) Startup() error {
	gate.parseConfig()

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

	port, ok := gate.cfgFile.Get("listen_addr", "port")
	if !ok {
		gate.close()
		return errors.New("can not find listen_addr/port in gate config file")
	}
	val, err := strconv.Atoi(port)
	if err != nil {
		gate.close()
		return errors.New("listen port is invalid")
	}
	ac, err := NewAcceptor(val)
	if err != nil {
		gate.close()
		return err
	}
	gate.Acceptor = ac

	sm, err := NewSessionMgr()
	if err != nil {
		gate.close()
		return err
	}
	gate.SessionMgr = sm

	gate.running = true
	return nil
}

func (gate *Gate) Run() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("runtime exception catched!")

			err := leoerror.CreateLeoError(leoerror.ErrRuntimeExcept, 
				"runtime exception failed")
			gate.Logger.Error(err)
		}

		gate.save()
		gate.close()
	}()

	for gate.running {
		//todo:
	}
}

func (gate *Gate) Shutdown() {
	gate.running = false
}

func (gate *Gate) close() {
	if gate.Logger != nil {
		gate.Logger.Close()
		gate.Logger = nil
	}
	if gate.Acceptor != nil {
		gate.Acceptor.Close()
		gate.Acceptor = nil
	}
	if gate.SessionMgr != nil {
		gate.SessionMgr.Close()
		gate.SessionMgr = nil
	}

	gate.cfgFile = nil
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