package mm1k

import (
	"log"
	"math"
	"sort"
	"sync"
)

// SJF/SJN queue implementation

// SJF implements a queue with shortest job first semantics behaviour
type SJF struct {
	a        []Customer
	capacity int
	lock     sync.Mutex
}

// receives a pointer so it can modify
func (q *SJF) push(c Customer) {
	q.a = append(q.a, c)
}

// receives a pointer so it can modify
func (q *SJF) pop() (c Customer) {
	c, q.a = q.peek(), (q.a)[1:]
	return
}

func (q *SJF) peek() (n Customer) {
	n = (q.a)[0]
	return
}

// Len implements mm1k.Queue.Len
func (q *SJF) Len() int {
	return len(q.a)
}

// Full implements mm1k.Queue.Full
func (q *SJF) Full() bool {
	return q.Len() == q.capacity
}

// NewSJF returns a reference to a new SJF
func NewSJF(c int) (sjf *SJF) {
	return &SJF{a: make([]Customer, 0), capacity: c}
}

// Dequeue implements mm1k.Queue.Dequeue
func (q *SJF) Dequeue() (c Customer) {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.pop()
}

// Enqueue implements mm1k.Queue.Enqueue
func (q *SJF) Enqueue(customer Customer) (cus Customer) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.Full() {
		log.Panicln("queue is full")
	}

	// TODO: stop preemption
	// compose queue so that
	//    [a...] + cus + [b...],
	// where a is the list customers with a startTime > cus.arrival and b is the list of remaining customers
	customer.Start = customer.Arrival
	q.push(customer)
	sort.Sort(byService(q.a))

	// recalculate start times starting from beginning of queue
	for i := 1; i < q.Len()-1; i++ {
		q.a[i].Start = math.Max(q.a[i+1].Start+q.a[i+1].Service, q.a[i].Arrival)
		log.Printf("non-preempting customer %d Start to %.03f\n", q.a[i].ID, q.a[i].Start)
	}
	return customer
}

// NextCompletion implements mm1k.Queue.NextCompletion
func (q *SJF) NextCompletion() (next float64) {
	if q.Len() > 0 {
		next = q.peek().Start + q.peek().Service
	} else {
		next = math.Inf(+1)
	}
	return
}
