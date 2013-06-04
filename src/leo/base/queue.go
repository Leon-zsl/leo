/* this is a simple queue
*/

package base

import (
	"sync"
)

type Queue struct {
	lock sync.Mutex
	rb *RingBuffer
}

func NewQueue() (q *Queue) {
	q = new(Queue)
	q.rb = NewRingBuffer(0)
	return
}

func (q *Queue) Push(v interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.rb.Push(v)
}

func (q *Queue) Pop() interface{} {
	q.lock.Lock()
	q.lock.Unlock()
	v := q.rb.Pop()

	if v != nil {
		return v
	} else {
		return nil
	}
}

func (q *Queue) Count() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.rb.Count()
}

func (q *Queue) Empty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.rb.Empty()
}

func (q *Queue) Peek() interface{} {
	q.lock.Lock()
	q.lock.Unlock()
	v := q.rb.Peek()

	if v != nil {
		return v
	} else {
		return nil
	}
}