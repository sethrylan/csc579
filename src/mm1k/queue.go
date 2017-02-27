package mm1k

import (
	"log"
	"math"
)

// FIFO queue implementation based on https://gist.github.com/moraes/2141121

// receives a pointer so it can modify
func (q *FIFOQueue) push(c Customer) {
	q.len += 1
	q.a = append(q.a, c)
}

// receives a pointer so it can modify
func (q *FIFOQueue) pop() (c Customer) {
	c = (q.a)[0]
	q.a = (q.a)[1:]
	q.len -= 1
	return
}

func (q *FIFOQueue) peek() (n Customer) {
	n = (q.a)[0]
	return
}

func (q *FIFOQueue) Len() int {
	return q.len
}

func (q *FIFOQueue) Full() bool {
	return q.len == q.size
}

type FIFOQueue struct {
	a    []Customer
	size int
	len  int
}

func NewFIFOQueue(c int) (fifo *FIFOQueue) {
	return &FIFOQueue{a: make([]Customer, 0), size: c, len: 0}
}

func (q *FIFOQueue) Dequeue() (c Customer) {
	return q.pop()
}

// Given the time to next arrival and time to next service opening, returns next available time and queue position
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

func (q *FIFOQueue) NextCompletion() (next float64) {
	if q.Len() > 0 {
		next = q.peek().Start + q.peek().Service
	} else {
		next = math.Inf(+1)
	}
	return
}
