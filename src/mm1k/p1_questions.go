package mm1k

import (
	"fmt"
	"sort"
	"time"
)

// P1Question1 : Let the queue capacity K = 20. Plot the CLR against the value of ρ,
// for ρ = 0.05 to ρ = 0.95, in increments of 0.10. Submit two graphs: one for
// C = 1000 and one for C = 100000.
func P1Question1(seed int64) {
	K := 20
	for _, C := range []int{1000, 100000} {
		for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {
			fmt.Printf("%f, %d, %d | ", ρ, K, C)
			completes, rejects := Simulate(ρ, 1.0, NewFIFOQueue(K), C, seed)
			sorted := append(rejects, completes...)
			totalEvents := sorted[len(sorted)-1].ID + 1
			sort.Sort(ByID(sorted))
			// fmt.Printf("X/N = %d/%d\n", len(rejects), sorted[len(sorted)-1].ID+1)
			fmt.Printf("CLR (Empirical) = %.4f\n", EmpiricalCLR(len(rejects), totalEvents))
		}
	}
}

// P1Question2 : Now let us fix ρ = 0.85. Plot the CLR against the value of the queue
// capacity K, as K increases from 10 to 100 in increments of 10. Again, submit
// two graphs: one for C = 1000 and one for C = 100000.
func P1Question2(seed int64) {
	ρ := 0.85
	for _, C := range []int{1000, 100000} {
		for K := 10; K <= 100; K += 10 {
			fmt.Printf("%f, %d, %d | ", ρ, K, C)
			completes, rejects := Simulate(ρ, 1.0, NewFIFOQueue(K), C, seed)
			sorted := append(rejects, completes...)
			totalEvents := sorted[len(sorted)-1].ID + 1
			sort.Sort(ByID(sorted))
			fmt.Printf("CLR (Empirical) = %.4f\n", EmpiricalCLR(len(rejects), totalEvents))
		}
	}
}

// P1Question3 : Let K = 20. For C = 100000 and for ρ = 0.05 to 0.95 (in increments of 0.10),
// plot the simulation and analytical values of CLR on the same graph.
func P1Question3(seed int64) {
	K := 20
	C := 100000
	for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {
		fmt.Printf("%f, %d, %d | ", ρ, K, C)
		completes, rejects := Simulate(ρ, 1.0, NewFIFOQueue(K), C, seed)
		sorted := append(rejects, completes...)
		totalEvents := sorted[len(sorted)-1].ID + 1
		sort.Sort(ByID(sorted))
		fmt.Printf("CLR (Empirical, Analytical) =  (%.4f, %.4f\n",
			EmpiricalCLR(len(rejects), totalEvents),
			AnalyticalCLR(ρ, K))
	}
}

// P1Question4 : Let us set K = 100 and C = 100000. Compute the average waiting time W of the
// C customers that have received service at the end of the simulation (i.e.,
// ignore any lost customers or customers waiting in the queue when the
// simulation ends). Plot W against the value of ρ for ρ = 0.05 to 0.95.
func P1Question4(seed int64) {
	K := 100
	C := 100000
	for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {
		fmt.Printf("%f, %d, %d | ", ρ, K, C)
		completes, _ := Simulate(ρ, 1.0, NewFIFOQueue(K), C, seed)
		fmt.Printf("W̄ = %.3f\n", Mean(completes, Wait))
	}
}

// P1Question5 : Let us again set K = 40 and C = 100000. Time the running time of your
// simulation for ρ = 0.05 to 0.95. Plot the running time against the value of
// ρ. Note: turn off I/O
func P1Question5(seed int64) {
	K := 40
	C := 100000
	for ρ := 0.05; ρ <= 0.95; ρ += 0.10 {
		func() {
			defer timeTrack(time.Now(), fmt.Sprintf("\nρ = %.2f", ρ))
			Simulate(ρ, 1.0, NewFIFOQueue(K), C, seed)
		}()
	}
	fmt.Printf("\n")
}
