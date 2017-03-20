package mm1k

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

var µ = 1.0
var discard = 1000
var replications = 30

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
			var waitTimes []float64
			fmt.Printf("% 10s %f, %d, %d | ", getFunctionName(makerFunc), ρ, K, C)
			waitTimesChannel := make(chan float64, replications)

			var wg sync.WaitGroup
			for i := 0; i < replications; i++ {
				wg.Add(1)
				queue := makerFunc(K)
				go replication(&wg, i, ρ, µ, queue, C, seed, waitTimesChannel)
			}
			wg.Wait()
			close(waitTimesChannel)

			for elem := range waitTimesChannel {
				waitTimes = append(waitTimes, elem)
			}

			sampleMean := MeanFloats(waitTimes)
			sampleStdDev := StdDev(waitTimes, sampleMean)
			fmt.Printf("  W̄ = %.3f ±%.3f@95%%\n", sampleMean, sampleStdDev*2)
		}
	}
}

func replication(wg *sync.WaitGroup, i int, ρ float64, µ float64, queue Queue, C int, seed int64, ch chan<- float64) {
	defer wg.Done()
	defer fmt.Printf(".")
	completes, _ := Simulate(ρ, µ, queue, C, seed+int64(i))
	completes = RemoveFirstNByDeparture(completes, discard)
	ch <- Mean(completes, Wait)
	return
}
