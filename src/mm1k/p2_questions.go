package mm1k

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
)

var µ = 1.0
var discard = 1000
var replications = 30

type SimMetrics struct {
	w float64
	s float64
}

type SimMetricsList []SimMetrics

func AverageWait(m SimMetrics) float64 {
	return m.w
}

func AverageService(m SimMetrics) float64 {
	return m.s
}

func (metrics SimMetricsList) MeanAndStdDev(fn func(c SimMetrics) float64) (mean float64, stdDev float64) {
	n := len(metrics)
	if n == 0 {
		return 0.0, 0.0
	}
	sum, squareSum := 0.0, 0.0
	for _, m := range metrics {
		sum += fn(m)
		squareSum += math.Pow(fn(m), 2)
	}
	mean = sum / float64(n)
	variance := squareSum/float64(n) - mean*mean
	return mean, math.Sqrt(variance)
}

func getFunctionName(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	return name[strings.LastIndex(name, ".")+1:]
}

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
			var metricsByQueue = make(map[int]SimMetricsList)
			fmt.Printf("% 10s %f, %d, %d | ", getFunctionName(makerFunc), ρ, K, C)
			metricsChannel := make(chan SimMetricsList, replications) // the metrics list will be len 0 < x <= p, where x is the number of queues for priority queues.

			var wg sync.WaitGroup
			for i := 0; i < replications; i++ {
				wg.Add(1)
				queue := makerFunc(K)
				go replication(&wg, i, ρ, µ, queue, C, seed, metricsChannel)
			}
			wg.Wait()
			close(metricsChannel)

			for waitTimesArray := range metricsChannel {
				for queueIndex, w := range waitTimesArray {
					metricsByQueue[queueIndex] = append(metricsByQueue[queueIndex], w)
				}
			}

			var keys []int
			for k := range metricsByQueue {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			for _, k := range keys {
				sampleMean, sampleStdDev := metricsByQueue[k].MeanAndStdDev(AverageWait)
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
func replication(wg *sync.WaitGroup, i int, ρ float64, µ float64, queue Queue, C int, seed int64, ch chan<- SimMetricsList) {
	defer wg.Done()
	defer fmt.Printf(".")
	completes, _ := Simulate(ρ, µ, queue, C, seed+int64(i))
	completes = RemoveFirstNByDeparture(completes, discard)

	customersGroupedByQueue := make(map[int][]Customer)
	for _, c := range completes {
		customersGroupedByQueue[c.PriorityQueue] = append(customersGroupedByQueue[c.PriorityQueue], c)
	}

	averageWaitTimesByQueue := make(SimMetricsList, len(customersGroupedByQueue))
	for k := range customersGroupedByQueue {
		averageWaitTimesByQueue[k].w = Mean(customersGroupedByQueue[k], Wait)
		log.Printf("W[%d]=%f\n, ", k, averageWaitTimesByQueue[k])
	}
	ch <- averageWaitTimesByQueue
	return
}
