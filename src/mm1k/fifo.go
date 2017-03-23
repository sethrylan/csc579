package mm1k

import (
	"log"
	"math"
	"sync"
)

// FIFO queue implementation based on https://gist.github.com/moraes/2141121

// FIFO implements a queue with first-in-first-out behaviour
type FIFO struct {
	a        []Customer
	capacity int
	lock     sync.Mutex
}

// receives a pointer so it can modify
func (q *FIFO) push(c Customer) {
	q.a = append(q.a, c)
}

// receives a pointer so it can modify
func (q *FIFO) pop() (c Customer) {
	c, q.a = q.peek(), (q.a)[1:]
	return
}

func (q *FIFO) peek() (n Customer) {
	n = (q.a)[0]
	return
}

func (q *FIFO) last() (n Customer) {
	if len(q.a) > 0 {
		n = (q.a)[len(q.a)-1]
	} else {
		n = Customer{}
	}
	return
}

// Len implements mm1k.Queue.Len
func (q *FIFO) Len() int {
	return len(q.a)
}

// Full implements mm1k.Queue.Full
func (q *FIFO) Full() bool {
	return q.Len() == q.capacity
}

// NewFIFO returns a reference to a new FIFO queue
func NewFIFO(c int) (fifo *FIFO) {
	return &FIFO{a: make([]Customer, 0), capacity: c}
}

// Dequeue implements mm1k.Queue.Dequeue
func (q *FIFO) Dequeue() (c Customer) {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.pop()
}

// Enqueue implements mm1k.Queue.Enqueue
func (q *FIFO) Enqueue(customer Customer) (cus Customer) {
	return q.enqueue(customer, q.last())
}

func (q *FIFO) enqueue(customer Customer, last Customer) (cus Customer) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.Full() {
		log.Panicln("queue is full")
	}
	customer.Position = q.Len()
	if last != (Customer{}) {
		customer.Start = math.Max(last.Start+last.Service, customer.Arrival)
	} else {
		customer.Start = customer.Arrival
	}
	q.push(customer)
	return customer
}

// NextCompletion implements mm1k.Queue.NextCompletion
func (q *FIFO) NextCompletion() (next float64) {
	if q.Len() > 0 {
		next = q.peek().Start + q.peek().Service
	} else {
		next = math.Inf(+1)
	}
	return
}
