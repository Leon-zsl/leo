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
	defer rb.lock.Unlock()
	return rb.queue.Pop().(*Packet)
}

func (rb *RingBuffer) Count() int {
	return rb.queue.Count()
}

func (rb *RingBuffer) Empty() bool {
	return rb.queue.Empty()
}

func (rb *RingBuffer) Peek() *Packet {
	rb.lock.Lock()
	defer rb.lock.Unlock()
	return rb.queue.Peek().(*Packet)
}