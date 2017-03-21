package mm1k

import (
	"log"
	"math"
	"math/rand"
)

const p int = 4

// Priority implements a queue with first-in-first-out behaviour with multiple queues
type Priority struct {
	a          [p]FIFO
	preemptive bool
	next       int
	generator  *rand.Rand
}

// receives a pointer so it can modify
func (q *Priority) push(c Customer) {
	q.a[q.next].push(c)
}

// receives a pointer so it can modify
func (q *Priority) pop() (c Customer) {
	for i := 0; i < p; i++ {
		if q.a[i].Len() > 0 {
			c = q.a[i].pop()
		}
	}
	return
}

func (q *Priority) peek() (n Customer) {
	for i := 0; i < p; i++ {
		if q.a[i].Len() > 0 {
			n = q.a[i].peek()
		}
	}
	return
}

// Len implements mm1k.Queue.Len
func (q *Priority) Len() int {
	sum := 0
	for i := 0; i < p; i++ {
		sum += q.a[i].Len()
	}
	return sum
}

// Full implements mm1k.Queue.Full
func (q *Priority) Full() bool {
	return q.a[q.next].Full()
}

// NewPriority returns a reference to a new Priority queue
func NewPriority(c int, preemptive bool) (priority *Priority) {
	priority = &Priority{a: [p]FIFO{*NewFIFO(c), *NewFIFO(c), *NewFIFO(c), *NewFIFO(c)}, generator: rand.New(rand.NewSource(2))}
	// for i := 0; i < p; i++ {
	// 	priority.a[i] = *NewFIFO(c)
	// }
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
	customer.Start = customer.Arrival
	q.push(customer)
	q.next = q.generator.Perm(p)[0] // values comes from default source

	// sort.Sort(byService(q.a))

	// TODO: stop preemption
	// recalculate start times starting from beginning of queue
	// for i := 1; i < q.Len()-1; i++ {
	// 	q.a[i].Start = math.Max(q.a[i+1].Start+q.a[i+1].Service, q.a[i].Arrival)
	// 	log.Printf("non-preempting customer %d Start to %.03f\n", q.a[i].ID, q.a[i].Start)
	// }
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
