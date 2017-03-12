package mm1k

import (
	"fmt"
	"log"
	"math"
)

// A Customer contains tracking and history for a customer
type Customer struct {
	ID               int     // montonically increasing identifier
	Arrival          float64 // absolute time of arrival of the customer
	Service          float64 // relative time of service interal
	Start            float64 // absolution start time of service of the customer
	Departure        float64 // absolute time of completion/departure of the customer
	Position         int     // position of this customer in the queue at time of arrival
	QueueAtDeparture int     // number of customers in queue at time of departure
}

// A Queue type defines the common operations for a service queue
type Queue interface {
	// Enqueue adds customer to the queue
	Enqueue(customer Customer) Customer

	// Dequeue removes and returns the next customer determined by the queue service discipline
	Dequeue() Customer

	// Len return the length of the queue
	Len() int

	// NextCompletion returns the time of next service completion, or +∞ if no next service
	NextCompletion() float64

	// Full returns true if queue is full
	Full() bool
}

// Simulate will terminate once C customers have completed. Assume that
// at time t = 0 the system is empty. Draw a random number to decide when the
// first arrival will occur.
func Simulate(λ float64, µ float64, q Queue, C int, seed int64) (completes []Customer, rejects []Customer) {
	var customer Customer
	var rejected, completed <-chan Customer
	rejected, completed = Run(
		NewExpDistribution(λ, seed),
		q,
		NewExpDistribution(µ, seed+1),
	)
	for len(completes) < C {
		select {
		case customer = <-rejected:
			rejects = append(rejects, customer)
			logCustomer(customer)
		case customer = <-completed:
			completes = append(completes, customer)
			logCustomer(customer)
		}
	}
	return
}

// Run will continually add and service customers using an event loop
func Run(arrivalDistribution Distribution, q Queue, serviceDistribution Distribution) (rejects, completes <-chan Customer) {
	rejected := make(chan Customer)  // Unbuffered channels ensure deterministic simulation
	completed := make(chan Customer) //
	var clock float64                // master clock

	go func() {
		var t1 = arrivalDistribution.Get() // time of next arrival
		var t2 = math.Inf(1)               // time of next completion (∞ for no schedule Customer)
		var id int                         // Incremented Customer ID
		for {                              // Do forever
			if t1 < t2 { // If next arrival is before next completion -> Event: Arrival
				clock = t1 // Set clock to time of next arrival.
				if q.Full() {
					rejected <- Customer{ID: id,
						Arrival:          t1,
						Service:          serviceDistribution.Get(),
						Departure:        t1,
						QueueAtDeparture: q.Len(),
						Position:         -1}
				} else {
					q.Enqueue(Customer{ID: id,
						Arrival: t1,
						Service: serviceDistribution.Get()})
				}
				id++
				t1 = clock + arrivalDistribution.Get() // Set t1 to time of next arrival.
				if q.Len() == 1 {                      // If queue is not empty
					t2 = q.NextCompletion() // then set t2 to time of next completion.
				}
			} else { // If next arrival is after next completion -> Event: Departure
				if !math.IsInf(t2, 1) { // if next completion exists
					clock = t2                             // Set time to time of next completion.
					customer := q.Dequeue()                // Remove customer from queue
					customer.Departure = t2                // Set completion time
					customer.QueueAtDeparture = q.Len()    // Set queue size
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

// PrintCustomer prints "the arrival time, service time, time of departure of
// customers, as well as the number of customers in the system immediately
// after the departure of each of these customers"
func PrintCustomer(c Customer) {
	fmt.Printf("Customer %d (%d) | ", c.ID, c.Position)
	fmt.Printf("Arrival, Service, Departure = %.3f, %.3f, %.3f\n", c.Arrival, c.Service, c.Departure)
}

func logCustomer(c Customer) {
	log.Printf("Customer %02d (%02d) | Arrival, Service, Departure = %.3f, %.3f, %.3f\n", c.ID, c.Position, c.Arrival, c.Service, c.Departure)
}
