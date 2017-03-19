package mm1k

import (
	"log"
	"math"
	"sync"
)

// LIFO stack implementation

// LIFO implements a stack with last-in-first-out behaviour
type LIFO struct {
	a        []Customer
	capacity int
	lock     sync.Mutex
}

// receives a pointer so it can modify
func (q *LIFO) push(c Customer) {
	q.a = append(q.a, c)
}

// receives a pointer so it can modify
func (q *LIFO) pop() (c Customer) {
	c, q.a = q.peek(), (q.a)[:q.Len()-1]
	return
}

func (q *LIFO) peek() (n Customer) {
	n = (q.a)[q.Len()-1]
	return
}

// Len implements mm1k.Queue.Len
func (q *LIFO) Len() int {
	return len(q.a)
}

// Full implements mm1k.Queue.Full
func (q *LIFO) Full() bool {
	return q.Len() == q.capacity
}

// NewLIFO returns a reference to a new FIFO
func NewLIFO(c int) (lifo *LIFO) {
	return &LIFO{a: make([]Customer, 0), capacity: c}
}

// Dequeue implements mm1k.Queue.Dequeue
func (q *LIFO) Dequeue() (c Customer) {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.pop()
}

// Enqueue implements mm1k.Queue.Enqueue
func (q *LIFO) Enqueue(customer Customer) (cus Customer) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.Full() {
		log.Panicln("queue is full")
	}
	customer.Position = q.Len()

	// If the new customer's arrival time is before the Start time of an existing
	// customer, then we insert the new customer after (at an earlier index) than
	// the existing customer. If we insert the new customer after the existing
	// customer, then we have pre-empted the existing customer's place in the
	// queue.
	for i := q.Len() - 1; i >= 0 && customer.Arrival <= q.a[i].Start+q.a[i].Service; i-- {
		// on last pass, i is equal to the intended position of the new customer
		customer.Position = i
	}

	if customer.Position == q.Len() {
		customer.Start = customer.Arrival
	} else {
		customer.Start = q.a[customer.Position].Start + q.a[customer.Position].Service
	}
	// append everything up to and including the nonPreemptive position + customer + rest of stack
	q.a = append(q.a[:customer.Position], append([]Customer{customer}, q.a[customer.Position:]...)...)

	// recalculate start times starting from end of stack
	for i := q.Len() - 2; i >= 0; i-- {
		q.a[i].Start = math.Max(q.a[i+1].Start+q.a[i+1].Service, q.a[i].Arrival)
		log.Printf("non-preempting customer %d Start to %.03f\n", q.a[i].ID, q.a[i].Start)
	}
	return customer
}

// NextCompletion implements mm1k.Queue.NextCompletion
func (q *LIFO) NextCompletion() (next float64) {
	if q.Len() > 0 {
		next = q.peek().Start + q.peek().Service
	} else {
		next = math.Inf(+1)
	}
	return
}
