package mm1k_test

import (
	"fmt"
	"mm1k"
)

func ExampleNewSFJNP() {
	var cus mm1k.Customer
	q := mm1k.NewSJF(2, false)

	// Empty Queue
	fmt.Printf("Len() = %d\n", q.Len())
	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())
	fmt.Printf("Full() = %v\n", q.Full())

	// Add 1 customer
	cus = q.Enqueue(mm1k.Customer{ID: 1, Arrival: .5, Service: .2})
	fmt.Printf("%v\n", cus)
	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())

	// Add 1 more customer
	cus = q.Enqueue(mm1k.Customer{ID: 2, Arrival: .5, Service: .1})
	fmt.Printf("%v\n", cus)
	fmt.Printf("NextCompletion() = %.2f\n", q.NextCompletion())

	// Output:
	// Len() = 0
	// NextCompletion() = +Inf
	// Full() = false
	// {1 0.5 0.2 0.5 0 0 0 0}
	// NextCompletion() = 0.70
	// {2 0.5 0.1 0 0 1 0 0}
	// NextCompletion() = 0.80
}

// func ExampleNewSFJNonPreemptive2() {
// 	// var t1 float64
// 	var cus mm1k.Customer
// 	q := mm1k.NewSJF(5, false)
//
// 	for i := 0; i < 5; i++ {
// 		cus = q.Enqueue(mm1k.Customer{ID: i, Arrival: .5 * float64(i), Service: .1 * float64(i)})
// 		fmt.Printf("%v\n", cus)
// 	}
//
// 	// Output:
// 	// Len() = 0
// 	// NextCompletion() = +Inf
// 	// Full() = false
// 	// pos = 0
// 	// Len() = 1
// 	// NextCompletion() = 0.60
// 	// Full() = true
// 	// Len() = 0
// 	// NextCompletion() = +Inf
// 	// Full() = false
// }
