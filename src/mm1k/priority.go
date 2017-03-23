package mm1k

import (
	"log"
	"math"
	"math/rand"
)

// Priority implements a queue with first-in-first-out behaviour with multiple queues
type Priority struct {
	a          []FIFO
	p          int
	preemptive bool
	next       int
	generator  *rand.Rand
}

// receives a pointer so it can modify
func (q *Priority) pop() (c Customer) {
	for i := 0; i < q.p; i++ {
		if q.a[i].Len() > 0 {
			c = q.a[i].pop()
			return
		}
	}
	return
}

func (q *Priority) Peek() (c Customer) {
	return q.peek()
}

func (q *Priority) peek() (c Customer) {
	for i := 0; i < q.p; i++ {
		if q.a[i].Len() > 0 {
			return q.a[i].peek()
		}
	}
	return
}

func (q *Priority) last(maxQueue int) (last Customer) {
	for i := 0; i <= maxQueue; i++ {
		if q.a[i].Len() > 0 {
			last = q.a[i].last()
		}
	}
	return
}

// Len implements mm1k.Queue.Len
func (q *Priority) Len() int {
	sum := 0
	for i := 0; i < q.p; i++ {
		sum += q.a[i].Len()
	}
	return sum
}

// Full implements mm1k.Queue.Full
func (q *Priority) Full() bool {
	return q.a[q.next].Full()
}

// NewPriority returns a reference to a new Priority queue
func NewPriority(c int, p int, preemptive bool) (priority *Priority) {
	priority = &Priority{a: make([]FIFO, p), p: p, preemptive: preemptive, generator: rand.New(rand.NewSource(2))}
	for i := 0; i < p; i++ {
		priority.a[i] = *NewFIFO(c)
	}
	return
}

// Dequeue implements mm1k.Queue.Dequeue
func (q *Priority) Dequeue() (c Customer) {
	return q.pop()
}

// Enqueue implements mm1k.Queue.Enqueue
func (q *Priority) Enqueue(customer Customer) (cus Customer) {
	if q.Full() {
		log.Panicln("queue is full")
	}

	customer.PriorityQueue = q.next
	// customer.Start = math.Max(customer.Arrival, q.last(q.next).Start+q.last(q.next).Service)
	// fmt.Printf("customer.Start = %f\n", customer.Start)
	customer = q.a[q.next].enqueue(customer, q.last(q.next))
	// fmt.Printf("customer.Start = %f\n", customer.Start)

	if q.preemptive {

	} else {
		var lastCompletion float64
		// recalculate start times of lower priority queues
		for i := 0; i < q.p-1; i++ { // For each lower priority queue,
			if q.a[i].Len() > 0 {
				lastCompletion = q.a[i].last().Start + q.a[i].last().Service // find the last completion for queue i
				// fmt.Printf("lastCompletion = %f\n", lastCompletion)
			}
			if q.a[i+1].Len() > 0 { // if the queue has items
				firstStart := q.a[i+1].peek().Start
				if firstStart < lastCompletion { // and first start in queue is before lastCompletion in higher priority queue
					for j := 0; j < q.a[i+1].Len(); j++ { // then for all al
						q.a[i+1].a[j].Start += lastCompletion - firstStart
					}
				}
			}
		}
	}

	q.next = q.generator.Perm(q.p)[0] // values comes from default source

	return customer
}

// NextCompletion implements mm1k.Queue.NextCompletion
func (q *Priority) NextCompletion() (next float64) {
	if q.Len() > 0 {
		n := q.peek()
		next = n.Start + n.Service
	} else {
		next = math.Inf(+1)
	}
	return
}
