/* this is the gameplay room
 */

package world

import (
	"fmt"
	"sync"
	"time"
	//	"uuid"
	"leo/base"
	"leo/world/world/module"
	"runtime/debug"
)

import uuid "code.google.com/p/go-uuid/uuid"

//one room one goroutine
type Room struct {
	running bool
	id      string

	lock_actor sync.Mutex
	actor_map  map[string]*Actor

	lock_mod sync.Mutex
	mod_map  map[string]module.Module

	dispatcher *base.Dispatcher
}

func NewRoom() (room *Room, err error) {
	room = new(Room)
	err = room.init()
	return
}

func (room *Room) init() error {
	room.running = false
	room.id = uuid.New()
	room.actor_map = make(map[string]*Actor)
	room.mod_map = make(map[string]module.Module)
	room.dispatcher, _ = base.NewDispatcher()
	return nil
}

func (room *Room) Start() error {
	room.running = true
	go room.run()
	return nil
}

func (room *Room) Close() error {
	room.running = false
	return nil
}

func (room *Room) Tick() error {
	//do nothing now
	// handle in goroutine
	return nil
}

func (room *Room) Save() error {
	return nil
}

func (room *Room) ID() string {
	return room.id
}

func (room *Room) Actor(id string) (*Actor, bool) {
	//do not need lock for read-only op
	a, ok := room.actor_map[id]
	return a, ok
}

func (room *Room) AddActor(a *Actor) {
	if a == nil {
		return
	}
	room.lock_actor.Lock()
	defer room.lock_actor.Unlock()
	room.actor_map[a.ID()] = a
}

func (room *Room) DelActor(a *Actor) {
	if a == nil {
		return
	}
	room.lock_actor.Lock()
	defer room.lock_actor.Unlock()
	delete(room.actor_map, a.ID())
}

func (room *Room) DelActorByID(id string) {
	room.lock_actor.Lock()
	defer room.lock_actor.Unlock()
	delete(room.actor_map, id)
}

func (room *Room) Module(id string) (module.Module, bool) {
	//do not need lock for read-only op
	m, ok := room.mod_map[id]
	return m, ok
}

func (room *Room) AddModule(m module.Module) {
	if m == nil {
		return
	}
	room.lock_mod.Lock()
	defer room.lock_mod.Unlock()
	room.mod_map[m.ID()] = m
}

func (room *Room) DelModule(m module.Module) {
	if m == nil {
		return
	}
	room.lock_mod.Lock()
	defer room.lock_mod.Unlock()
	delete(room.mod_map, m.ID())
}

func (room *Room) DelModuleByID(id string) {
	room.lock_mod.Lock()
	defer room.lock_mod.Unlock()
	delete(room.mod_map, id)
}

func (room *Room) run() {
	defer func() {
		if r := recover(); r != nil {
			if base.LoggerIns != nil {
				base.LoggerIns.Critical(r, string(debug.Stack()))
			} else {
				fmt.Println("runtine exception:", r, string(debug.Stack()))
			}
		}
		room.doclose()
	}()

	c := time.Tick(60 * time.Millisecond)
	for _ = range c {
		room.lock_mod.Lock()
		for _, m := range room.mod_map {
			m.Tick()
		}
		room.lock_mod.Unlock()

		room.lock_actor.Lock()
		for _, a := range room.actor_map {
			a.Tick()
		}
		room.lock_actor.Unlock()

		if !room.running {
			break
		}
	}

	room.Save()
}

func (room *Room) doclose() {
	room.running = false
	room.lock_mod.Lock()
	for _, m := range room.mod_map {
		m.Close()
	}
	room.lock_mod.Unlock()

	room.lock_actor.Lock()
	for _, a := range room.actor_map {
		a.Close()
	}
	room.lock_actor.Unlock()
}
