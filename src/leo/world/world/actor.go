/* this is actor
 */

package world

import (
	"leo/world/world/comp"
)

type Actor struct {
	id       string
	comp_map map[string]comp.Component
}

func NewActor() (actor *Actor, err error) {
	actor = new(Actor)
	err = actor.init()
	return
}

func (actor *Actor) init() error {
	return nil
}

func (actor *Actor) Start() error {
	for _, v := range actor.comp_map {
		v.Start()
	}
	return nil
}

func (actor *Actor) Close() error {
	for _, v := range actor.comp_map {
		v.Close()
	}
	return nil
}

func (actor *Actor) Tick() error {
	for _, v := range actor.comp_map {
		v.Tick()
	}
	return nil
}

func (actor *Actor) ID() string {
	return actor.id
}

func (actor *Actor) Comp(id string) (comp.Component, bool) {
	c, ok := actor.comp_map[id]
	return c, ok
}

func (actor *Actor) AddComp(c comp.Component) {
	if c == nil {
		return
	}
	actor.comp_map[c.ID()] = c
}

func (actor *Actor) DelComp(c comp.Component) {
	if c == nil {
		return
	}
	delete(actor.comp_map, c.ID())
}

func (actor *Actor) DelCompByID(id string) {
	delete(actor.comp_map, id)
}
