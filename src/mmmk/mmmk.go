package mmmk

import (
	"fmt"
	"math"
	"mm1k"
	"sort"
	"sync"
)

type Queuer interface {
	Enqueue(arriveAt float64, serverAvailableAt float64) (serveAt float64, seatID int)
}

type Server interface {
	Serve(startAt float64) (departAt float64, serverID int)
}

type Nexter interface {
	Next() (nextAvailableAt float64)
}

type Queue interface {
	Queuer
	Nexter
}

type Service interface {
	Server
	Nexter
}

func Run(a mm1k.Distribution, q Queue, s Service) (rejs, deps <-chan mm1k.Customer) {
	rej := make(chan mm1k.Customer) // Unbuffered channels ensure deterministic simulation
	dep := make(chan mm1k.Customer) //
	var clock float64               // master clock

	go func() {
		var t0 float64           // time of next arrival
		var t1 float64           // time of next start
		var t2 float64           // time of next completion (∞ for no schedule Customer)
		var id int               // Incremented Customer ID
		var chs float64          // time of next available service
		var position, server int // position and server id
		for {
			t0 += a.Get()
			// no rejections

			t1, position = q.Enqueue(t0, chs)
			q.Next()
			// waited from t0 to t1

			t2, server = s.Serve(t1)
			chs = s.Next()
			// served from t1 to t2 by server

			dep <- mm1k.Customer{
				ID:        id,
				Arrival:   t0,
				Service:   t2 - t1,
				Start:     t1,
				Departure: t2,
				Position:  position,
				Server:    server,
			} // departed
			id++
			clock = t2
		}
	}()

	rejs = rej
	deps = dep
	return
}

func SimulateReplicationsMMMK(λ float64, queueType int, numServers int, µ float64, C int, replications int, seed int64) mm1k.SimMetricsList {
	metricsChannel := make(chan mm1k.SimMetrics, replications)
	metricsList := make([]mm1k.SimMetrics, replications)
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
	metricsList := make([]mm1k.SimMetrics, replications)
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
	_, completes := Simulate(λ, server, C, seed)

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
func Simulate(λ float64, server Service, C int, seed int64) (rejects []mm1k.Customer, completes []mm1k.Customer) {
	var customer mm1k.Customer
	var rejected, completed <-chan mm1k.Customer
	rejected, completed = Run(
		mm1k.NewExpDistribution(λ, seed),
		NewRing(1000000),
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
		maxDeparture = math.Max(metrics.LastDeparture, maxDeparture)
	}
	fmt.Printf("\nClock     = %.3f (max of all replications)\n", maxDeparture)

	sampleMean, sampleStdDev = metricsList.MeanAndStdDev(mm1k.AverageWait)
	fmt.Printf("W̄ait      = %.3f±%.3f\n", sampleMean, sampleStdDev*2) // Print 95% confidence interval

	sampleMean, sampleStdDev = metricsList.MeanAndStdDev(mm1k.AverageSystem)
	fmt.Printf("S̄ystem    = %.3f±%.3f\n", sampleMean, sampleStdDev*2) // Print 95% confidence interval

	sampleMean, sampleStdDev = metricsList.MeanAndStdDev(mm1k.CLR)
	fmt.Printf("CLR       = %.3f±%.3f\n", sampleMean, sampleStdDev*2) // Print 95% confidence interval
}
