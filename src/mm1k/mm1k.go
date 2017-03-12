package mm1k

import (
	"fmt"
	"log"
	"math"
)

type Customer struct {
	ID               int     // montonically increasing identifier
	Arrival          float64 // absolute time of arrival of the customer
	Service          float64 // relative time of service interal
	Start            float64 // absolution start time of service of the customer
	Departure        float64 // absolute time of completion/departure of the customer
	Position         int     // position of this customer in the queue at time of arrival
	QueueAtDeparture int     // number of customers in queue at time of departure
}

type Queue interface {
	Enqueue(customer Customer) Customer
	Dequeue() Customer
	Len() int
	NextCompletion() float64
	Full() bool
}

// Your simulation program will terminate once C customers have completed
// service, where C is an input parameter. For initial conditions, assume that
// at time t = 0 the system is empty. Draw a random number to decide when the
// first arrival will occur.
func Simulate(λ float64, µ float64, K int, C int, seed int64) (completes []Customer, rejects []Customer) {
	var customer Customer
	var rejected, completed <-chan Customer
	rejected, completed = Run(
		NewExpDistribution(λ, seed),
		NewFIFOQueue(K),
		NewExpDistribution(µ, seed+1),
	)
	for len(completes) < C {
		select {
		case customer = <-rejected:
			rejects = append(rejects, customer)
			LogCustomer(customer)
		case customer = <-completed:
			completes = append(completes, customer)
			LogCustomer(customer)
		}
	}
	return
}

func Run(arrivalDistribution Distribution, q Queue, serviceDistribution Distribution) (rejects, completes <-chan Customer) {
	rejected := make(chan Customer)  // Unbuffered channels ensure deterministic simulation
	completed := make(chan Customer) //
	var clock float64                // master clock

	go func() {
		var t1 float64 = arrivalDistribution.Get() // time of next arrival
		var t2 float64 = math.Inf(1)               // time of next completion (∞ for no schedule Customer)
		var id int                                 // Incremented Customer ID
		for {                                      // Do forever
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
				id += 1
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

// Print the arrival time, service time, time of departure of customers, as well
// as the number of customers in the system immediately after the departure of
// each of these customers
func PrintCustomer(c Customer) {
	fmt.Printf("Customer %d (%d) | ", c.ID, c.Position)
	fmt.Printf("Arrival, Service, Departure = %.3f, %.3f, %.3f\n", c.Arrival, c.Service, c.Departure)
}

func LogCustomer(c Customer) {
	log.Printf("Customer %02d (%02d) | Arrival, Service, Departure = %.3f, %.3f, %.3f\n", c.ID, c.Position, c.Arrival, c.Service, c.Departure)
}
