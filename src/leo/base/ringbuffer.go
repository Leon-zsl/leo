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
	rb.lock.Unlock()
	v := rb.queue.Pop()

	if v != nil {
		return v.(*Packet)
	} else {
		return nil
	}
}

func (rb *RingBuffer) Count() int {
	rb.lock.Lock()
	defer rb.lock.Unlock()
	return rb.queue.Count()
}

func (rb *RingBuffer) Empty() bool {
	rb.lock.Lock()
	defer rb.lock.Unlock()
	return rb.queue.Empty()
}

func (rb *RingBuffer) Peek() *Packet {
	rb.lock.Lock()
	rb.lock.Unlock()
	v := rb.queue.Peek()

	if v != nil {
		return v.(*Packet)
	} else {
		return nil
	}
}