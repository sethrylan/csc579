package mm1k

import (
	"fmt"
	"math"
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
	fmt.Printf("\n=======Starting P2Question3=======\n")
	kcpu := 50
	kio := 30
	C := 100000
	for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {
		fmt.Printf("\nρ=%f, Kcpu=%d, Kio=%d, C=%d\n", ρ, kcpu, kio, C)
		metricsListByQueue := SimulateReplicationsCPUIO(ρ, []float64{µcpu, µio, µio, µio}, []int{kcpu, kio, kio, kio}, C, replications, discard, seed)
		PrintMetricsListQueueMapCPUIO(metricsListByQueue)
	}
	fmt.Println()
}

func P2Question4(replications int, seed int64) {
	fmt.Printf("\n=======Starting P2Question4=======\n")
	kcpu := 40
	C := 100000
	for ρ := 0.10; ρ <= .60; ρ += 0.10 {
		minimumClrForIOQueues := math.Inf(1)
		for kio := 1; minimumClrForIOQueues >= 0.0001; kio++ {
			fmt.Printf("\nρ=%f, Kcpu=%d, Kio=%d, C=%d: ", ρ, kcpu, kio, C)
			// Using 1 replication
			metricsListByQueue := SimulateReplicationsCPUIO(ρ, []float64{µcpu, µio, µio, µio}, []int{kcpu, kio, kio, kio}, C, 1, discard, seed)
			for i := range []int{1, 2, 3} {
				minimumClrForIOQueues = math.Min(metricsListByQueue[i][0].CLR, minimumClrForIOQueues)
			}
			fmt.Printf("minimum CLR for IO Queues = %f ", minimumClrForIOQueues)
		}
	}
	fmt.Println()
}

func P2Question5(replications int, seed int64) {
	fmt.Printf("\n=======Starting P2Question5=======\n")
	kcpu := 50
	kio := 30
	C := 100000
	for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {
		fmt.Printf("\nρ=%f, Kcpu=%d, Kio=%d, C=%d\n", ρ, kcpu, kio, C)
		_, completes, _ := SimulateCPUIO(ρ, []float64{µcpu, µio, µio, µio}, []Queue{NewFIFO(kcpu), NewFIFO(kio), NewFIFO(kio), NewFIFO(kio)}, C, seed)
		completes = RemoveFirstNByDeparture(completes, discard)
		// exits = RemoveFirstNByDeparture(exits, discard)

		idsSet := make(map[int]struct{}) // unique list of ids
		visitsPerQueue := make([]int, 4)
		for _, c := range completes {
			idsSet[c.ID] = struct{}{}
			visitsPerQueue[c.PriorityQueue]++
		}
		fmt.Printf("Average per queue : ")
		for _, visits := range visitsPerQueue {
			fmt.Printf("%.3f  ", float64(visits)/float64(len(idsSet)))
		}

	}
	fmt.Println()
}
