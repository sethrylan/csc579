package mm1k_test

import (
	"fmt"
	"mm1k"
)

func ExampleDifferentQueues() {
	var cus mm1k.Customer
	q := mm1k.NewPriority(10, true)
	for i := 0; i < 10; i++ {
		cus = q.Enqueue(mm1k.Customer{ID: 1, Arrival: .5, Service: .1})
		fmt.Printf("queue = %d\n", cus.PriorityQueue)
	}
	// Output:
	//queue = 0
	// queue = 3
	// queue = 3
	// queue = 1
	// queue = 1
	// queue = 1
	// queue = 1
	// queue = 0
	// queue = 2
	// queue = 3
}

// func ExampleNewPriority() {
// 	// var t1 float64
// 	var cus mm1k.Customer
// 	q := mm1k.NewPriority(1, true)
//
// 	// Empty Queue
// 	fmt.Printf("Len() = %d\n", q.Len())
// 	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())
// 	fmt.Printf("Full() = %v\n", q.Full())
//
// 	// Add 1 item
// 	cus = q.Enqueue(mm1k.Customer{ID: 1, Arrival: .5, Service: .1})
// 	fmt.Printf("pos = %d\n", cus.Position)
// 	fmt.Printf("queue = %d\n", cus.PriorityQueue)
// 	fmt.Printf("Len() = %d\n", q.Len())
// 	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())
// 	fmt.Printf("Full() = %v\n", q.Full())
//
// 	// Remove 1 item
// 	q.Dequeue()
// 	// fmt.Printf("customer = %v", customer)
// 	fmt.Printf("Len() = %d\n", q.Len())
// 	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())
// 	fmt.Printf("Full() = %v\n", q.Full())
//
// 	// Output:
// 	// Len() = 0
// 	// NextCompletion() = +Inf
// 	// Full() = false
// 	// pos = 0
// 	// queue = 0
// 	// Len() = 1
// 	// NextCompletion() = 0.60
// 	// Full() = true
// 	// Len() = 0
// 	// NextCompletion() = +Inf
// 	// Full() = false
// }
