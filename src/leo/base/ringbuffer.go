/* this is a simple ring buffer
*/

package base

import (
	"sync"
)

type RingBuffer struct {
	lock sync.Mutex
	queue *Queue
}

func NewRingBuffer() (rb *RingBuffer) {
	rb = new(RingBuffer)
	rb.queue = NewQueue(0)
	return
}

func (rb *RingBuffer) Push(pk *Packet) {
	rb.lock.Lock()
	defer rb.lock.Unlock()
	rb.queue.Push(pk)
}

func (rb *RingBuffer) Pop() *Packet {
	rb.lock.Lock()
	v := rb.queue.Pop()
	rb.lock.Unlock()

	if v != nil {
		return v.(*Packet)
	} else {
		return nil
	}
}

func (rb *RingBuffer) Count() int {
	return rb.queue.Count()
}

func (rb *RingBuffer) Empty() bool {
	return rb.queue.Empty()
}

func (rb *RingBuffer) Peek() *Packet {
	rb.lock.Lock()
	v := rb.queue.Peek()
	rb.lock.Unlock()

	if v != nil {
		return v.(*Packet)
	} else {
		return nil
	}
}