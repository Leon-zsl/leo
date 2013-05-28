/* this is queue data structure
*/

package base

type Queue struct {
	nodes []interface{}
	head int
	tail int
	count int
}

func NewQueue(size int) *Queue {
	if size <= 0 {
		size = 8
	}

	q := new(Queue)
	q.head = 0
	q.tail = 0
	q.count = 0
	q.nodes = make([]interface{}, size)
	return q
}

func (q *Queue) Push(v interface{}) {
	if q.count > 0 && q.head == q.tail {
		ns := make([]interface{}, len(q.nodes) * 2)
		copy(ns, q.nodes[q.head:])
		copy(ns[len(q.nodes) - q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = q.count
		q.nodes = ns
	}

	q.nodes[q.tail] = v
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

func (q *Queue) Pop() interface{} {
	if q.count == 0 {
		return nil
	}

	n := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return n
}

func (q *Queue) Peek() interface{} {
	if q.count == 0 {
		return nil
	}	
	return q.nodes[q.head]
}

func (q *Queue) Count() int {
	return q.count
}

func (q *Queue) Empty() bool {
	return q.count == 0
}

