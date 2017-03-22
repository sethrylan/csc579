package mm1k_test

import (
	"fmt"
	"mm1k"
)

func ExampleNewSFJPreemptive() {
	// var t1 float64
	var cus mm1k.Customer
	q := mm1k.NewSJF(1, true)

	// Empty Queue
	fmt.Printf("Len() = %d\n", q.Len())
	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())
	fmt.Printf("Full() = %v\n", q.Full())

	// Add 1 item
	cus = q.Enqueue(mm1k.Customer{ID: 1, Arrival: .5, Service: .1})
	fmt.Printf("pos = %d\n", cus.Position)
	fmt.Printf("Len() = %d\n", q.Len())
	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())
	fmt.Printf("Full() = %v\n", q.Full())

	// Remove 1 item
	q.Dequeue()
	// fmt.Printf("customer = %v", customer)
	fmt.Printf("Len() = %d\n", q.Len())
	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())
	fmt.Printf("Full() = %v\n", q.Full())

	// Output:
	// Len() = 0
	// NextCompletion() = +Inf
	// Full() = false
	// pos = 0
	// Len() = 1
	// NextCompletion() = 0.60
	// Full() = true
	// Len() = 0
	// NextCompletion() = +Inf
	// Full() = false
}

func ExampleNewSFJNonPreemptive() {
	// var t1 float64
	var cus mm1k.Customer
	q := mm1k.NewSJF(1, false)

	// Empty Queue
	fmt.Printf("Len() = %d\n", q.Len())
	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())
	fmt.Printf("Full() = %v\n", q.Full())

	// Add 1 item
	cus = q.Enqueue(mm1k.Customer{ID: 1, Arrival: .5, Service: .1})
	fmt.Printf("pos = %d\n", cus.Position)
	fmt.Printf("Len() = %d\n", q.Len())
	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())
	fmt.Printf("Full() = %v\n", q.Full())

	// Remove 1 item
	q.Dequeue()
	// fmt.Printf("customer = %v", customer)
	fmt.Printf("Len() = %d\n", q.Len())
	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())
	fmt.Printf("Full() = %v\n", q.Full())

	// Output:
	// Len() = 0
	// NextCompletion() = +Inf
	// Full() = false
	// pos = 0
	// Len() = 1
	// NextCompletion() = 0.60
	// Full() = true
	// Len() = 0
	// NextCompletion() = +Inf
	// Full() = false
}
