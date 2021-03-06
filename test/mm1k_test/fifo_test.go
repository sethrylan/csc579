package mm1k_test

import (
	"fmt"
	"mm1k"
)

func ExampleNewFIFO() {
	// var t1 float64
	var cus mm1k.Customer
	q := mm1k.NewFIFO(1)

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

func ExampleNewFIFOAdd2() {
	q := mm1k.NewFIFO(2)

	// Add 2 items
	q.Enqueue(mm1k.Customer{ID: 1, Arrival: .5, Service: .1})
	q.Enqueue(mm1k.Customer{ID: 2, Arrival: 1, Service: .1})

	fmt.Printf("Len() = %d\n", q.Len())
	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())
	fmt.Printf("Full() = %v\n", q.Full())

	// Remove 1 item
	customer := q.Dequeue()
	fmt.Printf("customer.(ID, Position, Start) = (%d, %d, %.2f)\n", customer.ID, customer.Position, customer.Start)
	fmt.Printf("Len() = %d\n", q.Len())
	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())
	fmt.Printf("Full() = %v\n", q.Full())

	// Output:
	// Len() = 2
	// NextCompletion() = 0.60
	// Full() = true
	// customer.(ID, Position, Start) = (1, 0, 0.50)
	// Len() = 1
	// NextCompletion() = 1.10
	// Full() = false
}

func ExampleNewFIFOAdd2Overlap() {
	var cus mm1k.Customer
	q := mm1k.NewFIFO(2)

	// Add 2 items
	cus = q.Enqueue(mm1k.Customer{ID: 1, Arrival: .5, Service: .5})
	fmt.Printf("customer.(ID, Position, Start) = (%d, %d, %.2f)\n", cus.ID, cus.Position, cus.Start)
	cus = q.Enqueue(mm1k.Customer{ID: 2, Arrival: .6, Service: .5})
	fmt.Printf("customer.(ID, Position, Start) = (%d, %d, %.2f)\n", cus.ID, cus.Position, cus.Start)

	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())

	// // Remove 1 item
	customer := q.Dequeue()
	fmt.Printf("customer.(ID, Position, Start) = (%d, %d, %.2f)\n", customer.ID, customer.Position, customer.Start)
	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())

	// Output:
	// customer.(ID, Position, Start) = (1, 0, 0.50)
	// customer.(ID, Position, Start) = (2, 1, 1.00)
	// NextCompletion() = 1.00
	// customer.(ID, Position, Start) = (1, 0, 0.50)
	// NextCompletion() = 1.50
}
