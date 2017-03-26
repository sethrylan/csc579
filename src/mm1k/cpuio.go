package mm1k

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"sync"
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

// Yield the metrics per queue (after discard)
func replicationCPUIO(wg *sync.WaitGroup, i int, λ float64, µs []float64, Ks []int, c int, seed int64, ch chan<- SimMetricsList) {
	defer wg.Done()
	defer fmt.Printf(".")
	queues := []Queue{}
	for _, k := range Ks {
		queues = append(queues, NewFIFO(k))
	}

	rejects, completes, exits := SimulateCPUIO(λ, µs, queues, c, seed)
	completes = RemoveFirstNByDeparture(completes, discard)
	exits = RemoveFirstNByDeparture(exits, discard)

	completesByQueue := make(map[int][]Customer)
	rejectsByQueue := make(map[int][]Customer)
	exitsByQueue := make(map[int][]Customer)

	for _, c := range completes {
		completesByQueue[c.PriorityQueue] = append(completesByQueue[c.PriorityQueue], c)
	}

	for _, c := range rejects {
		rejectsByQueue[c.PriorityQueue] = append(rejectsByQueue[c.PriorityQueue], c)
	}

	for _, c := range exits {
		exitsByQueue[c.PriorityQueue] = append(exitsByQueue[c.PriorityQueue], c)
	}

	metricsListByQueue := make(SimMetricsList, len(completesByQueue))
	for k, completes := range completesByQueue {
		metricsListByQueue[k].wait = Mean(completes, Wait)
		metricsListByQueue[k].system = Mean(completes, System)
		sort.Sort(byDeparture(completes))
		metricsListByQueue[k].lastDeparture = completes[len(completes)-1].Departure
		metricsListByQueue[k].clr = EmpiricalCLR(len(rejectsByQueue[k]), len(rejectsByQueue[k])+len(completes))
	}
	ch <- metricsListByQueue
	return
}

// RunCPUIO will continually add and service customers using an event loop. At time
// t = 0 the system is empty. Draw a random number to decide when the
// first arrival will occur.
func RunCPUIO(arrivalDistribution Distribution, queues []Queue, serviceDistributions []Distribution, qProbabilities map[int]float64) (rejects, completes, exits <-chan Customer) {
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
				if queues[0].Full() {
					rejected <- Customer{ID: id,
						Arrival:          t1,
						Service:          serviceDistributions[0].Get(),
						Departure:        t1,
						QueueAtDeparture: queues[0].Len(),
						Position:         -1,
						PriorityQueue:    0}
				} else {
					queues[0].Enqueue(Customer{ID: id,
						Arrival:       t1,
						Service:       serviceDistributions[0].Get(),
						PriorityQueue: 0})
				}
				id++
				t1 = clock + arrivalDistribution.Get() // Set t1 to time of next arrival.
				t2q, t2 = nextCompletion(queues)       // then set t2 to time of next completion.
			} else { // If next arrival is after next completion -> Event: Departure
				if !math.IsInf(t2, 1) { // if next completion exists
					clock = t2                                    // Set time to time of next completion.
					customer := queues[t2q].Dequeue()             // Remove customer from queue
					customer.Departure = t2                       // Set completion time
					customer.QueueAtDeparture = queues[t2q].Len() // Set queue size
					customer.Start = t2 - customer.Service        // and start time
					completed <- customer                         // and add the customer to the completed channel
					if t2q >= 1 {
						customer.PriorityQueue = 0
						if queues[0].Full() {
							rejected <- customer
						} else {
							queues[0].Enqueue(customer)
						}
					} else {
						r := rand.Float64()
						if r <= qProbabilities[0] {
							customer.PriorityQueue = 0
							exited <- customer
						} else {
							for i := 0; i < len(qProbabilities)-1; i++ {
								if r > qProbabilities[i] && r <= qProbabilities[i+1] {
									customer.PriorityQueue = i + 1
									if queues[i+1].Full() {
										rejected <- customer
									} else {
										customer.Service = serviceDistributions[i+1].Get()
										queues[i+1].Enqueue(customer)
									}
								}
							}
						}
					}
				}
				t2q, t2 = nextCompletion(queues) // then set t2 to time of next completion.
			}
		}
	}()

	rejects = rejected
	completes = completed
	exits = exited
	return
}

