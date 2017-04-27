package mm1k

import (
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
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
	PriorityQueue    int     // the number of priority queue that the customer occupies
	Server           int     // id of the server that served this customer
}

// ByID implements sort.Interface for []Customer
type ByID []Customer

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

// byService implements sort.Interface for []Customer
type ByService []Customer

func (a ByService) Len() int           { return len(a) }
func (a ByService) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByService) Less(i, j int) bool { return a[i].Service < a[j].Service }

// byDeparture implements sort.Interface for []Customer
type ByDeparture []Customer

func (a ByDeparture) Len() int           { return len(a) }
func (a ByDeparture) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDeparture) Less(i, j int) bool { return a[i].Departure < a[j].Departure }

// ByWait implements sort.Interface for []Customer
type ByWait []Customer

func (a ByWait) Len() int           { return len(a) }
func (a ByWait) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWait) Less(i, j int) bool { return Wait(a[i]) < Wait(a[j]) }

type field func(c Customer) float64

// Service is a field for sorting customer.
func Service(c Customer) float64 {
	return c.Service
}

// Wait is a field for sorting customer.
func Wait(c Customer) float64 {
	return c.Start - c.Arrival
}

// Departure is a field for sorting customer.
func Departure(c Customer) float64 {
	return c.Departure - c.Arrival
}

// System is a field for sorting customer.
func System(c Customer) float64 {
	return c.Departure - c.Arrival
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

	NextQueue() int
}

// SimulateReplications will terminate once C customers have completed.
func SimulateReplications(λ float64, µ float64, makerFunc func(int) Queue, K int, C int, replications int, discard int, seed int64) map[int]SimMetricsList {
	var metricsByQueue = make(map[int]SimMetricsList)
	metricsChannel := make(chan SimMetricsList, replications) // the metrics list will be len 0 < x <= p, where x is the number of queues for priority queues.

	var wg sync.WaitGroup
	for i := 0; i < replications; i++ {
		wg.Add(1)
		go replication(&wg, i, λ, µ, makerFunc(K), C, seed, metricsChannel)
	}
	wg.Wait()
	close(metricsChannel)

	for waitTimesArray := range metricsChannel {
		for queueIndex, w := range waitTimesArray {
			metricsByQueue[queueIndex] = append(metricsByQueue[queueIndex], w)
		}
	}
	return metricsByQueue
}

// Yield the metrics per queue (after discard)
func replication(wg *sync.WaitGroup, i int, λ float64, µ float64, queue Queue, C int, seed int64, ch chan<- SimMetricsList) {
	defer wg.Done()
	defer fmt.Printf(".")
	rejects, completes := Simulate(λ, µ, queue, C, seed+int64(i))
	completes = RemoveFirstNByDeparture(completes, discard)

	completesByQueue := make(map[int][]Customer)
	rejectsByQueue := make(map[int][]Customer)

	for _, c := range completes {
		completesByQueue[c.PriorityQueue] = append(completesByQueue[c.PriorityQueue], c)
	}

	for _, c := range rejects {
		rejectsByQueue[c.PriorityQueue] = append(rejectsByQueue[c.PriorityQueue], c)
	}

	metricsListByQueue := make(SimMetricsList, len(completesByQueue))
	for k, completes := range completesByQueue {
		metricsListByQueue[k].Wait = Mean(completes, Wait)
		metricsListByQueue[k].System = Mean(completes, System)
		sort.Sort(ByDeparture(completes))
		metricsListByQueue[k].LastDeparture = completes[len(completes)-1].Departure
		metricsListByQueue[k].CLR = EmpiricalCLR(len(rejectsByQueue[k]), len(rejectsByQueue[k])+len(completes))
	}
	ch <- metricsListByQueue
	return
}

// Simulate will terminate once C customers have completed.
func Simulate(λ float64, µ float64, q Queue, C int, seed int64) (rejects []Customer, completes []Customer) {
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

// Run will continually add and service customers using an event loop. At time
// t = 0 the system is empty. Draw a random number to decide when the
// first arrival will occur.
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
						Position:         -1,
						PriorityQueue:    q.NextQueue()}
				} else {
					q.Enqueue(Customer{ID: id,
						Arrival: t1,
						Service: serviceDistribution.Get()})
				}
				id++
				t1 = clock + arrivalDistribution.Get() // Set t1 to time of next arrival.
				t2 = q.NextCompletion()                // then set t2 to time of next completion.
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
	fmt.Printf("Customer %02d (%02d) | Arrival, Service, [Start, Departure] = %.3f, %.3f, [%.3f, %.3f]\n", c.ID, c.Position, c.Arrival, c.Service, c.Start, c.Departure)
}

func logCustomer(c Customer) {
	log.Printf("Customer %02d (%02d) | Arrival, Service, [Start, Departure] = %.3f, %.3f, [%.3f, %.3f]\n", c.ID, c.Position, c.Arrival, c.Service, c.Start, c.Departure)
}

func LogCustomer(c Customer) {
	logCustomer(c)
}

func PrintMetricsListQueueMap(metricsListByQueue map[int]SimMetricsList) {
	var keys []int
	var sampleMean, sampleStdDev float64
	var maxDeparture float64

	for k := range metricsListByQueue {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		maxDeparture = math.Max(metricsListByQueue[k][len(metricsListByQueue[k])-1].LastDeparture, maxDeparture)
	}

	fmt.Printf("\nClock        = %.3f (all queues, last replication)\n", maxDeparture)

	totalMetrics := make([]SimMetrics, 0)
	for _, k := range keys {
		for _, simMetric := range metricsListByQueue[k] {
			totalMetrics = append(totalMetrics, simMetric)
		}
		sampleMean, sampleStdDev = metricsListByQueue[k].MeanAndStdDev(AverageWait)
		if k == 0 {
			fmt.Printf("W̄ait         =")
		}
		fmt.Printf(" %.3f±%.3f", sampleMean, sampleStdDev*2) // Print 95% confidence interval
	}
	sampleMean, sampleStdDev = MeanAndStdDev(totalMetrics, AverageWait)
	fmt.Printf(" (Overall : %.3f±%.3f)", sampleMean, sampleStdDev*2) // Print 95% confidence interval
	fmt.Println()

	for _, k := range keys {
		sampleMean, sampleStdDev = metricsListByQueue[k].MeanAndStdDev(AverageSystem)
		if k == 0 {
			fmt.Printf("S̄ystem       =")
		}
		fmt.Printf(" %.3f±%.3f", sampleMean, sampleStdDev*2) // Print 95% confidence interval
	}
	fmt.Println()

	for _, k := range keys {
		sampleMean, sampleStdDev = metricsListByQueue[k].MeanAndStdDev(CLR)
		if k == 0 {
			fmt.Printf("CLR          =")
		}
		fmt.Printf(" %.3f±%.3f", sampleMean, sampleStdDev*2) // Print 95% confidence interval
	}
	sampleMean, sampleStdDev = MeanAndStdDev(totalMetrics, CLR)
	fmt.Printf(" (Overall : %.3f±%.3f)", sampleMean, sampleStdDev*2) // Print 95% confidence interval
	fmt.Println()
}

// QueueMakers is a sorted list of queue generator functions
var QueueMakers = []func(int) Queue{fifo, lifo, sjf, prioNP, prioP}

func fifo(K int) Queue {
	return NewFIFO(K)
}

func lifo(K int) Queue {
	return NewLIFO(K)
}

func sjf(K int) Queue {
	return NewSJF(K, false)
}

func prioNP(K int) Queue {
	return NewPriority(K, 4, false)
}

func prioP(K int) Queue {
	return NewPriority(K, 4, true)
}
