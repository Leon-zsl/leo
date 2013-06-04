/* this is the gameplay room
*/

package world

import (
	"uuid"
	"leo/base"
	"leo/world/world/module"
)

//one room one goroutine
type Room struct {
	id string
	actor_map map[string] *Actor
	mod_map map[string] module.Module

	dispatcher *base.Dispatcher
}

func NewRoom() (room *Room, err error) {
	room = new(Room)
	err = room.init()
	return
}

func (room *Room) init() error {
	room.id = uuid.New()
	room.dispatcher, _ = base.NewDispatcher()
	room.actor_map = make(map[string] *Actor)
	return nil
}

func (room *Room) Start() error {
	return nil
}

func (room *Room) Close() error {
	return nil
}

func (room *Room) Tick() error {
	return nil
}

func (room *Room) ID() string {
	return room.id
}

func (room *Room) Actor(id string) (*Actor, bool) {
	a, ok := room.actor_map[id]
	return a, ok
}

func (room *Room) AddActor(a *Actor) {
	if a == nil {
		return
	}
	room.actor_map[a.ID()] = a
}

func (room *Room) DelActor(a *Actor) {
	if a == nil {
		return
	}
	delete(room.actor_map, a.ID())
}

func (room *Room) DelActorByID(id string) {
	delete(room.actor_map, id)
}

func (room *Room) Module(id string) (module.Module, bool) {
	m, ok := room.mod_map[id]
	return m, ok
}

func (room *Room) AddModule(m module.Module) {
	if m == nil {
		return
	}
	room.mod_map[m.ID()] = m
}

func (room *Room) DelModule(m module.Module) {
	if m == nil {
		return
	}
	delete(room.mod_map, m.ID())
}

func (room *Room) DelModuleByID(id string) {
	delete(room.mod_map, id)
}