// Simulate will terminate once C customers have completed.
func SimulateCPUIO(λ float64, µs []float64, queues []Queue, C int, seed int64) (rejects []Customer, completes []Customer, exits []Customer) {
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
	rejected, completed, exited = RunCPUIO(
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

// SimulateReplicationsCPUIO will terminate once C customers have completed.
func SimulateReplicationsCPUIO(λ float64, µs []float64, Ks []int, C int, replications int, discard int, seed int64) map[int]SimMetricsList {
	var metricsByQueue = make(map[int]SimMetricsList)
	metricsChannel := make(chan SimMetricsList, replications) // the metrics list will be len 0 < x <= p, where x is the number of queues for priority queues.

	var wg sync.WaitGroup
	for i := 0; i < replications; i++ {
		wg.Add(1)
		go replicationCPUIO(&wg, i, λ, µs, Ks, C, seed, metricsChannel)
	}
	wg.Wait()
	close(metricsChannel)

	for metricsArray := range metricsChannel {
		for queueIndex, w := range metricsArray {
			metricsByQueue[queueIndex] = append(metricsByQueue[queueIndex], w)
		}
	}
	return metricsByQueue
}

func PrintMetricsListQueueMapCPUIO(metricsListByQueue map[int]SimMetricsList) {
	var keys []int
	var metricsForAllQueues SimMetricsList
	var sampleMean, sampleStdDev float64
	var maxDeparture float64

	for k := range metricsListByQueue {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		maxDeparture = math.Max(metricsListByQueue[0][len(metricsListByQueue[k])-1].lastDeparture, maxDeparture)
	}

	fmt.Printf("\nClock        = %.3f (last exit of CPU queue, last replication)\n", maxDeparture)

	for _, k := range keys {
		for _, metrics := range metricsListByQueue[k] {
			metricsForAllQueues = append(metricsForAllQueues, metrics)
		}
	}

	for _, k := range keys {
		sampleMean, sampleStdDev = metricsListByQueue[k].MeanAndStdDev(AverageWait)
		if k == 0 {
			fmt.Printf("W̄ait         =")
		}
		fmt.Printf(" %.3f±%.3f", sampleMean, sampleStdDev*2) // Print 95% confidence interval
	}
	sampleMean, sampleStdDev = metricsForAllQueues.MeanAndStdDev(AverageWait)
	fmt.Printf("\nW̄ait (All)   = %.3f±%.3f\n", sampleMean, sampleStdDev*2)

	for _, k := range keys {
		sampleMean, sampleStdDev = metricsListByQueue[k].MeanAndStdDev(AverageSystem)
		if k == 0 {
			fmt.Printf("S̄ystem       =")
		}
		fmt.Printf(" %.3f±%.3f", sampleMean, sampleStdDev*2) // Print 95% confidence interval
	}
	sampleMean, sampleStdDev = metricsForAllQueues.MeanAndStdDev(AverageSystem)
	fmt.Printf("\nS̄ystem (All) = %.3f±%.3f\n", sampleMean, sampleStdDev*2)

	for _, k := range keys {
		sampleMean, sampleStdDev = metricsListByQueue[k].MeanAndStdDev(CLR)
		if k == 0 {
			fmt.Printf("CLR          =")
		}
		fmt.Printf(" %.3f±%.3f", sampleMean, sampleStdDev*2) // Print 95% confidence interval
	}
	sampleMean, sampleStdDev = metricsForAllQueues.MeanAndStdDev(CLR)
	fmt.Printf("\nCLR (All)    = %.3f±%.3f\n", sampleMean, sampleStdDev*2)
}
