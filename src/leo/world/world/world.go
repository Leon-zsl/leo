/* this is stage module
 */

package world

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

type World struct {
	running bool

	Port    *base.Port
	Service base.Service
}

var (
	Root *World = nil
)

func NewWorld() (world *World, err error) {
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)
	fmt.Println("number if cpu: ", cpu)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("create world failed!", r, string(debug.Stack()))
			world = nil
		}
	}()

	if Root != nil {
		world = Root
		return
	}

	Root = new(World)
	err = Root.init()
	if err != nil {
		fmt.Println("init world failed", err)
		debug.PrintStack()
		Root = nil
		return
	}

	world = Root
	return
}

func (world *World) init() error {
	//parse config file
	confile := path.Join(CONF_PATH, CONF_FILE)
	conf, err := ini.LoadFile(confile)
	if err != nil {
		world.close()
		return err
	}

	//init logger
	file, ok := conf.Get("logger", "config_file")
	if !ok {
		world.close()
		return errors.New("can not find logger/config_file in world config file")
	}
	ty, ok := conf.Get("logger", "log_type")
	if !ok {
		world.close()
		return errors.New("can not find logger/log_type in world config file")
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
		world.close()
		return err
	}

	//init port
	cid, ok := conf.Get("port_server", "id")
	if !ok {
		world.close()
		return errors.New("can not find port_server/id in world config file")
	}
	id, err := strconv.Atoi(cid)
	if err != nil {
		world.close()
		return err
	}
	ip, ok := conf.Get("port_server", "ip")
	if !ok {
		world.close()
		return errors.New("can not find port_server/ip in world config file")
	}
	pt, ok := conf.Get("port_server", "port")
	if !ok {
		world.close()
		return errors.New("can not find port_server/port in world config file")
	}
	port, err := strconv.Atoi(pt)
	if err != nil {
		world.close()
		return err
	}
	p, err := base.NewPort(id, ip, port)
	if err != nil {
		world.close()
		return err
	}
	world.Port = p

	//init service
	sv, err := NewWorldService()
	if err != nil {
		world.close()
		return err
	}
	world.Service = sv

	return nil
}

func (world *World) Run() {
	defer func() {
		if r := recover(); r != nil {
			if base.LoggerIns != nil {
				base.LoggerIns.Critical(r, string(debug.Stack()))
			} else {
				fmt.Println("runtine exception:", r, string(debug.Stack()))
			}
		}

		world.close()
	}()

	world.start()
	c := time.Tick(60 * time.Millisecond)
	for _ = range c {
		world.Service.Tick()
		if !world.running {
			break
		}
	}
	world.save()
}

func (world *World) Shutdown() {
	world.running = false
}

func (world *World) close() {
	world.running = false

	if world.Port != nil {
		world.Port.Close()
		world.Port = nil
	}
	if world.Service != nil {
		world.Service.Close()
		world.Service = nil
	}

	Root = nil
}

func (world *World) start() {
	world.running = true
	world.Port.Start()
	world.Service.Start()
}

func (world *World) save() {
	world.Service.Save()
}
