package mm1k

import (
	"log"
	"math"
)

// FIFO queue implementation based on https://gist.github.com/moraes/2141121

// receives a pointer so it can modify
func (q *FIFOQueue) push(c Customer) {
	q.len++
	q.a = append(q.a, c)
}

// receives a pointer so it can modify
func (q *FIFOQueue) pop() (c Customer) {
	c = (q.a)[0]
	q.a = (q.a)[1:]
	q.len--
	return
}

func (q *FIFOQueue) peek() (n Customer) {
	n = (q.a)[0]
	return
}

// Len returns the length of the queue
func (q *FIFOQueue) Len() int {
	return q.len
}

// Full returns true is the queue is full
func (q *FIFOQueue) Full() bool {
	return q.len == q.size
}

// FIFOQueue implements a queue with first-in-first-out behaviour
type FIFOQueue struct {
	a    []Customer
	size int
	len  int
}

// NewFIFOQueue returns a reference to a new FIFOQueue
func NewFIFOQueue(c int) (fifo *FIFOQueue) {
	return &FIFOQueue{a: make([]Customer, 0), size: c, len: 0}
}

// Dequeue implements mm1k.Queue.Dequeue
func (q *FIFOQueue) Dequeue() (c Customer) {
	return q.pop()
}

// Enqueue implements mm1k.Queue.Enqueue
func (q *FIFOQueue) Enqueue(customer Customer) (cus Customer) {
	if q.len == q.size {
		log.Panicln("queue is full")
	}
	customer.Position = q.Len()
	if q.len > 0 {
		customer.Start = math.Max(q.a[q.len-1].Start+q.a[q.len-1].Service, customer.Arrival)
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
