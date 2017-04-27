package mmmk

import (
	"fmt"
	"mm1k"
	"sort"
)

func P3Question3(replications int, seed int64) {
	C := 50000
	µ := float64(1) / 3000
	α := 1.1
	k := 332
	p := 100000 * 100000

	fmt.Printf("========== M/G/3 ==========\n")

	for λ := 0.0001; λ <= 0.0010; λ += 0.0001 {
		fmt.Printf("fifo %f, %d | ", λ, C)
		metricsList := SimulateReplicationsMGMK(λ, 0, 3, α, k, p, C, replications, seed)
		PrintMetricsList(metricsList)
	}

	for λ := 0.0001; λ <= 0.0010; λ += 0.0001 {
		fmt.Printf("sjf %f, %d | ", λ, C)
		metricsList := SimulateReplicationsMGMK(λ, 1, 3, α, k, p, C, replications, seed)
		PrintMetricsList(metricsList)
	}

	fmt.Printf("========== M/M/3 ==========\n")
	for λ := 0.0001; λ <= 0.0010; λ += 0.0001 {
		fmt.Printf("fifo %f, %d | ", λ, C)
		metricsList := SimulateReplicationsMMMK(λ, 0, 3, µ, C, replications, seed)
		PrintMetricsList(metricsList)
	}

	for λ := 0.0001; λ <= 0.0010; λ += 0.0001 {
		fmt.Printf("sfj %f, %d | ", λ, C)
		metricsList := SimulateReplicationsMMMK(λ, 1, 3, µ, C, replications, seed)
		PrintMetricsList(metricsList)
	}

}

func PrintPercentiles(customers []mm1k.Customer, binSize int) {
	sort.Sort(mm1k.ByWait(customers))
	for percentile, events := range split(customers, binSize) {
		mean, stddev := mm1k.MeanAndStd(events, mm1k.Wait)
		// fmt.Printf("%d, %.0f±%.0f\n", percentile+1, mean, stddev*2)
		// fmt.Printf("(%d, %.0f) +- (0.0, %.0f)\n", percentile+1, mean, stddev*2)
		fmt.Printf("%d\t%.0f\t%.0f\n", percentile+1, mean, stddev*2)
	}
}

func P3Question4(replications int, seed int64) {
	C := 100000
	µ := float64(1) / 3000
	α := 1.1
	k := 332
	p := 100000 * 100000
	λ := 0.0005
	bins := 100
	binSize := C / bins
	var server Service
	var completes []mm1k.Customer

	fmt.Printf("========== M/G/1 ==========\n")

	fmt.Printf("========== FIFO ========== \n")
	server = MakeMinService(1, mm1k.NewParetoDistribution(α, k, p, seed))
	_, completes = Simulate(λ, 0, server, C, seed)
	PrintPercentiles(completes, binSize)

	fmt.Printf("========== SJF ========== \n")
	server = MakeMinService(1, mm1k.NewParetoDistribution(α, k, p, seed))
	_, completes = Simulate(λ, 1, server, C, seed)
	PrintPercentiles(completes, binSize)

	fmt.Printf("========== M/M/1 ==========\n")

	fmt.Printf("========== FIFO ========== \n")
	server = MakeMinService(1, mm1k.NewExpDistribution(µ, seed))
	_, completes = Simulate(λ, 0, server, C, seed)
	PrintPercentiles(completes, binSize)

	fmt.Printf("========== SJF ========== \n")
	server = MakeMinService(1, mm1k.NewExpDistribution(µ, seed))
	_, completes = Simulate(λ, 1, server, C, seed)
	PrintPercentiles(completes, binSize)

}

func split(buf []mm1k.Customer, binSize int) [][]mm1k.Customer {
	var chunk []mm1k.Customer
	chunks := make([][]mm1k.Customer, 0, len(buf)/binSize+1)
	for len(buf) >= binSize {
		chunk, buf = buf[:binSize], buf[binSize:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
