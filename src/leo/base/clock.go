/* this is clock service
*/

package base

import (
	"time"
)

//use this in one goroutine
type Clock struct {
	running bool

	pretime time.Time
	curtime time.Time
	deltatime time.Duration
}

func NewClock() (clock *Clock, err error) {
	clock = new(Clock)
	err = nil
	return
}

func (clock *Clock) Start() error {
	clock.pretime = time.Now()
	clock.curtime = time.Now()
	clock.deltatime = 0
	clock.running = true
	return nil
}

func (clock *Clock) Close() error {
	clock.running = false
	return nil
}

func (clock *Clock) CurTime() time.Time {
	return clock.curtime
}

func (clock *Clock) DeltaTime() time.Duration {
	return clock.deltatime
}

func (clock *Clock) Tick() error {
	if clock.running {
		clock.curtime = time.Now()
		clock.deltatime = clock.curtime.Sub(clock.pretime)
		clock.pretime = clock.curtime
	}
	return nil
}