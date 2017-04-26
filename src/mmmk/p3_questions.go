package mmmk

import (
	"fmt"
)

func P3Question3(replications int, seed int64) {
	C := 50000
	µ := float64(1) / 3000
	α := 1.1
	k := 332
	p := 100000 * 100000

	fmt.Printf("========== M/G/3 ==========\n")

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

	fmt.Printf("========== M/M/3 ==========\n")
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
