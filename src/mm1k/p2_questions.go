package mm1k

import "fmt"

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
		fmt.Printf("FIFO %f, %d, %d | ", ρ, K, C)
		completes, _ := Simulate(ρ, 1.0, NewFIFO(K), C, seed)
		fmt.Printf("W̄ = %.3f\n", Mean(completes, Wait))

		fmt.Printf("LIFO %f, %d, %d | ", ρ, K, C)
		completes, _ = Simulate(ρ, 1.0, NewLIFO(K), C, seed)
		fmt.Printf("W̄ = %.3f\n", Mean(completes, Wait))

		// TODO: SJF
		// TODO: ...
	}
}
