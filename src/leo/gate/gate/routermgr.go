/* this is session acceptor */

package gate

import (
	"sync"

	"leo/base"
)

type RouterMgr struct {
	lock sync.RWMutex
	rt_map map[string] *Router
}

func NewRouterMgr() (mgr *RouterMgr, err error) {
	mgr = new(RouterMgr)
	err = mgr.init()
	return
}

func (mgr *RouterMgr) init() error {
	mgr.rt_map = make(map[string]*Router)
	return nil
}

func (mgr *RouterMgr) Close() error {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	for _, v := range(mgr.rt_map) {
		v.Session().Close()
	}
	mgr.rt_map = make(map[string]*Router)
	return nil
}

func (mgr *RouterMgr) HandleAcceptedSession(ssn *base.Session) {
	rt, err := NewRouter(mgr, ssn)
	if err != nil {
		base.LoggerIns.Error("create router failed:", err)
		return
	}

	mgr.lock.Lock()
	mgr.rt_map[ssn.SID()] = rt
	mgr.lock.Unlock()
}

func (mgr *RouterMgr) Router(sid string) (*Router, bool) {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	rt, ok := mgr.rt_map[sid]
	return rt, ok
}

func (mgr *RouterMgr) CloseRouter(sid string) {
	mgr.lock.Lock()
	rt, ok := mgr.rt_map[sid]
	if !ok {
		return
	}
	delete(mgr.rt_map, sid)
	mgr.lock.Unlock()

	rt.Session().Close()
}

//do not close the ssn
func (mgr *RouterMgr) DelRouter(sid string) {
	mgr.lock.Lock()
	delete(mgr.rt_map, sid)
	mgr.lock.Unlock()
}

func (mgr *RouterMgr) Routers() []*Router {
	rts := make([]*Router, 0)
	mgr.lock.RLock()
	for _, v := range(mgr.rt_map) {
		rts = append(rts, v)
	}
	mgr.lock.RUnlock()
	return rts
}