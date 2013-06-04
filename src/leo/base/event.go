/* this is the event interface
*/

package base

type Event interface {
	Op() int
	Args() []interface{}
	//stop to continue
	Handled() bool
}

type EventHandler interface {
	HandleEvent(e Event)
}