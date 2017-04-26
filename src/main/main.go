package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mm1k"
	"mmmk"
	"os"
	"sort"
	"strconv"
)

var λ, µ, µcpu, µio, α float64
var k, kcpu, kio, c, l, m, p int

const seed int64 = 42
const discard int = 1000
const replications int = 30
const p3 = true

const usageMsg string = "λ K C L\n" +
	"λ = distribution of interarrival times\n" +
	"Kcpu = customers that the CPU queue may hold\n" +
	"Kio = customers that the IO queue may hold\n" +
	"C = customers served before the program terminates\n" +
	"L = 0–M/M/1 system, 1–CPU with I/O disks\n" +
	"M = 1–FCFS, 2–LCFS-NP, 3–SJF-NP, 4–Prio-NP, 5–Prio-P"

const usageMsgP3 string = "λ C L M\n" +
	"λ = distribution of interarrival times\n" +
	"C = customers served before the program terminates\n" +
	"L = 0-FCFS, 1-SJF\n" +
	"M = 0-M/M/3 system, 1-M/G/3"

func init() {

	if len(os.Args) < 7 && !p3 {
		fmt.Printf("usage: %s %s\n", os.Args[0], usageMsg)
		os.Exit(1)
	}

	if len(os.Args) < 5 && p3 {
		fmt.Printf("usage: %s %s\n", os.Args[0], usageMsgP3)
		os.Exit(1)
	}

	debugPtr := flag.Bool("debug", false, "a bool")
	flag.Parse()
	args := flag.Args()

	if p3 {
		λ, _ = strconv.ParseFloat(args[0], 64)
		c, _ = strconv.Atoi(args[1])
		l, _ = strconv.Atoi(args[2])
		m, _ = strconv.Atoi(args[3])

		µ = float64(1) / 3000
		α = 1.1
		k = 332
		p = 100000 * 100000
	} else {
		λ, _ = strconv.ParseFloat(args[0], 64)
		kcpu, _ = strconv.Atoi(args[1])
		kio, _ = strconv.Atoi(args[2])
		c, _ = strconv.Atoi(args[3])
		l, _ = strconv.Atoi(args[4])
		m, _ = strconv.Atoi(args[5])

		µ = 1.0
		µcpu = 1.0
		µio = 0.5
	}

	log.SetFlags(log.Lshortfile)
	if *debugPtr {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	if p3 {
		if c <= discard {
			fmt.Printf("WARNING: first %d events are discarded in metric calculations\n", discard)
		}
		mgmkSimulationWithReplication(seed)

	} else {
		switch l {
		case 0:
			if c <= discard {
				fmt.Printf("WARNING: first %d events are discarded in metric calculations\n", discard)
			}
			mm1kSimulationWithReplication(seed)
		case 1:
			cpuSimulationWithReplication(seed)
		default:
			fmt.Printf("usage: %s %s\n", os.Args[0], usageMsg)
			os.Exit(1)
		}
	}

	// mm1k.P2Question1(replications, seed)
	// mm1k.P2Question2a(replications, seed)
	// mm1k.P2Question3(replications, seed)
	// mm1k.P2Question4(replications, seed)
	// mm1k.P2Question5(replications, seed)
}

// P2 implementation
func mm1kSimulationWithReplication(seed int64) {
	fmt.Printf("======= Running m/m/1/k Simulation =======\n")
	if m > 5 || m < 1 {
		fmt.Printf("usage: %s %s\n", os.Args[0], usageMsg)
		os.Exit(1)
	}
	fmt.Printf("λ =     %.3f\n", λ)
	fmt.Printf("µ =     %.3f\n", µcpu)
	fmt.Printf("K =     %d\n", kcpu)
	fmt.Printf("C =     %d\n", c)
	fmt.Printf("L =     %d\n", l)
	fmt.Printf("M =     %s\n", mm1k.GetFunctionName(mm1k.QueueMakers[m-1]))
	fmt.Printf("%s ", mm1k.GetFunctionName(mm1k.QueueMakers[m-1]))
	metricsListByQueue := mm1k.SimulateReplications(λ, µ, mm1k.QueueMakers[m-1], kcpu, c, replications, discard, seed)
	mm1k.PrintMetricsListQueueMap(metricsListByQueue)
}

// P2 implementation
func cpuSimulationWithReplication(seed int64) {
	fmt.Printf("======= Running CPU/IO Simulation =======\n")
	fmt.Printf("λ =     %.3f\n", λ)
	fmt.Printf("µcpu =  %.3f\n", µcpu)
	fmt.Printf("µio =   %.3f\n", µio)
	fmt.Printf("Kcpu =  %d\n", kcpu)
	fmt.Printf("Kio =   %d\n", kio)
	fmt.Printf("C =     %d\n", c)
	fmt.Printf("L =     %d\n", l)

	metricsListByQueue := mm1k.SimulateReplicationsCPUIO(λ, []float64{µcpu, µio, µio, µio}, []int{kcpu, kio, kio, kio}, c, replications, discard, seed)
	mm1k.PrintMetricsListQueueMapCPUIO(metricsListByQueue)
}

// P2 implementation
// func cpuSimulation(seed int64) {
// 	fmt.Printf("======= Running CPU/IO Simulation =======\n")
// 	fmt.Printf("λ =     %.3f\n", λ)
// 	fmt.Printf("µcpu =  %.3f\n", µcpu)
// 	fmt.Printf("µio =   %.3f\n", µio)
// 	fmt.Printf("Kcpu =  %d\n", kcpu)
// 	fmt.Printf("Kio =   %d\n", kio)
// 	fmt.Printf("C =     %d\n", c)
// 	fmt.Printf("L =     %d\n", l)
//
// 	rejects, completes, exits := mm1k.SimulateCPUIO(λ, []float64{µcpu, µio, µio, µio}, []mm1k.Queue{mm1k.NewFIFO(kcpu), mm1k.NewFIFO(kio), mm1k.NewFIFO(kio), mm1k.NewFIFO(kio)}, c, seed)
// 	sorted := append(rejects, exits...)
// 	sorted = append(sorted, completes...)
// 	sort.Sort(mm1k.ByID(sorted))
// 	totalEvents := sorted[len(sorted)-1].ID + 1
// 	sampleServiceTimeMean := mm1k.Mean(completes, mm1k.Service)
//
// 	fmt.Printf("Master clock =                   %.3f\n", completes[len(completes)-1].Departure)
// 	fmt.Printf("CLR (Analytical) =               %.3f\n", mm1k.AnalyticalCLR(λ, kcpu))
// 	fmt.Printf("CLR (Empirical; X/N = %d/%d) =   %.3f\n", len(rejects), totalEvents, mm1k.EmpiricalCLR(len(rejects), totalEvents))
// 	fmt.Printf("Mean Service Time (S̄) =          %.3f\n", sampleServiceTimeMean)
// 	fmt.Printf("Mean Wait Time (W̄) =             %.3f\n", mm1k.Mean(completes, mm1k.Wait))
// 	fmt.Printf("Analytical Wait Time (W̄) =       %.3f\n", mm1k.AnalyticalWaitTime(λ, kcpu))
// }

// P1 implementation
func mm1kSimulation(seed int64) {
	rejects, completes := mm1k.Simulate(λ, µ, mm1k.QueueMakers[m-1](kcpu), c, seed)
	sorted := append(rejects, completes...)
	sort.Sort(mm1k.ByID(sorted))
	totalEvents := sorted[len(sorted)-1].ID + 1
	sampleServiceTimeMean := mm1k.Mean(completes, mm1k.Service)

	fmt.Printf("Master clock =                   %.3f\n", completes[len(completes)-1].Departure)
	fmt.Printf("CLR (Analytical) =               %.3f\n", mm1k.AnalyticalCLR(λ, kcpu))
	fmt.Printf("CLR (Empirical; X/N = %d/%d) =   %.3f\n", len(rejects), totalEvents, mm1k.EmpiricalCLR(len(rejects), totalEvents))
	fmt.Printf("Mean Service Time (S̄) =          %.3f\n", sampleServiceTimeMean)
	fmt.Printf("Mean Wait Time (W̄) =             %.3f\n", mm1k.Mean(completes, mm1k.Wait))
	fmt.Printf("Analytical Wait Time (W̄) =       %.3f\n", mm1k.AnalyticalWaitTime(λ, kcpu))

	for _, customer := range sorted {
		// L, L + 1, L + 10, and L + 11
		if customer.ID == l || customer.ID == l+1 || customer.ID == l+10 || customer.ID == l+11 {
			mm1k.PrintCustomer(customer)
		}
	}

}

func mgmkSimulationWithReplication(seed int64) {
	// var rejects, completes []mm1k.Customer
	fmt.Printf("======= Running m/m/1/k Simulation =======\n")
	if m > 1 || m < 0 {
		fmt.Printf("usage: %s %s\n", os.Args[0], usageMsgP3)
		os.Exit(1)
	}
	fmt.Printf("λ =     %.4f\n", λ)
	fmt.Printf("C =     %d\n", c)
	fmt.Printf("L =     %s\n", mm1k.GetFunctionName(mm1k.QueueMakers[l]))

	var metricsList mm1k.SimMetricsList
	if m == 0 { // MMMK (exponential service times)
		fmt.Printf("M =     M/M/3\n")
		metricsList = mmmk.SimulateReplicationsMMMK(λ, l, 3, µ, c, replications, seed)
		mmmk.PrintMetricsList(metricsList)
	} else { // MMGK (pareto service times)
		fmt.Printf("M =     M/G/3\n")
		metricsList = mmmk.SimulateReplicationsMGMK(λ, l, 3, α, k, p, c, replications, seed)
		mmmk.PrintMetricsList(metricsList)
	}
}
