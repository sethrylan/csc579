package mm1k

import "fmt"

var µ = 1.0
var discard = 1000
var replications = 30

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

		for _, queue := range []Queue{NewFIFO(K), NewLIFO(K), NewSJF(K)} {
			waitTimes := make([]float64, replications)
			fmt.Printf("% 10s %f, %d, %d | ", getType(queue), ρ, K, C)
			for i := 0; i < replications; i++ {
				fmt.Printf(".")
				completes, _ := Simulate(ρ, µ, queue, C, seed+int64(i))
				completes = removeFirstNByDeparture(completes, discard)
				waitTimes[i] = Mean(completes, Wait)
			}
			sampleMean := MeanFloats(waitTimes)
			sampleStdDev := StdDev(waitTimes, sampleMean)
			fmt.Printf("W̄ = %.3f ±%.3f@95%%\n", sampleMean, sampleStdDev*2)
		}

		//
		// fmt.Printf("LIFO %f, %d, %d | ", ρ, K, C)
		// completes, _ = Simulate(ρ, 1.0, NewLIFO(K), C, seed)
		// fmt.Printf("W̄ = %.3f\n", Mean(completes, Wait))
		//
		// fmt.Printf("SJF  %f, %d, %d | ", ρ, K, C)
		// completes, _ = Simulate(ρ, 1.0, NewSJF(K), C, seed)
		// fmt.Printf("W̄ = %.3f\n", Mean(completes, Wait))

		// TODO: ...
	}
}
