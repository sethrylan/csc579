package mm1k

import (
	"fmt"
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

func fifo(K int) Queue {
	return NewFIFO(K)
}

func lifo(K int) Queue {
	return NewLIFO(K)
}

func sjf(K int) Queue {
	return NewSJF(K)
}

func prioNP(K int) Queue {
	return NewPriority(K, true)
}

func prioP(K int) Queue {
	return NewPriority(K, true)
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

		for _, makerFunc := range []func(int) Queue{fifo, lifo, sjf, prioNP} {
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
				fmt.Printf("  W̄%d = %.3f ±%.3f@95%%", k, sampleMean, sampleStdDev*2)
			}
			fmt.Println()
			// } else {
			//  sampleMean := MeanFloats(waitTimes[0])
			//  sampleStdDev := StdDev(waitTimes[0], sampleMean)
			//  fmt.Printf("  W̄ = %.3f ±%.3f@95%%\n", sampleMean, sampleStdDev*2)
			// }
		}
	}
}

// Yield the average wait time per queue (after discard)
func replication(wg *sync.WaitGroup, i int, ρ float64, µ float64, queue Queue, C int, seed int64, ch chan<- []float64) {
	defer wg.Done()
	defer fmt.Printf(".")
	completes, _ := Simulate(ρ, µ, queue, C, seed+int64(i))
	completes = RemoveFirstNByDeparture(completes, discard)

	m := make(map[int][]Customer)
	for _, c := range completes {
		m[c.PriorityQueue] = append(m[c.PriorityQueue], c)
	}

	averageWaitTimes := make([]float64, len(m))
	for k := range m {
		averageWaitTimes[k] = Mean(m[k], Wait)
		// fmt.Printf("key[%s] value[%s]\n", k, m[k])
	}
	ch <- averageWaitTimes
	return
}
