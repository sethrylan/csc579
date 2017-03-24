package mm1k

import (
	"log"
	"math"
	"math/rand"
)

// Return earliest NextCompletion for the list of queues
func nextCompletion(queues []Queue) (nextQueue int, nextCompletion float64) {
	nextCompletion = math.Inf(1)
	nextQueue = -1
	for i, q := range queues {
		if q.NextCompletion() < nextCompletion {
			nextQueue, nextCompletion = i, q.NextCompletion()
		}
	}
	log.Printf("nextQueue, nextCompletion = %d, %.3f\n", nextQueue, nextCompletion)
	return
}

// RunCPUIo will continually add and service customers using an event loop. At time
// t = 0 the system is empty. Draw a random number to decide when the
// first arrival will occur.
func RunCPUIo(arrivalDistribution Distribution, q []Queue, serviceDistributions []Distribution, qProbabilities map[int]float64) (rejects, completes, exits <-chan Customer) {
	rejected := make(chan Customer)  // Unbuffered channels ensure deterministic simulation
	completed := make(chan Customer) //
	exited := make(chan Customer)    //
	var clock float64                // master clock

	go func() {
		var t1 = arrivalDistribution.Get() // time of next arrival
		var t2q = -1                       // next queue with completing customer
		var t2 = math.Inf(1)               // time of next completion (∞ for no schedule Customer)
		var id int                         // Incremented Customer ID
		for {                              // Do forever
			if t1 < t2 { // If next arrival is before next completion -> Event: Arrival
				clock = t1 // Set clock to time of next arrival.
				if q[0].Full() {
					rejected <- Customer{ID: id,
						Arrival:          t1,
						Service:          serviceDistributions[0].Get(),
						Departure:        t1,
						QueueAtDeparture: q[0].Len(),
						Position:         -1}
				} else {
					q[0].Enqueue(Customer{ID: id,
						Arrival: t1,
						Service: serviceDistributions[0].Get()})
				}
				id++
				t1 = clock + arrivalDistribution.Get() // Set t1 to time of next arrival.
				t2q, t2 = nextCompletion(q)            // then set t2 to time of next completion.
			} else { // If next arrival is after next completion -> Event: Departure
				if !math.IsInf(t2, 1) { // if next completion exists
					clock = t2                               // Set time to time of next completion.
					customer := q[t2q].Dequeue()             // Remove customer from queue
					customer.Departure = t2                  // Set completion time
					customer.QueueAtDeparture = q[t2q].Len() // Set queue size
					customer.Start = t2 - customer.Service   // and start time
					completed <- customer                    // and add the customer to the completed channel
					if t2q >= 1 {
						q[0].Enqueue(customer)
					} else {
						r := rand.Float64()
						if r <= qProbabilities[0] {
							exited <- customer
						} else {
							for i := 0; i < len(qProbabilities)-1; i++ {
								if r > qProbabilities[i] && r <= qProbabilities[i+1] {
									if q[i+1].Full() {
										// TODO: reject
									} else {
										q[i+1].Enqueue(customer)
									}
								}
							}
						}
					}
				}
				t2q, t2 = nextCompletion(q) // then set t2 to time of next completion.
			}
		}
	}()

	rejects = rejected
	completes = completed
	exits = exited
	return
}

// Simulate will terminate once C customers have completed.
func SimulateCPU(λ float64, µs []float64, queues []Queue, C int, seed int64) (completes []Customer, rejects []Customer, exits []Customer) {
	transitions := map[int]float64{
		0: 0.7,
		1: 0.8,
		2: 0.9,
		3: 1.0,
	}
	serviceDistributions := []Distribution{}
	for i, µ := range µs {
		serviceDistributions = append(serviceDistributions, NewExpDistribution(µ, seed+int64(i)+1))
	}

	var customer Customer
	var rejected, completed, exited <-chan Customer
	rejected, completed, exited = RunCPUIo(
		NewExpDistribution(λ, seed),
		queues,
		serviceDistributions,
		transitions,
	)
	for len(exits) < C {
		select {
		case customer = <-rejected:
			rejects = append(rejects, customer)
			logCustomer(customer)
		case customer = <-completed:
			completes = append(completes, customer)
			logCustomer(customer)
		case customer = <-exited:
			exits = append(exits, customer)
			logCustomer(customer)
		}
	}
	return
}
