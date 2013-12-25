/* this is event dispatcher
 */

package base

import (
	"container/list"
	"sync"
)

//goroutine safe
type Dispatcher struct {
	lock_e     sync.Mutex
	event_list *list.List

	lock_h      sync.Mutex
	handler_map map[int][]EventHandler
}

func NewDispatcher() (disp *Dispatcher, err error) {
	disp = new(Dispatcher)
	err = disp.init()
	return
}

func (disp *Dispatcher) init() error {
	disp.handler_map = make(map[int][]EventHandler)
	disp.event_list = list.New()
	return nil
}

func (disp *Dispatcher) RegisterHandler(op int, h EventHandler) {
	if h == nil {
		return
	}
	if _, ok := disp.handler_map[op]; !ok {
		disp.handler_map[op] = make([]EventHandler, 0)
	}
	disp.handler_map[op] = append(disp.handler_map[op], h)
}

func (disp *Dispatcher) UnRegisterHandler(op int, h EventHandler) {
	if h == nil {
		return
	}
	if l, ok := disp.handler_map[op]; ok {
		for i, v := range l {
			if v == h {
				disp.handler_map[op] = append(l[:i], l[i+1:]...)
				break
			}
		}
	}
}

func (disp *Dispatcher) SendEvent(e Event) {
	if e == nil {
		return
	}
	disp.lock_e.Lock()
	defer disp.lock_e.Unlock()
	disp.event_list.PushBack(e)
}

func (disp *Dispatcher) Dispatch() {
	disp.lock_e.Lock()
	defer disp.lock_e.Unlock()

	disp.lock_h.Lock()
	defer disp.lock_h.Unlock()

	for elem := disp.event_list.Front(); elem != nil; elem = elem.Next() {
		if e, ok := elem.Value.(Event); ok {
			if l, ok := disp.handler_map[e.Op()]; ok {
				if ok {
					for _, h := range l {
						if e.Handled() {
							break
						}
						h.HandleEvent(e)
					}
				}
			}
		}
	}
	disp.event_list.Init()
}
