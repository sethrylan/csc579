package mm1k

import (
	"log"
	"math"
	"sync"
)

// FIFO queue implementation based on https://gist.github.com/moraes/2141121

// FIFOQueue implements a queue with first-in-first-out behaviour
type FIFOQueue struct {
	a        []Customer
	capacity int
	lock     sync.Mutex
}

// receives a pointer so it can modify
func (q *FIFOQueue) push(c Customer) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.a = append(q.a, c)
}

// receives a pointer so it can modify
func (q *FIFOQueue) pop() (c Customer) {
	q.lock.Lock()
	defer q.lock.Unlock()
	c = q.peek()
	q.a = (q.a)[1:]
	return
}

func (q *FIFOQueue) peek() (n Customer) {
	n = (q.a)[0]
	return
}

// Len implements mm1k.Queue.Len
func (q *FIFOQueue) Len() int {
	return len(q.a)
}

// Full implements mm1k.Queue.Full
func (q *FIFOQueue) Full() bool {
	return q.Len() == q.capacity
}

// NewFIFOQueue returns a reference to a new FIFOQueue
func NewFIFOQueue(c int) (fifo *FIFOQueue) {
	return &FIFOQueue{a: make([]Customer, 0), capacity: c}
}

// Dequeue implements mm1k.Queue.Dequeue
func (q *FIFOQueue) Dequeue() (c Customer) {
	return q.pop()
}

// Enqueue implements mm1k.Queue.Enqueue
func (q *FIFOQueue) Enqueue(customer Customer) (cus Customer) {
	if q.Full() {
		log.Panicln("queue is full")
	}
	customer.Position = q.Len()
	if q.Len() > 0 {
		customer.Start = math.Max(q.a[q.Len()-1].Start+q.a[q.Len()-1].Service, customer.Arrival)
	} else {
		customer.Start = customer.Arrival
	}
	q.push(customer)
	return customer
}

// NextCompletion implements mm1k.Queue.NextCompletion
func (q *FIFOQueue) NextCompletion() (next float64) {
	if q.Len() > 0 {
		next = q.peek().Start + q.peek().Service
	} else {
		next = math.Inf(+1)
	}
	return
}
