/* this is ring buffer data structure
*/

package base

type RingBuffer struct {
	nodes []interface{}
	head int
	tail int
	count int
}

func NewRingBuffer(size int) *RingBuffer {
	if size <= 0 {
		size = 8
	}

	rb := new(RingBuffer)
	rb.head = 0
	rb.tail = 0
	rb.count = 0
	rb.nodes = make([]interface{}, size)
	return rb
}

func (rb *RingBuffer) Push(v interface{}) {
	if rb.count > 0 && rb.head == rb.tail {
		ns := make([]interface{}, len(rb.nodes) * 2)
		copy(ns, rb.nodes[rb.head:])
		copy(ns[len(rb.nodes) - rb.head:], rb.nodes[:rb.head])
		rb.head = 0
		rb.tail = rb.count
		rb.nodes = ns
	}

	rb.nodes[rb.tail] = v
	rb.tail = (rb.tail + 1) % len(rb.nodes)
	rb.count++
}

func (rb *RingBuffer) Pop() interface{} {
	if rb.count == 0 {
		return nil
	}

	n := rb.nodes[rb.head]
	rb.head = (rb.head + 1) % len(rb.nodes)
	rb.count--
	return n
}

func (rb *RingBuffer) Peek() interface{} {
	if rb.count == 0 {
		return nil
	}	
	return rb.nodes[rb.head]
}

func (rb *RingBuffer) Count() int {
	return rb.count
}

func (rb *RingBuffer) Empty() bool {
	return rb.count == 0
}

