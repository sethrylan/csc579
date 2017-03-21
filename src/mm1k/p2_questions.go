package mm1k

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
)

var µ = 1.0
var discard = 1000
var replications = 30

type queueTime struct {
	queue int
	time  float64
}

func getFunctionName(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	return name[strings.LastIndex(name, ".")+1:]
}

type queueTimeMap map[int][]float64

// NewFIFO(K), NewLIFO(K), NewSJF(K), NewPriority(K, true)

// P2Question1 : Let µ = 1 as before, the size of the queue K = 40, and the number
// of customers served before a simulation run terminates C = 100,000. Plot the
// average customer waiting time against the value of ρ, for ρ = 0.05 to
// ρ = 0.95, in increments of 0.10 (remember to include the confidence
// intervals). Compile five plots, one for each of the service disciplines. In
// particular, for the Priority service disciplines, plot the average waiting
// time of each of the four classes of customers, as well as the overall
// average.
func P2Question1(seed int64) {
	K := 40
	C := 100000
	for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {

		for _, makerFunc := range QueueMakers {
			var averageWaitTimes = make(queueTimeMap)
			fmt.Printf("% 10s %f, %d, %d | ", getFunctionName(makerFunc), ρ, K, C)
			averageWaitTimesChannel := make(chan []float64, replications)

			var wg sync.WaitGroup
			for i := 0; i < replications; i++ {
				wg.Add(1)
				queue := makerFunc(K)
				go replication(&wg, i, ρ, µ, queue, C, seed, averageWaitTimesChannel)
			}
			wg.Wait()
			close(averageWaitTimesChannel)

			for waitTimesArray := range averageWaitTimesChannel {
				for queueIndex, w := range waitTimesArray {
					averageWaitTimes[queueIndex] = append(averageWaitTimes[queueIndex], w)
				}
			}

			var keys []int
			for k := range averageWaitTimes {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			for _, k := range keys {
				sampleMean := MeanFloats(averageWaitTimes[k])
				sampleStdDev := StdDev(averageWaitTimes[k], sampleMean)
				fmt.Printf("  W̄%d = %.3f ±%.3f", k, sampleMean, sampleStdDev*2)
				if k == 0 {
					fmt.Printf("@95%%")
				}
			}
			fmt.Println()
		}
	}
}

// Yield the average wait time per queue (after discard)
func replication(wg *sync.WaitGroup, i int, ρ float64, µ float64, queue Queue, C int, seed int64, ch chan<- []float64) {
	defer wg.Done()
	defer fmt.Printf(".")
	completes, _ := Simulate(ρ, µ, queue, C, seed+int64(i))
	completes = RemoveFirstNByDeparture(completes, discard)

	customersGroupedByQueue := make(map[int][]Customer)
	for _, c := range completes {
		customersGroupedByQueue[c.PriorityQueue] = append(customersGroupedByQueue[c.PriorityQueue], c)
	}

	averageWaitTimes := make([]float64, len(customersGroupedByQueue))
	for k := range customersGroupedByQueue {
		averageWaitTimes[k] = Mean(customersGroupedByQueue[k], Wait)
		log.Printf("W[%d]=%f\n, ", k, averageWaitTimes[k])
	}
	ch <- averageWaitTimes
	return
}
