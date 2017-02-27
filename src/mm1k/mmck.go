package mm1k

import (
	"log"
	"math"
)

type Customer struct {
	ID         int
	Arrival    float64 // time of arrival of the customer
	Service    float64
	Start      float64 // start time of service of the customer
	Completion float64 // time of completion/departure of the customer
	Position   int     // identifier of the seat that sits the waiting customer
}

type Queue interface {
	Enqueue(customer Customer) Customer
	Dequeue() Customer
	Len() int
	NextCompletion() float64
	Full() bool
}

func Run(arrivalDistribution Distribution, q Queue, serviceDistribution Distribution) (rejects, completes <-chan Customer) {
	rejected := make(chan Customer) // Unbuffered channels ensure deterministic simulation
	completed := make(chan Customer)
	var clock float64 // master clock

	go func() {
		var t1 float64 = arrivalDistribution.Get() // time of next arrival
		var t2 float64 = math.Inf(1)               // time of next completion
		var id int                                 // Incremented Customer ID
		for {                                      // Do forever
			if t1 < t2 { // If next arrival is before next completion -> Event: Arrival
				clock = t1 // Set clock to time of next arrival.
				if q.Full() {
					rejected <- Customer{ID: id, Arrival: t1, Service: serviceDistribution.Get()}
				} else {
					q.Enqueue(Customer{ID: id, Arrival: t1, Service: serviceDistribution.Get()})
				}
				id += 1
				t1 = clock + arrivalDistribution.Get() // Set t1 to time of next arrival.
				if q.Len() == 1 {                      // If queue is not empty
					t2 = q.NextCompletion() // then set t2 to time of next completion.
				}
			} else { // If next arrival is after next completion -> Event: Departure
				if !math.IsInf(t2, 1) { // if next completion exists
					clock = t2 // Set time to time of next completion.
					log.Printf("Len() == %d\n", q.Len())
					log.Printf("t1() == %f\n", t1)
					log.Printf("t2() == %f\n", t2)
					customer := q.Dequeue()                // Remove customer from queue
					customer.Completion = t2               // Set completion time
					customer.Start = t2 - customer.Service // and start time
					completed <- customer                  // and add the customer to the completed channel
				}
				t2 = q.NextCompletion()
			}
		}
	}()

	rejects = rejected
	completes = completed
	return
}
