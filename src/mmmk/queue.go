package mmmk

//
// import (
// 	"math"
// 	"mm1k"
// 	"sync"
// )
//
// // FIFO queue implementation based on https://gist.github.com/moraes/2141121
//
// // Queue implements a queue with first-in-first-out behaviour
// type Queue struct {
// 	a    []mm1k.Customer
// 	lock sync.Mutex
// }
//
// // receives a pointer so it can modify
// func (q *Queue) push(c mm1k.Customer) {
// 	q.a = append(q.a, c)
// }
//
// // receives a pointer so it can modify
// func (q *Queue) pop() (c mm1k.Customer) {
// 	c, q.a = q.peek(), (q.a)[1:]
// 	return
// }
//
// func (q *Queue) peek() (n mm1k.Customer) {
// 	n = (q.a)[0]
// 	return
// }
//
// func (q *Queue) last() (n mm1k.Customer) {
// 	if len(q.a) > 0 {
// 		n = (q.a)[len(q.a)-1]
// 	} else {
// 		n = mm1k.Customer{}
// 	}
// 	return
// }
//
// // Len implements mm1k.Queue.Len
// func (q *Queue) Len() int {
// 	return len(q.a)
// }
//
// // NewQueue returns a reference to a new queue
// func NewQueue() (queue *Queue) {
// 	return &Queue{a: make([]mm1k.Customer, 1000000)}
// }
//
// // Dequeue implements mm1k.Queue.Dequeue
// func (q *Queue) Dequeue() (c mm1k.Customer) {
// 	q.lock.Lock()
// 	defer q.lock.Unlock()
// 	return q.pop()
// }
//
// // Enqueue implements mm1k.Queue.Enqueue
// func (q *Queue) Enqueue(customer mm1k.Customer) (cus mm1k.Customer) {
// 	return q.enqueue(customer, q.last())
// }
//
// func (q *Queue) enqueue(customer mm1k.Customer, last mm1k.Customer) (cus mm1k.Customer) {
// 	q.lock.Lock()
// 	defer q.lock.Unlock()
// 	customer.Position = q.Len()
// 	q.push(customer)
// 	return customer
// }
//
// // NextCompletion implements mm1k.Queue.NextCompletion
// func (q *Queue) NextCompletion() (next float64) {
// 	if q.Len() > 0 {
// 		next = q.peek().Start + q.peek().Service
// 	} else {
// 		next = math.Inf(+1)
// 	}
// 	return
// }
