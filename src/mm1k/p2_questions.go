package mm1k

import (
	"fmt"
	"sort"
)

var µ = 1.0
var discard = 1000

// P2Question1 : Let µ = 1 as before, the size of the queue K = 40, and the number
// of customers served before a simulation run terminates C = 100,000. Plot the
// average customer waiting time against the value of ρ, for ρ = 0.05 to
// ρ = 0.95, in increments of 0.10 (remember to include the confidence
// intervals). Compile five plots, one for each of the service disciplines. In
// particular, for the Priority service disciplines, plot the average waiting
// time of each of the four classes of customers, as well as the overall
// average.
func P2Question1(replications int, seed int64) {
	K := 40
	C := 100000
	for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {
		for _, makerFunc := range QueueMakers {
			fmt.Printf("% 10s %f, %d, %d | ", getFunctionName(makerFunc), ρ, K, C)

			metricsByQueue := SimulateReplications(ρ, µ, makerFunc, K, C, replications, discard, seed)
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
