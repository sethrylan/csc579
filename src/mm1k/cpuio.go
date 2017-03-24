package mm1k

import "math"

// Run will continually add and service customers using an event loop. At time
// t = 0 the system is empty. Draw a random number to decide when the
// first arrival will occur.
func RunCpuIo(arrivalDistribution Distribution, queues []Queue, serviceDistributions []Distribution, queueTransitionProbabilities map[int]float64) (rejects, completes <-chan Customer) {
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
				// if q.Len() > 0 {                       // If queue is not empty
				t2 = q.NextCompletion() // then set t2 to time of next completion.
				// }
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
