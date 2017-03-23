package mm1k_test

import (
	"fmt"
	"mm1k"
)

func ExampleDifferentQueuesNP() {
	var cus mm1k.Customer
	var q mm1k.Queue
	q = mm1k.NewPriority(10, 1, false)
	for i := 0; i < 3; i++ {
		cus = q.Enqueue(mm1k.Customer{ID: 1, Arrival: .5, Service: .1})
		fmt.Printf("%v\n", cus)
	}

	fmt.Println()

	q = mm1k.NewPriority(10, 1, true)
	for i := 0; i < 3; i++ {
		cus = q.Enqueue(mm1k.Customer{ID: 1, Arrival: .5, Service: .1})
		fmt.Printf("%v\n", cus)
	}

	// Output:
	// {1 0.5 0.1 0.5 0 0 0 0}
	// {1 0.5 0.1 0.6 0 1 0 0}
	// {1 0.5 0.1 0.7 0 2 0 0}
	//
	// {1 0.5 0.1 0.5 0 0 0 0}
	// {1 0.5 0.1 0.6 0 1 0 0}
	// {1 0.5 0.1 0.7 0 2 0 0}

}

func ExampleDifferentQueuesNP2() {
	var cus mm1k.Customer
	var q mm1k.Queue
	q = mm1k.NewPriority(10, 2, false)
	for i := 0; i < 3; i++ {
		cus = q.Enqueue(mm1k.Customer{ID: 1, Arrival: .5, Service: .1})
		fmt.Printf("%v\n", cus)
	}

	fmt.Println()

	q = mm1k.NewPriority(10, 2, true)
	for i := 0; i < 3; i++ {
		cus = q.Enqueue(mm1k.Customer{ID: 1, Arrival: .5, Service: .1})
		fmt.Printf("%v\n", cus)
	}
	// Output:
	// {1 0.5 0.1 0.5 0 0 0 0}
	// {1 0.5 0.1 0.6 0 0 0 1}
	// {1 0.5 0.1 0.7 0 1 0 1}
	//
	// {1 0.5 0.1 0.5 0 0 0 0}
	// {1 0.5 0.1 0.6 0 0 0 1}
	// {1 0.5 0.1 0.7 0 1 0 1}
}
