package mmmk

import (
	"fmt"
	"math"
	"mm1k"
	"sort"
	"sync"
)

type Queuer interface {
	Enqueue(customer mm1k.Customer)
	Dequeue() mm1k.Customer
	Empty() bool
}

type Server interface {
	Add(customer mm1k.Customer)
	Dequeue() mm1k.Customer
	NextCompletion() mm1k.Customer
	Get() float64 // return RV from service distribution
	isFull() bool
}

type Queue interface {
	Queuer
}

type Service interface {
	Server
}

func Run(arrivalDistribution mm1k.Distribution, q Queue, s Service) (rejects, completes <-chan mm1k.Customer) {
	rejected := make(chan mm1k.Customer)  // Unbuffered channels ensure deterministic simulation
	completed := make(chan mm1k.Customer) //
	var clock float64                     // master clock

	go func() {
		var t1 = arrivalDistribution.Get() // time of next arrival
		var t2 = math.Inf(1)               // time of next completion (∞ for no schedule Customer)
		var id int                         // Incremented Customer ID
		for {                              // Do forever
			if t1 < t2 { // If next arrival is before next completion -> Event: Arrival
				clock = t1 // Set clock to time of next arrival.
				// never reject
				q.Enqueue(
					mm1k.Customer{
						ID:      id,
						Arrival: t1,
						Service: s.Get(),
					})
				id++
				if !s.isFull() {
					c := q.Dequeue()                  // Dequeue from queue
					c.Start = clock                   // set start to current time
					c.Departure = c.Start + c.Service // calculate depart time
					s.Add(c)                          // Add to server
				}
				t1 = clock + arrivalDistribution.Get() // Set t1 to time of next arrival.
				t2 = s.NextCompletion().Departure      // then set t2 to time of next completion.
			} else { // If next arrival is after next completion -> Event: Departure
				if !math.IsInf(t2, 1) { // if next completion exists
					clock = t2                                             // Set time to time of next completion.
					customer := s.Dequeue()                                // Remove customer from queue
					customer.Departure = customer.Start + customer.Service // Set completion time
					// customer.Start = t2 - customer.Service // and start time
					// mm1k.PrintCustomer(customer)
					completed <- customer // and add the customer to the completed channel
					if !q.Empty() {
						c := q.Dequeue()
						c.Start = t2
						c.Departure = c.Start + c.Service
						s.Add(c)
					}
				}
				t2 = s.NextCompletion().Departure
			}
		}
	}()
	rejects = rejected
	completes = completed
	return
}

func SimulateReplicationsMMMK(λ float64, queueType int, numServers int, µ float64, C int, replications int, seed int64) mm1k.SimMetricsList {
	metricsChannel := make(chan mm1k.SimMetrics, replications)
	metricsList := make([]mm1k.SimMetrics, 0)
	var wg sync.WaitGroup
	for i := 0; i < replications; i++ {
		wg.Add(1)
		server := MakeMinService(3, mm1k.NewExpDistribution(µ, seed+int64(i)))
		go replication(&wg, i, λ, queueType, server, C, seed, metricsChannel)
	}
	wg.Wait()
	close(metricsChannel)

	for metrics := range metricsChannel {
		metricsList = append(metricsList, metrics)
	}
	return metricsList
}

func SimulateReplicationsMGMK(λ float64, queueType int, numServers int, α float64, k int, p int, C int, replications int, seed int64) mm1k.SimMetricsList {
	metricsChannel := make(chan mm1k.SimMetrics, replications)
	metricsList := make([]mm1k.SimMetrics, 0)
	var wg sync.WaitGroup
	for i := 0; i < replications; i++ {
		wg.Add(1)
		server := MakeMinService(3, mm1k.NewParetoDistribution(α, k, p, seed+int64(i)))
		go replication(&wg, i, λ, queueType, server, C, seed+int64(i), metricsChannel)
	}
	wg.Wait()
	close(metricsChannel)

	for metrics := range metricsChannel {
		metricsList = append(metricsList, metrics)
	}
	return metricsList
}

func replication(wg *sync.WaitGroup, i int, λ float64, queueType int, server Service, C int, seed int64, ch chan<- mm1k.SimMetrics) {
	defer wg.Done()
	defer fmt.Printf(".")
	_, completes := Simulate(λ, queueType, server, C, seed)

	var metrics mm1k.SimMetrics
	metrics.Wait = mm1k.Mean(completes, mm1k.Wait)
	metrics.System = mm1k.Mean(completes, mm1k.System)
	sort.Sort(mm1k.ByDeparture(completes))
	metrics.LastDeparture = completes[len(completes)-1].Departure
	// metrics.LastDeparture = EmpiricalCLR(len(), len(rejectsByQueue[k])+len(completes))
	ch <- metrics
	return
}

// Simulate will terminate once C customers have completed.
func Simulate(λ float64, queueType int, server Service, C int, seed int64) (rejects []mm1k.Customer, completes []mm1k.Customer) {
	var customer mm1k.Customer
	var queue Queue
	var rejected, completed <-chan mm1k.Customer
	switch queueType {
	case 0:
		queue = NewFIFO()
	case 1:
		queue = NewSJF()
	}
	rejected, completed = Run(
		mm1k.NewExpDistribution(λ, seed),
		queue,
		server,
	)
	for len(completes) < C {
		select {
		case customer = <-rejected:
			rejects = append(rejects, customer)
			mm1k.LogCustomer(customer)
		case customer = <-completed:
			completes = append(completes, customer)
			mm1k.LogCustomer(customer)
		}
	}
	return
}

func PrintMetricsList(metricsList mm1k.SimMetricsList) {
	var maxDeparture, sampleMean, sampleStdDev float64
	for _, metrics := range metricsList {
		// fmt.Printf("metrics = %v\n", metrics)
		maxDeparture = math.Max(metrics.LastDeparture, maxDeparture)
	}
	fmt.Printf("\nClock     = %.0f (max of all replications)\n", maxDeparture)

	sampleMean, sampleStdDev = metricsList.MeanAndStdDev(mm1k.AverageWait)
	fmt.Printf("W̄ait      = %.0f±%.0f\n", sampleMean, sampleStdDev*2) // Print 95% confidence interval

	sampleMean, sampleStdDev = metricsList.MeanAndStdDev(mm1k.AverageSystem)
	fmt.Printf("S̄ystem    = %.0f±%.0f\n", sampleMean, sampleStdDev*2) // Print 95% confidence interval

	// fmt.Printf("  %.0f) +- (0.0, %.0f)\n", sampleMean, sampleStdDev*2) // Print 95% confidence interval

}
