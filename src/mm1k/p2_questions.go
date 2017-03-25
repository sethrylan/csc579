package mm1k

import (
	"fmt"
	"time"
)

var µ = 1.0
var µcpu = 1.0
var µio = 0.5
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
	fmt.Printf("\n=======Starting P2Question1=======\n")
	K := 40
	C := 100000
	for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {
		for _, makerFunc := range QueueMakers {
			fmt.Printf("\n%s ρ=%f, K=%d, C=%d | ", GetFunctionName(makerFunc), ρ, K, C)
			metricsByQueue := SimulateReplications(ρ, µ, makerFunc, K, C, replications, discard, seed)
			PrintMetricsListQueueMap(metricsByQueue)
		}
	}
}

func P2Question2(replications int, seed int64) {
	fmt.Printf("\n=======Starting P2Question2=======\n")
	K := 20
	C := 1000000
	for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {
		fmt.Printf("\nρ=%f, K=%d, C=%d (", ρ, K, C)
		for _, makerFunc := range QueueMakers {
			fmt.Printf("%s ", GetFunctionName(makerFunc))
		}
		fmt.Printf("\b) : ")

		for _, makerFunc := range QueueMakers {
			// fmt.Printf("\n%s %f, %d, %d | ", GetFunctionName(makerFunc), ρ, K, C)
			func() {
				defer timeTrack(time.Now(), "\t")
				Simulate(ρ, 1.0, makerFunc(K), C, seed)
			}()
		}
	}
	fmt.Println()
}

func P2Question3(replications int, seed int64) {
	fmt.Printf("\n=======Starting P2Question2=======\n")
	kcpu := 50
	kio := 30
	C := 100000
	for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {
		fmt.Printf("\nρ=%f, Kcpu=%d, Kio=%d, C=%d\n", ρ, kcpu, kio, C)
		func() {
			metricsListByQueue := SimulateReplicationsCPUIO(ρ, []float64{µcpu, µio, µio, µio}, []int{kcpu, kio, kio, kio}, C, replications, discard, seed)
			PrintMetricsListQueueMapCPUIO(metricsListByQueue)
		}()
	}
	fmt.Println()
}
