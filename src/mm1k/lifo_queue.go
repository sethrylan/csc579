package mm1k

import (
	"log"
	"math"
	"sync"
)

// LIFO stack implementation

// LIFOQueue implements a stack with last-in-first-out behaviour
type LIFOQueue struct {
	a        []Customer
	capacity int
	lock     sync.Mutex
}

// receives a pointer so it can modify
func (q *LIFOQueue) push(c Customer) {
	q.a = append(q.a, c)
}

// receives a pointer so it can modify
func (q *LIFOQueue) pop() (c Customer) {
	c = q.peek()
	q.a = (q.a)[:q.Len()-1]
	return
}

func (q *LIFOQueue) peek() (n Customer) {
	n = (q.a)[q.Len()-1]
	return
}

// Len implements mm1k.Queue.Len
func (q *LIFOQueue) Len() int {
	return len(q.a)
}

// Full implements mm1k.Queue.Full
func (q *LIFOQueue) Full() bool {
	return q.Len() == q.capacity
}

// NewLIFOQueue returns a reference to a new FIFOQueue
func NewLIFOQueue(c int) (lifo *LIFOQueue) {
	return &LIFOQueue{a: make([]Customer, 0), capacity: c}
}

// Dequeue implements mm1k.Queue.Dequeue
func (q *LIFOQueue) Dequeue() (c Customer) {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.pop()
}

// Enqueue implements mm1k.Queue.Enqueue
func (q *LIFOQueue) Enqueue(customer Customer) (cus Customer) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.Full() {
		log.Panicln("queue is full")
	}
	customer.Position = q.Len()
	// if q.Len() > 0 {
	// 	customer.Start = math.Max(q.a[q.Len()-1].Start+q.a[q.Len()-1].Service, customer.Arrival)
	// } else {
	customer.Start = customer.Arrival
	// }
	q.push(customer)

	// TODO: recalculate start times starting from end of stack
	for i := q.Len() - 2; i >= 0; i-- {
		q.a[i].Start = math.Max(q.a[i+1].Start+q.a[i+1].Service, q.a[i].Arrival)
	}
	return customer
}

// NextCompletion implements mm1k.Queue.NextCompletion
func (q *LIFOQueue) NextCompletion() (next float64) {
	if q.Len() > 0 {
		next = q.peek().Start + q.peek().Service
	} else {
		next = math.Inf(+1)
	}
	return
}
