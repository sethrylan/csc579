package mmmk

import (
	"fmt"
)

// P1Question1 : Let the queue capacity K = 20. Plot the CLR against the value of ρ,
// for ρ = 0.05 to ρ = 0.95, in increments of 0.10. Submit two graphs: one for
// C = 1000 and one for C = 100000.
func P3Question3(seed int64) {
	C := 50000
	µ := float64(1) / 3000
	α := 1.1
	k := 332
	p := 100000 * 100000
	replications := 30
	for λ := 0.0001; λ <= 0.0009; λ += 0.0001 {
		fmt.Printf("fifo %f, %d | ", λ, C)
		metricsList := SimulateReplicationsMGMK(λ, 0, 3, α, k, p, C, replications, seed)
		PrintMetricsList(metricsList)
	}

	for λ := 0.0001; λ <= 0.0009; λ += 0.0001 {
		fmt.Printf("sjf %f, %d | ", λ, C)
		metricsList := SimulateReplicationsMGMK(λ, 1, 3, α, k, p, C, replications, seed)
		PrintMetricsList(metricsList)
	}

	for λ := 0.0001; λ <= 0.0009; λ += 0.0001 {
		fmt.Printf("fifo %f, %d | ", λ, C)
		metricsList := SimulateReplicationsMMMK(λ, 0, 3, µ, C, replications, seed)
		PrintMetricsList(metricsList)
	}

	for λ := 0.0001; λ <= 0.0009; λ += 0.0001 {
		fmt.Printf("sfj %f, %d | ", λ, C)
		metricsList := SimulateReplicationsMMMK(λ, 1, 3, µ, C, replications, seed)
		PrintMetricsList(metricsList)
	}

}
